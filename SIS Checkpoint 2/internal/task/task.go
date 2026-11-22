package task

import (
	"errors"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrQueueFull         = errors.New("queue is full")
	ErrPoolStopped       = errors.New("worker pool is stopped")
	ErrInvalidConfig     = errors.New("invalid pool configuration")
)

type State string

const (
	Pending    State = "Pending"
	Queued     State = "Queued"
	Processing State = "Processing"
	Completed  State = "Completed"
	Failed     State = "Failed"
)

var transitions = map[State][]State{
	Pending:    {Queued, Processing},
	Queued:     {Processing, Pending},
	Processing: {Completed, Failed},
}

func CanTransition(from, to State) bool {
	for _, next := range transitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

type Level string

const (
	LevelAAA  Level = "AAA"
	LevelAA   Level = "AA"
	LevelFail Level = "Fail"
)

type Input struct {
	Foreground string
	Background string
	PointSize  float64
	Bold       bool
}

type Result struct {
	Ratio     float64 // 1..21
	LargeText bool
	PassAA    bool
	PassAAA   bool
	Level     Level
}

type Task struct {
	ID     string
	Input  Input
	State  State
	Result *Result
	Error  string
}

// clone returns a deep copy (Result is a pointer)
func (t Task) clone() Task {
	if t.Result != nil {
		r := *t.Result
		t.Result = &r
	}
	return t
}
