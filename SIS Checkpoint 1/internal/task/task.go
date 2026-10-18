package task

import (
	"fmt"
	"math"

	"sis/internal/contrast"
)

type Input struct {
	Foreground string
	Background string
	PointSize  float64
	Bold       bool
}

func (in Input) Validate() error {
	if _, err := contrast.ParseHex(in.Foreground); err != nil {
		return fmt.Errorf("%w: foreground: %v", ErrInvalidInput, err)
	}
	if _, err := contrast.ParseHex(in.Background); err != nil {
		return fmt.Errorf("%w: background: %v", ErrInvalidInput, err)
	}
	if !(in.PointSize > 0) || math.IsInf(in.PointSize, 0) {
		return fmt.Errorf("%w: point size must be a positive number, got %v", ErrInvalidInput, in.PointSize)
	}
	return nil
}

type Result struct {
	Ratio     float64
	LargeText bool
	AA        bool
	AAA       bool
}

func (r Result) Level() string {
	switch {
	case r.AAA:
		return "AAA"
	case r.AA:
		return "AA"
	default:
		return "Fail"
	}
}

type Task struct {
	ID      string
	Input   Input
	State   State
	Result  *Result
	Failure string
}

func (t *Task) transition(next State) error {
	if !t.State.CanTransitionTo(next) {
		return fmt.Errorf("%w: task %s cannot go from %s to %s", ErrInvalidTransition, t.ID, t.State, next)
	}
	t.State = next
	return nil
}

func (t *Task) snapshot() Task {
	c := *t
	if t.Result != nil {
		r := *t.Result
		c.Result = &r
	}
	return c
}
