package task

import (
	"errors"
	"math"
	"testing"
)

type failingProcessor struct{}

func (failingProcessor) Process(Input) (Result, error) {
	return Result{}, errors.New("boom")
}

var (
	blackOnWhite = Input{Foreground: "#000000", Background: "#FFFFFF", PointSize: 12}
	greyOnWhite  = Input{Foreground: "#777777", Background: "#FFFFFF", PointSize: 12}
)

func newTestService() *Service { return NewService(WCAGProcessor{}) }

func TestCreate_Success(t *testing.T) {
	s := newTestService()

	got, err := s.Create(blackOnWhite)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}
	if got.ID == "" {
		t.Error("Create() returned an empty ID")
	}
	if got.State != Pending {
		t.Errorf("state = %v, want Pending", got.State)
	}
	if got.Result != nil {
		t.Errorf("result = %+v, want nil before processing", got.Result)
	}
	if got.Input != blackOnWhite {
		t.Errorf("input = %+v, want %+v", got.Input, blackOnWhite)
	}
}

func TestCreate_UniqueIDs(t *testing.T) {
	s := newTestService()
	seen := map[string]bool{}

	for i := 0; i < 50; i++ {
		got, err := s.Create(blackOnWhite)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if seen[got.ID] {
			t.Fatalf("duplicate ID %q", got.ID)
		}
		seen[got.ID] = true
	}
}

func TestCreate_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		in   Input
	}{
		{"empty foreground", Input{"", "#FFFFFF", 12, false}},
		{"empty background", Input{"#000000", "", 12, false}},
		{"whitespace-only foreground", Input{"   ", "#FFFFFF", 12, false}},
		{"malformed foreground", Input{"#12", "#FFFFFF", 12, false}},
		{"non-hex background", Input{"#000000", "#ZZZZZZ", 12, false}},
		{"zero size", Input{"#000000", "#FFFFFF", 0, false}},
		{"negative size", Input{"#000000", "#FFFFFF", -4, false}},
		{"NaN size", Input{"#000000", "#FFFFFF", math.NaN(), false}},
		{"infinite size", Input{"#000000", "#FFFFFF", math.Inf(1), false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestService()

			_, err := s.Create(tt.in)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("Create() error = %v, want ErrInvalidInput", err)
			}
			if n := len(s.List()); n != 0 {
				t.Errorf("rejected input was stored: %d tasks in service", n)
			}
		})
	}
}

func TestGet_Success(t *testing.T) {
	s := newTestService()
	created, _ := s.Create(greyOnWhite)

	got, err := s.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if got.ID != created.ID || got.Input != greyOnWhite {
		t.Errorf("Get() = %+v, want the task created with %+v", got, greyOnWhite)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := newTestService()
	s.Create(blackOnWhite)

	_, err := s.Get("task-999")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("Get() error = %v, want ErrTaskNotFound", err)
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	s := newTestService()
	created, _ := s.Create(blackOnWhite)
	s.Process(created.ID)

	got, _ := s.Get(created.ID)
	got.State = Pending
	got.Result.AA = false
	got.Input.Foreground = "#123456"

	again, _ := s.Get(created.ID)
	if again.State != Completed || !again.Result.AA || again.Input.Foreground != "#000000" {
		t.Errorf("changing a returned copy altered the stored task: %+v", again)
	}
}

func TestProcess_Success(t *testing.T) {
	tests := []struct {
		name      string
		in        Input
		wantRatio float64
		wantLarge bool
		wantAA    bool
		wantAAA   bool
		wantLevel string
	}{
		{
			name:      "black on white passes everything",
			in:        Input{"#000000", "#FFFFFF", 12, false},
			wantRatio: 21, wantLarge: false, wantAA: true, wantAAA: true, wantLevel: "AAA",
		},
		{
			name:      "#767676 on white passes AA but not AAA",
			in:        Input{"#767676", "#FFFFFF", 12, false},
			wantRatio: 4.54, wantLarge: false, wantAA: true, wantAAA: false, wantLevel: "AA",
		},
		{
			name:      "#777777 on white fails at body size",
			in:        Input{"#777777", "#FFFFFF", 12, false},
			wantRatio: 4.48, wantLarge: false, wantAA: false, wantAAA: false, wantLevel: "Fail",
		},
		{
			name:      "same grey passes AA once the text is large (18pt)",
			in:        Input{"#777777", "#FFFFFF", 18, false},
			wantRatio: 4.48, wantLarge: true, wantAA: true, wantAAA: false, wantLevel: "AA",
		},
		{
			name:      "14pt bold counts as large text",
			in:        Input{"#777777", "#FFFFFF", 14, true},
			wantRatio: 4.48, wantLarge: true, wantAA: true, wantAAA: false, wantLevel: "AA",
		},
		{
			name:      "14pt regular is not large text",
			in:        Input{"#777777", "#FFFFFF", 14, false},
			wantRatio: 4.48, wantLarge: false, wantAA: false, wantAAA: false, wantLevel: "Fail",
		},
		{
			name:      "pure red on white is 4:1: fails body text",
			in:        Input{"#FF0000", "#FFFFFF", 12, false},
			wantRatio: 4.00, wantLarge: false, wantAA: false, wantAAA: false, wantLevel: "Fail",
		},
		{
			name:      "colours given in short form and swapped still work",
			in:        Input{"#fff", "#000", 12, false},
			wantRatio: 21, wantLarge: false, wantAA: true, wantAAA: true, wantLevel: "AAA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestService()
			created, err := s.Create(tt.in)
			if err != nil {
				t.Fatalf("Create() unexpected error: %v", err)
			}

			got, err := s.Process(created.ID)
			if err != nil {
				t.Fatalf("Process() unexpected error: %v", err)
			}

			if got.State != Completed {
				t.Fatalf("state = %v, want Completed", got.State)
			}
			if got.Result == nil {
				t.Fatal("result is nil after processing")
			}
			if math.Abs(got.Result.Ratio-tt.wantRatio) > 0.005 {
				t.Errorf("ratio = %.4f, want %.2f", got.Result.Ratio, tt.wantRatio)
			}
			if got.Result.LargeText != tt.wantLarge {
				t.Errorf("LargeText = %v, want %v", got.Result.LargeText, tt.wantLarge)
			}
			if got.Result.AA != tt.wantAA {
				t.Errorf("AA = %v, want %v", got.Result.AA, tt.wantAA)
			}
			if got.Result.AAA != tt.wantAAA {
				t.Errorf("AAA = %v, want %v", got.Result.AAA, tt.wantAAA)
			}
			if got.Result.Level() != tt.wantLevel {
				t.Errorf("Level() = %q, want %q", got.Result.Level(), tt.wantLevel)
			}

			stored, _ := s.Get(created.ID)
			if stored.State != Completed || stored.Result == nil || *stored.Result != *got.Result {
				t.Errorf("stored task = %+v, does not match processed task %+v", stored, got)
			}
		})
	}
}

func TestProcess_NotFound(t *testing.T) {
	s := newTestService()

	_, err := s.Process("task-1")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("Process() error = %v, want ErrTaskNotFound", err)
	}
}

func TestProcess_AlreadyCompleted(t *testing.T) {
	s := newTestService()
	created, _ := s.Create(greyOnWhite)
	first, err := s.Process(created.ID)
	if err != nil {
		t.Fatalf("first Process() unexpected error: %v", err)
	}

	_, err = s.Process(created.ID)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("second Process() error = %v, want ErrInvalidTransition", err)
	}

	after, _ := s.Get(created.ID)
	if after.State != Completed {
		t.Errorf("state = %v, want it to stay Completed", after.State)
	}
	if after.Result == nil || *after.Result != *first.Result {
		t.Errorf("result changed after rejected re-process: %+v", after.Result)
	}
}

func TestProcess_ProcessorFailure(t *testing.T) {
	s := NewService(failingProcessor{})
	created, _ := s.Create(blackOnWhite)

	got, err := s.Process(created.ID)
	if err == nil {
		t.Fatal("Process() error = nil, want the processor's error")
	}
	if got.State != Failed {
		t.Errorf("state = %v, want Failed", got.State)
	}
	if got.Result != nil {
		t.Errorf("result = %+v, want nil for a failed task", got.Result)
	}
	if got.Failure == "" {
		t.Error("Failure message is empty")
	}

	if _, err := s.Process(created.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("re-processing a Failed task: error = %v, want ErrInvalidTransition", err)
	}
}

func TestMultipleIndependentTasks(t *testing.T) {
	s := newTestService()

	a, _ := s.Create(Input{"#000000", "#FFFFFF", 12, false}) // 21:1
	b, _ := s.Create(Input{"#777777", "#FFFFFF", 12, false}) // 4.48:1
	c, _ := s.Create(Input{"#FF0000", "#FFFFFF", 12, false}) // 4.00:1

	if _, err := s.Process(a.ID); err != nil {
		t.Fatalf("Process(a): %v", err)
	}
	if _, err := s.Process(c.ID); err != nil {
		t.Fatalf("Process(c): %v", err)
	}

	gotA, _ := s.Get(a.ID)
	gotB, _ := s.Get(b.ID)
	gotC, _ := s.Get(c.ID)

	if gotA.State != Completed || gotA.Result == nil || !gotA.Result.AAA {
		t.Errorf("task a = %+v, want Completed with AAA result", gotA)
	}
	if gotB.State != Pending || gotB.Result != nil {
		t.Errorf("task b = %+v, want untouched Pending without result", gotB)
	}
	if gotC.State != Completed || gotC.Result == nil || gotC.Result.AA {
		t.Errorf("task c = %+v, want Completed with failing result", gotC)
	}
	if gotA.Result.Ratio == gotC.Result.Ratio {
		t.Error("tasks a and c share the same ratio: results were mixed up")
	}

	if gotA.Input.Foreground != "#000000" || gotB.Input.Foreground != "#777777" || gotC.Input.Foreground != "#FF0000" {
		t.Errorf("inputs were mixed up: %v %v %v", gotA.Input, gotB.Input, gotC.Input)
	}
}

func TestList_CreationOrder(t *testing.T) {
	s := newTestService()
	var ids []string
	for i := 0; i < 5; i++ {
		got, _ := s.Create(blackOnWhite)
		ids = append(ids, got.ID)
	}

	list := s.List()
	if len(list) != len(ids) {
		t.Fatalf("List() returned %d tasks, want %d", len(list), len(ids))
	}
	for i, task := range list {
		if task.ID != ids[i] {
			t.Errorf("List()[%d].ID = %s, want %s", i, task.ID, ids[i])
		}
	}
}
