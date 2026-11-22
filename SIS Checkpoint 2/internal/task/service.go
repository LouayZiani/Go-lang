package task

import (
	"fmt"
	"math"
	"sync"

	"contrastchecker/internal/contrast"
)

type record struct {
	task Task
	done chan struct{}
}

type Service struct {
	processor Processor

	mu      sync.RWMutex
	records map[string]*record
	order   []string
	nextID  uint64
}

func NewService(p Processor) *Service {
	return &Service{
		processor: p,
		records:   make(map[string]*record),
	}
}

func (s *Service) Create(in Input) (Task, error) {
	if err := validate(in); err != nil {
		return Task{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	t := Task{ID: fmt.Sprintf("task-%d", s.nextID), Input: in, State: Pending}
	s.records[t.ID] = &record{task: t, done: make(chan struct{})}
	s.order = append(s.order, t.ID)
	return t.clone(), nil
}

func validate(in Input) error {
	if _, err := contrast.ParseHex(in.Foreground); err != nil {
		return fmt.Errorf("%w: foreground: %v", ErrInvalidInput, err)
	}
	if _, err := contrast.ParseHex(in.Background); err != nil {
		return fmt.Errorf("%w: background: %v", ErrInvalidInput, err)
	}
	if math.IsNaN(in.PointSize) || math.IsInf(in.PointSize, 0) || in.PointSize <= 0 {
		return fmt.Errorf("%w: point size must be a positive finite number, got %v", ErrInvalidInput, in.PointSize)
	}
	return nil
}

func (s *Service) Get(id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rec, ok := s.records[id]
	if !ok {
		return Task{}, fmt.Errorf("%w: %q", ErrTaskNotFound, id)
	}
	return rec.task.clone(), nil
}

func (s *Service) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Task, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, s.records[id].task.clone())
	}
	return out
}

func (s *Service) Done(id string) (<-chan struct{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rec, ok := s.records[id]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrTaskNotFound, id)
	}
	return rec.done, nil
}

func (s *Service) Process(id string) (Task, error) {
	in, err := s.move(id, Pending, Processing)
	if err != nil {
		return Task{}, err
	}
	return s.execute(id, in), nil
}
func (s *Service) move(id string, from, to State) (Input, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.records[id]
	if !ok {
		return Input{}, fmt.Errorf("%w: %q", ErrTaskNotFound, id)
	}
	if rec.task.State != from || !CanTransition(from, to) {
		return Input{}, fmt.Errorf("%w: task %s is %s, cannot go %s -> %s",
			ErrInvalidTransition, id, rec.task.State, from, to)
	}
	rec.task.State = to
	return rec.task.Input, nil
}

func (s *Service) execute(id string, in Input) Task {
	res, err := s.processor.Process(in)
	return s.finish(id, res, err)
}

func (s *Service) runQueued(id string) {
	in, err := s.move(id, Queued, Processing)
	if err != nil {
		return
	}
	s.execute(id, in)
}

func (s *Service) finish(id string, res Result, procErr error) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec := s.records[id]
	if procErr != nil {
		rec.task.State = Failed
		rec.task.Error = procErr.Error()
	} else {
		r := res
		rec.task.State = Completed
		rec.task.Result = &r
	}
	close(rec.done)
	return rec.task.clone()
}
