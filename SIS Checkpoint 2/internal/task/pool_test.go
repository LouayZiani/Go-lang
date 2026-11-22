package task

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewPoolConfig(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	tests := []struct {
		name     string
		svc      *Service
		workers  int
		capacity int
	}{
		{"nil service", nil, 2, 1},
		{"one worker", svc, 1, 1},
		{"zero workers", svc, 0, 1},
		{"negative workers", svc, -3, 1},
		{"zero capacity", svc, 2, 0},
		{"negative capacity", svc, 2, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPool(tt.svc, tt.workers, tt.capacity)
			if !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("want ErrInvalidConfig, got %v", err)
			}
			if p != nil {
				t.Error("pool must be nil on error")
			}
		})
	}
}

func TestPoolStartup(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	p := mustPool(t, svc, 3, 5)

	if p.Workers() != 3 || p.QueueCapacity() != 5 {
		t.Errorf("Workers=%d QueueCapacity=%d", p.Workers(), p.QueueCapacity())
	}
	if p.Running() != 3 {
		t.Errorf("Running = %d, want 3", p.Running())
	}
	if p.Stopped() {
		t.Error("new pool must not be stopped")
	}
	if p.QueueLen() != 0 {
		t.Errorf("QueueLen = %d, want 0", p.QueueLen())
	}
}

func TestSubmitDoesNotWaitForProcessing(t *testing.T) {
	g := newGate()
	svc := NewService(g)
	p := mustPool(t, svc, 2, 2)
	t.Cleanup(g.open) // runs before p.Stop so blocked workers are released

	tk := mustCreate(t, svc, input(12))
	// The gate is closed, so processing cannot finish. If Submit waited for the result, this call would hang forever
	if err := p.Submit(tk.ID); err != nil {
		t.Fatalf("Submit: %v", err)
	}
	g.waitStarted(t, 1)
	if got := mustGet(t, svc, tk.ID); got.State != Processing {
		t.Fatalf("state = %s, want Processing", got.State)
	}

	g.open()
	if got := waitDone(t, svc, tk.ID); got.State != Completed {
		t.Fatalf("state = %s, want Completed", got.State)
	}
}

func TestSubmitAndProcessMany(t *testing.T) {
	const n = 40
	svc := NewService(echoProcessor)
	p := mustPool(t, svc, 4, n)

	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = mustCreate(t, svc, input(float64(i+1))).ID
	}
	for _, id := range ids {
		if err := p.Submit(id); err != nil {
			t.Fatalf("Submit(%s): %v", id, err)
		}
	}
	for i, id := range ids {
		got := waitDone(t, svc, id)
		if got.State != Completed || got.Result == nil {
			t.Fatalf("task %s: %+v", id, got)
		}
		if want := float64(i + 1); got.Result.Ratio != want {
			t.Errorf("task %s: ratio = %v, want %v (result landed in the wrong task)", id, got.Result.Ratio, want)
		}
	}
}

func TestRealProcessorThroughPool(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	p := mustPool(t, svc, 2, 4)

	pass := mustCreate(t, svc, Input{"#000", "#fff", 12, false})
	fail := mustCreate(t, svc, Input{"#777777", "#ffffff", 12, false})
	for _, tk := range []Task{pass, fail} {
		if err := p.Submit(tk.ID); err != nil {
			t.Fatal(err)
		}
	}
	if got := waitDone(t, svc, pass.ID); got.State != Completed || got.Result.Level != LevelAAA {
		t.Errorf("pass: %+v", got)
	}
	//here a contrast "Fail" verdict is still a successfully processed task
	if got := waitDone(t, svc, fail.ID); got.State != Completed || got.Result.Level != LevelFail {
		t.Errorf("fail: %+v", got)
	}
}

// Proof of real overlap: two tasks are inside Process at the same moment. Neither can leave until the test opens the gate, so if the pool processed
// tasks one after another, the second "started" signal would never arrive.
func TestTasksProcessConcurrently(t *testing.T) {
	g := newGate()
	svc := NewService(g)
	p := mustPool(t, svc, 2, 4)
	t.Cleanup(g.open)

	a := mustCreate(t, svc, input(1))
	b := mustCreate(t, svc, input(2))
	if err := p.Submit(a.ID); err != nil {
		t.Fatal(err)
	}
	if err := p.Submit(b.ID); err != nil {
		t.Fatal(err)
	}

	g.waitStarted(t, 2) //both inside Process simultaneously

	if got := mustGet(t, svc, a.ID); got.State != Processing {
		t.Errorf("a: state = %s, want Processing", got.State)
	}
	if got := mustGet(t, svc, b.ID); got.State != Processing {
		t.Errorf("b: state = %s, want Processing", got.State)
	}

	g.open()
	for _, id := range []string{a.ID, b.ID} {
		if got := waitDone(t, svc, id); got.State != Completed {
			t.Errorf("%s: state = %s", id, got.State)
		}
	}
}

// The number of tasks being processed at once never exceeds the worker count, i.e. the pool reuses workers instead of spawning a goroutine per task.
func TestConcurrencyIsBoundedByWorkers(t *testing.T) {
	const workers, n = 3, 60
	var current, peak atomic.Int32
	proc := fakeProcessor(func(in Input) (Result, error) {
		c := current.Add(1)
		for {
			old := peak.Load()
			if c <= old || peak.CompareAndSwap(old, c) {
				break
			}
		}
		runtime.Gosched() //give other workers a chance to overlap
		current.Add(-1)
		return Result{Ratio: in.PointSize}, nil
	})
	svc := NewService(proc)
	p := mustPool(t, svc, workers, n)

	ids := make([]string, n)
	for i := range ids {
		ids[i] = mustCreate(t, svc, input(float64(i+1))).ID
	}
	for _, id := range ids {
		if err := p.Submit(id); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range ids {
		waitDone(t, svc, id)
	}
	if got := peak.Load(); got > workers {
		t.Fatalf("peak concurrency = %d, want <= %d", got, workers)
	}
}

func TestSubmitInvalid(t *testing.T) {
	g := newGate()
	svc := NewService(g)
	p := mustPool(t, svc, 2, 4)
	t.Cleanup(g.open)

	if err := p.Submit("nope"); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("unknown id: want ErrTaskNotFound, got %v", err)
	}

	tk := mustCreate(t, svc, input(12))
	if err := p.Submit(tk.ID); err != nil {
		t.Fatal(err)
	}
	g.waitStarted(t, 1) // now Processing

	// Duplicate submit while Processing
	if err := p.Submit(tk.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("submit while Processing: want ErrInvalidTransition, got %v", err)
	}
	// Synchronous Process on a task a worker is running
	if _, err := svc.Process(tk.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("Process while Processing: want ErrInvalidTransition, got %v", err)
	}

	g.open()
	done := waitDone(t, svc, tk.ID)
	// Submit after Completed
	if err := p.Submit(tk.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("submit when Completed: want ErrInvalidTransition, got %v", err)
	}
	after := mustGet(t, svc, tk.ID)
	if after.State != Completed || after.Result == nil || *after.Result != *done.Result {
		t.Fatalf("rejected submits changed the task: %+v", after)
	}
}

func TestSubmitQueuedTaskTwice(t *testing.T) {
	g := newGate()
	svc := NewService(g)
	p := mustPool(t, svc, 2, 4)
	t.Cleanup(g.open)

	//Occupy both workers so the next task stays in the queue.
	for i := 0; i < 2; i++ {
		if err := p.Submit(mustCreate(t, svc, input(float64(i+1))).ID); err != nil {
			t.Fatal(err)
		}
	}
	g.waitStarted(t, 2)

	queued := mustCreate(t, svc, input(10))
	if err := p.Submit(queued.ID); err != nil {
		t.Fatal(err)
	}
	if got := mustGet(t, svc, queued.ID); got.State != Queued {
		t.Fatalf("state = %s, want Queued", got.State)
	}
	if p.QueueLen() != 1 {
		t.Fatalf("QueueLen = %d, want 1", p.QueueLen())
	}
	if err := p.Submit(queued.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("want ErrInvalidTransition, got %v", err)
	}
	if p.QueueLen() != 1 {
		t.Fatalf("QueueLen = %d after rejected duplicate, want 1", p.QueueLen())
	}
}

func TestFullQueueRejectsAndRecovers(t *testing.T) {
	const workers, capacity = 2, 2
	g := newGate()
	svc := NewService(g)
	p := mustPool(t, svc, workers, capacity)
	t.Cleanup(g.open)

	create := func(size float64) string { return mustCreate(t, svc, input(size)).ID }
	var accepted []string

	//so first step is we fill both workers, and wait until both hold a task (queue is empty)
	for i := 0; i < workers; i++ {
		id := create(float64(i + 1))
		if err := p.Submit(id); err != nil {
			t.Fatal(err)
		}
		accepted = append(accepted, id)
	}
	g.waitStarted(t, workers)

	//second step is to fill the queue exactly
	for i := 0; i < capacity; i++ {
		id := create(float64(10 + i))
		if err := p.Submit(id); err != nil {
			t.Fatalf("Submit into free slot: %v", err)
		}
		accepted = append(accepted, id)
	}
	if p.QueueLen() != capacity {
		t.Fatalf("QueueLen = %d, want %d", p.QueueLen(), capacity)
	}

	//third one is when we do one more: rejected, and the task is still a clean Pending task
	extra := create(99)
	if err := p.Submit(extra); !errors.Is(err, ErrQueueFull) {
		t.Fatalf("want ErrQueueFull, got %v", err)
	}
	got := mustGet(t, svc, extra)
	if got.State != Pending || got.Result != nil {
		t.Fatalf("rejected task must stay Pending, got %+v", got)
	}
	if p.QueueLen() != capacity {
		t.Fatalf("queue grew beyond capacity: %d", p.QueueLen())
	}

	//fourth step is free space: everything accepted finishes, then the rejected task can be resubmitted
	g.open()
	for _, id := range accepted {
		if got := waitDone(t, svc, id); got.State != Completed {
			t.Fatalf("%s: %s", id, got.State)
		}
	}
	if err := p.Submit(extra); err != nil {
		t.Fatalf("resubmit after space freed: %v", err)
	}
	if got := waitDone(t, svc, extra); got.State != Completed || got.Result.Ratio != 99 {
		t.Fatalf("resubmitted task: %+v", got)
	}
}

func TestFailedTaskThroughPool(t *testing.T) {
	svc := NewService(fakeProcessor(func(in Input) (Result, error) {
		if in.PointSize == 13 {
			return Result{}, errors.New("unlucky")
		}
		return Result{Ratio: in.PointSize}, nil
	}))
	p := mustPool(t, svc, 2, 4)

	bad := mustCreate(t, svc, input(13))
	good := mustCreate(t, svc, input(14))
	for _, id := range []string{bad.ID, good.ID} {
		if err := p.Submit(id); err != nil {
			t.Fatal(err)
		}
	}
	if got := waitDone(t, svc, bad.ID); got.State != Failed || got.Result != nil || got.Error != "unlucky" {
		t.Errorf("bad: %+v", got)
	}
	if got := waitDone(t, svc, good.ID); got.State != Completed || got.Result.Ratio != 14 {
		t.Errorf("good: %+v", got)
	}
}

// Readers hammer Get/List while workers write results. Every snapshot must be internally consistent (a Completed task always has its own result). The real
// check is the race detector: run with -race.
func TestConcurrentReadsWhileProcessing(t *testing.T) {
	const n = 100
	svc := NewService(echoProcessor)
	p := mustPool(t, svc, 4, n)

	ids := make([]string, n)
	for i := range ids {
		ids[i] = mustCreate(t, svc, input(float64(i+1))).ID
	}

	check := func(tk Task) {
		switch tk.State {
		case Completed:
			if tk.Result == nil || tk.Result.Ratio != tk.Input.PointSize {
				t.Errorf("inconsistent completed task: %+v", tk)
			}
		case Pending, Queued, Processing:
			if tk.Result != nil {
				t.Errorf("task %s has a result while %s", tk.ID, tk.State)
			}
		default:
			t.Errorf("unexpected state %q", tk.State)
		}
	}

	stop := make(chan struct{})
	var readers sync.WaitGroup
	for r := 0; r < 4; r++ {
		readers.Add(1)
		go func() {
			defer readers.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				for _, id := range ids {
					tk, err := svc.Get(id)
					if err != nil {
						t.Errorf("Get: %v", err)
						return
					}
					check(tk)
				}
				for _, tk := range svc.List() {
					check(tk)
				}
			}
		}()
	}

	for _, id := range ids {
		if err := p.Submit(id); err != nil {
			t.Fatalf("Submit: %v", err)
		}
	}
	for _, id := range ids {
		waitDone(t, svc, id)
	}
	close(stop)
	readers.Wait()
}

// Many goroutines race to submit the SAME task: exactly one may succeed, and the task must be processed exactly once.
func TestConcurrentSubmitSameTask(t *testing.T) {
	var calls atomic.Int32
	svc := NewService(fakeProcessor(func(in Input) (Result, error) {
		calls.Add(1)
		return Result{Ratio: in.PointSize}, nil
	}))
	p := mustPool(t, svc, 4, 8)
	tk := mustCreate(t, svc, input(5))

	const n = 20
	var wg sync.WaitGroup
	var ok, invalid atomic.Int32
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch err := p.Submit(tk.ID); {
			case err == nil:
				ok.Add(1)
			case errors.Is(err, ErrInvalidTransition):
				invalid.Add(1)
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()
	waitDone(t, svc, tk.ID)
	p.Stop()

	if ok.Load() != 1 || invalid.Load() != n-1 {
		t.Fatalf("ok=%d invalid=%d, want 1 and %d", ok.Load(), invalid.Load(), n-1)
	}
	if calls.Load() != 1 {
		t.Fatalf("processor ran %d times, want 1", calls.Load())
	}
}

func TestStopWaitsForWorkersToExit(t *testing.T) {
	svc := NewService(echoProcessor)
	p := mustPool(t, svc, 4, 4)
	if p.Running() != 4 {
		t.Fatalf("Running = %d, want 4", p.Running())
	}

	finished := make(chan struct{})
	go func() {
		p.Stop()
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(guard):
		t.Fatal("Stop did not return (workers did not terminate)")
	}
	if p.Running() != 0 {
		t.Errorf("Running = %d after Stop, want 0", p.Running())
	}
	if !p.Stopped() {
		t.Error("Stopped() = false after Stop")
	}
}

// Tasks that were accepted before Stop must all be finished by the time Stop returns, even those still waiting in the queue when Stop was called.
func TestStopFinishesAcceptedTasks(t *testing.T) {
	const workers, n = 2, 8
	g := newGate()
	svc := NewService(g)
	p := mustPool(t, svc, workers, n)
	t.Cleanup(g.open)

	ids := make([]string, n)
	for i := range ids {
		ids[i] = mustCreate(t, svc, input(float64(i+1))).ID
		if err := p.Submit(ids[i]); err != nil {
			t.Fatal(err)
		}
	}
	g.waitStarted(t, workers) // 2 running, the rest (6) are queued

	stopped := make(chan struct{})
	go func() {
		p.Stop()
		close(stopped)
	}()
	// Wait until Stop has really begun (queue closed) before releasing the workers, so the queued tasks are drained AFTER shutdown started.
	deadline := time.Now().Add(guard)
	for !p.Stopped() {
		if time.Now().After(deadline) {
			t.Fatal("Stop never started")
		}
		runtime.Gosched()
	}
	g.open()

	select {
	case <-stopped:
	case <-time.After(guard):
		t.Fatal("Stop did not return")
	}
	for i, id := range ids {
		got := mustGet(t, svc, id)
		if got.State != Completed || got.Result == nil || got.Result.Ratio != float64(i+1) {
			t.Errorf("accepted task %s was lost or wrong after Stop: %+v", id, got)
		}
	}
}

func TestSubmitAfterStopIsRejected(t *testing.T) {
	svc := NewService(echoProcessor)
	p := mustPool(t, svc, 2, 4)
	tk := mustCreate(t, svc, input(12))

	p.Stop()

	if err := p.Submit(tk.ID); !errors.Is(err, ErrPoolStopped) {
		t.Fatalf("want ErrPoolStopped, got %v", err)
	}
	if got := mustGet(t, svc, tk.ID); got.State != Pending {
		t.Fatalf("task must stay Pending, got %s", got.State)
	}
}

func TestStopTwiceAndConcurrently(t *testing.T) {
	svc := NewService(echoProcessor)
	p := mustPool(t, svc, 3, 4)

	p.Stop()
	p.Stop() //second call: no panic (no double close), no hang

	p2 := mustPool(t, NewService(echoProcessor), 3, 4)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p2.Stop()
		}()
	}
	wg.Wait()
	if p2.Running() != 0 {
		t.Fatalf("Running = %d, want 0", p2.Running())
	}
}

// Submitters and Stop run at the same time. There must be no panic (send on a closed channel), no deadlock, and every task whose Submit returned nil must
// be Completed after Stop. Everything else must still be Pending.
func TestConcurrentSubmitAndStop(t *testing.T) {
	const submitters, perSubmitter = 8, 50
	svc := NewService(echoProcessor)
	p := mustPool(t, svc, 4, 16) // small queue on purpose: ErrQueueFull is expected too

	ids := make([][]string, submitters)
	for s := range ids {
		for i := 0; i < perSubmitter; i++ {
			ids[s] = append(ids[s], mustCreate(t, svc, input(float64(s*perSubmitter+i+1))).ID)
		}
	}

	var (
		mu       sync.Mutex
		accepted = map[string]bool{}
		wg       sync.WaitGroup
		started  = make(chan struct{}, submitters)
	)
	for s := 0; s < submitters; s++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, id := range ids[s] {
				err := p.Submit(id)
				if i == 0 {
					started <- struct{}{}
				}
				switch {
				case err == nil:
					mu.Lock()
					accepted[id] = true
					mu.Unlock()
				case errors.Is(err, ErrQueueFull), errors.Is(err, ErrPoolStopped):
				default:
					t.Errorf("unexpected error: %v", err)
				}
			}
		}()
	}

	// Stop while the submitters are still going.
	for i := 0; i < submitters/2; i++ {
		<-started
	}
	p.Stop()
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	for s := range ids {
		for _, id := range ids[s] {
			got := mustGet(t, svc, id)
			if accepted[id] {
				if got.State != Completed || got.Result == nil || got.Result.Ratio != got.Input.PointSize {
					t.Errorf("accepted task %s lost or wrong: %+v", id, got)
				}
			} else if got.State != Pending {
				t.Errorf("rejected task %s is %s, want Pending", id, got.State)
			}
		}
	}
	if p.Running() != 0 {
		t.Errorf("Running = %d after Stop", p.Running())
	}
}
