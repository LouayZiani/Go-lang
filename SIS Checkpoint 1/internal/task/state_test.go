package task

import "testing"

func TestStateCanTransitionTo(t *testing.T) {
	tests := []struct {
		from, to State
		want     bool
	}{
		{Pending, Processing, true},
		{Processing, Completed, true},
		{Processing, Failed, true},

		{Pending, Completed, false}, // cant skip Processing
		{Pending, Failed, false},
		{Processing, Pending, false}, // cant go back
		{Completed, Processing, false},
		{Completed, Pending, false},
		{Failed, Processing, false},
		{Failed, Completed, false},
		{Pending, Pending, false},
	}

	for _, tt := range tests {
		t.Run(tt.from.String()+"->"+tt.to.String(), func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("%v.CanTransitionTo(%v) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestStateString(t *testing.T) {
	if got := State(99).String(); got != "Unknown" {
		t.Errorf("State(99).String() = %q, want %q", got, "Unknown")
	}
}
