package task

import (
	"errors"
	"math"
	"sync"
	"testing"
)

func TestCreateValidation(t *testing.T) {
	tests := []struct {
		name    string
		in      Input
		wantErr bool
	}{
		{"valid long", Input{"#112233", "#ffffff", 12, false}, false},
		{"valid short", Input{"#123", "#fff", 12, true}, false},
		{"empty foreground", Input{"", "#fff", 12, false}, true},
		{"empty background", Input{"#fff", "", 12, false}, true},
		{"missing hash", Input{"ffffff", "#000", 12, false}, true},
		{"wrong length", Input{"#ffff", "#000", 12, false}, true},
		{"non hex digit", Input{"#gggggg", "#000", 12, false}, true},
		{"zero size", Input{"#fff", "#000", 0, false}, true},
		{"negative size", Input{"#fff", "#000", -1, false}, true},
		{"NaN size", Input{"#fff", "#000", math.NaN(), false}, true},
		{"+Inf size", Input{"#fff", "#000", math.Inf(1), false}, true},
		{"-Inf size", Input{"#fff", "#000", math.Inf(-1), false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(WCAGProcessor{})
			tk, err := svc.Create(tt.in)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidInput) {
					t.Fatalf("want ErrInvalidInput, got %v", err)
				}
				if len(svc.List()) != 0 {
					t.Error("rejected input must not be stored")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tk.ID == "" || tk.State != Pending || tk.Result != nil {
				t.Errorf("unexpected new task: %+v", tk)
			}
		})
	}
}

func TestUniqueIDsAndListOrder(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	var ids []string
	for i := 1; i <= 5; i++ {
		ids = append(ids, mustCreate(t, svc, input(float64(i))).ID)
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate ID %s", id)
		}
		seen[id] = true
	}
	list := svc.List()
	if len(list) != len(ids) {
		t.Fatalf("List len = %d, want %d", len(list), len(ids))
	}
	for i, tk := range list {
		if tk.ID != ids[i] {
			t.Errorf("List[%d] = %s, want %s", i, tk.ID, ids[i])
		}
	}
}

func TestGetNotFound(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	if _, err := svc.Get("nope"); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("want ErrTaskNotFound, got %v", err)
	}
	if _, err := svc.Done("nope"); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("Done: want ErrTaskNotFound, got %v", err)
	}
}

func TestProcessSync(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	tk := mustCreate(t, svc, Input{"#000", "#fff", 12, false})

	got, err := svc.Process(tk.ID)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if got.State != Completed || got.Result == nil || got.Result.Level != LevelAAA {
		t.Fatalf("unexpected result: %+v", got)
	}
	stored := mustGet(t, svc, tk.ID)
	if stored.State != Completed || stored.Result == nil {
		t.Fatalf("stored task not updated: %+v", stored)
	}
}

func TestProcessTwiceIsRejected(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	tk := mustCreate(t, svc, input(12))
	if _, err := svc.Process(tk.ID); err != nil {
		t.Fatal(err)
	}
	before := mustGet(t, svc, tk.ID)
	if _, err := svc.Process(tk.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("want ErrInvalidTransition, got %v", err)
	}
	after := mustGet(t, svc, tk.ID)
	if before.State != after.State || *before.Result != *after.Result {
		t.Error("rejected Process must leave the task unchanged")
	}
}

func TestProcessUnknown(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	if _, err := svc.Process("nope"); !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("want ErrTaskNotFound, got %v", err)
	}
}

func TestProcessorFailureMarksFailed(t *testing.T) {
	boom := errors.New("boom")
	svc := NewService(fakeProcessor(func(Input) (Result, error) { return Result{}, boom }))
	tk := mustCreate(t, svc, input(12))

	got, err := svc.Process(tk.ID)
	if err != nil {
		t.Fatalf("Process itself should succeed, failure lives on the task: %v", err)
	}
	if got.State != Failed || got.Result != nil || got.Error != "boom" {
		t.Fatalf("unexpected failed task: %+v", got)
	}
	// Failed is terminal.
	if _, err := svc.Process(tk.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("want ErrInvalidTransition, got %v", err)
	}
}

func TestSnapshotsAreCopies(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	tk := mustCreate(t, svc, input(12))
	if _, err := svc.Process(tk.ID); err != nil {
		t.Fatal(err)
	}

	got := mustGet(t, svc, tk.ID)
	got.State = Pending
	got.Result.Ratio = 999
	got.Input.PointSize = 1

	again := mustGet(t, svc, tk.ID)
	if again.State != Completed || again.Result.Ratio == 999 || again.Input.PointSize != 12 {
		t.Fatalf("stored task was modified through a snapshot: %+v", again)
	}

	list := svc.List()
	list[0].Result.Ratio = 555
	if mustGet(t, svc, tk.ID).Result.Ratio == 555 {
		t.Fatal("List leaked a pointer to the stored Result")
	}
}

func TestCanTransition(t *testing.T) {
	tests := []struct {
		from, to State
		want     bool
	}{
		{Pending, Queued, true},
		{Pending, Processing, true},
		{Queued, Processing, true},
		{Queued, Pending, true},
		{Processing, Completed, true},
		{Processing, Failed, true},
		{Pending, Completed, false},
		{Queued, Completed, false},
		{Processing, Pending, false},
		{Processing, Queued, false},
		{Completed, Processing, false},
		{Completed, Pending, false},
		{Failed, Processing, false},
		{Failed, Completed, false},
	}
	for _, tt := range tests {
		if got := CanTransition(tt.from, tt.to); got != tt.want {
			t.Errorf("CanTransition(%s, %s) = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestSyncProcessRejectsQueuedTask(t *testing.T) {
	svc := NewService(WCAGProcessor{})
	tk := mustCreate(t, svc, input(12))
	if _, err := svc.move(tk.ID, Pending, Queued); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Process(tk.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("want ErrInvalidTransition, got %v", err)
	}
	if got := mustGet(t, svc, tk.ID); got.State != Queued {
		t.Fatalf("state changed to %s", got.State)
	}
}

func TestConcurrentCreate(t *testing.T) {
	const n = 100
	svc := NewService(WCAGProcessor{})

	ids := make([]string, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tk, err := svc.Create(input(float64(i + 1)))
			if err != nil {
				t.Errorf("Create: %v", err)
				return
			}
			ids[i] = tk.ID
		}()
	}
	wg.Wait()

	seen := make(map[string]bool, n)
	for _, id := range ids {
		if id == "" || seen[id] {
			t.Fatalf("missing or duplicate ID %q", id)
		}
		seen[id] = true
	}
	if got := len(svc.List()); got != n {
		t.Fatalf("List len = %d, want %d", got, n)
	}
}

func TestConcurrentProcessSameTaskOnlyOneWins(t *testing.T) {
	svc := NewService(echoProcessor)
	tk := mustCreate(t, svc, input(7))

	const n = 20
	var wg sync.WaitGroup
	var mu sync.Mutex
	wins, rejects := 0, 0
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.Process(tk.ID)
			mu.Lock()
			defer mu.Unlock()
			switch {
			case err == nil:
				wins++
			case errors.Is(err, ErrInvalidTransition):
				rejects++
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	if wins != 1 || rejects != n-1 {
		t.Fatalf("wins=%d rejects=%d, want 1 and %d", wins, rejects, n-1)
	}
	if got := mustGet(t, svc, tk.ID); got.State != Completed || got.Result.Ratio != 7 {
		t.Fatalf("unexpected final task: %+v", got)
	}
}
