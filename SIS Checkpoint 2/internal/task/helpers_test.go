package task

import (
	"sync"
	"testing"
	"time"
)

// guard is only a safety net so a broken implementation fails the test instead of hanging it. Passing tests never wait for it.
const guard = 5 * time.Second

// fakeProcessor lets a test decide what processing does
type fakeProcessor func(Input) (Result, error)

func (f fakeProcessor) Process(in Input) (Result, error) { return f(in) }

// echoProcessor puts the point size into Result.Ratio. Every task in a test gets a unique size, so a result that landed in the wrong task is detectable
var echoProcessor = fakeProcessor(func(in Input) (Result, error) {
	return Result{Ratio: in.PointSize}, nil
})

func input(size float64) Input {
	return Input{Foreground: "#000000", Background: "#ffffff", PointSize: size}
}

func mustCreate(t *testing.T, svc *Service, in Input) Task {
	t.Helper()
	tk, err := svc.Create(in)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	return tk
}

func mustPool(t *testing.T, svc *Service, workers, capacity int) *Pool {
	t.Helper()
	p, err := NewPool(svc, workers, capacity)
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	t.Cleanup(p.Stop) //idempotent, so tests may also call Stop themselves
	return p
}

func waitDone(t *testing.T, svc *Service, id string) Task {
	t.Helper()
	done, err := svc.Done(id)
	if err != nil {
		t.Fatalf("Done(%s): %v", id, err)
	}
	select {
	case <-done:
	case <-time.After(guard):
		t.Fatalf("task %s did not finish in time", id)
	}
	tk, err := svc.Get(id)
	if err != nil {
		t.Fatalf("Get(%s): %v", id, err)
	}
	return tk
}

func mustGet(t *testing.T, svc *Service, id string) Task {
	t.Helper()
	tk, err := svc.Get(id)
	if err != nil {
		t.Fatalf("Get(%s): %v", id, err)
	}
	return tk
}

// gate is a processor whose tasks block until the test opens it. That lets a test freeze workers in the middle of a task and inspect the system.
type gate struct {
	started chan struct{} //one signal each time a task enters Process
	release chan struct{}
	once    sync.Once
}

func newGate() *gate {
	return &gate{started: make(chan struct{}, 1024), release: make(chan struct{})}
}

func (g *gate) Process(in Input) (Result, error) {
	g.started <- struct{}{} //buffered, never blocks
	<-g.release
	return Result{Ratio: in.PointSize}, nil
}

// open lets all blocked and future tasks run. Safe to call many times
func (g *gate) open() { g.once.Do(func() { close(g.release) }) }

// waitStarted blocks until n tasks have entered Process
func (g *gate) waitStarted(t *testing.T, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		select {
		case <-g.started:
		case <-time.After(guard):
			t.Fatalf("only %d of %d tasks started processing", i, n)
		}
	}
}
