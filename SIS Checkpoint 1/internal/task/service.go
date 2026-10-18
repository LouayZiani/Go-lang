package task

import "fmt"

type Service struct {
	processor Processor
	tasks     map[string]*Task
	order     []string
	nextID    int
}

func NewService(p Processor) *Service {
	return &Service{
		processor: p,
		tasks:     make(map[string]*Task),
	}
}

func (s *Service) Create(in Input) (Task, error) {
	if err := in.Validate(); err != nil {
		return Task{}, err
	}

	s.nextID++
	t := &Task{
		ID:    fmt.Sprintf("task-%d", s.nextID),
		Input: in,
		State: Pending,
	}
	s.tasks[t.ID] = t
	s.order = append(s.order, t.ID)
	return t.snapshot(), nil
}

func (s *Service) Get(id string) (Task, error) {
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("%w: %q", ErrTaskNotFound, id)
	}
	return t.snapshot(), nil
}

func (s *Service) List() []Task {
	out := make([]Task, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, s.tasks[id].snapshot())
	}
	return out
}

func (s *Service) Process(id string) (Task, error) {
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("%w: %q", ErrTaskNotFound, id)
	}

	if err := t.transition(Processing); err != nil {
		return Task{}, err
	}

	res, err := s.processor.Process(t.Input)
	if err != nil {
		if terr := t.transition(Failed); terr != nil {
			return Task{}, terr
		}
		t.Failure = err.Error()
		return t.snapshot(), fmt.Errorf("processing task %s: %w", id, err)
	}

	t.Result = &res
	if err := t.transition(Completed); err != nil {
		return Task{}, err
	}
	return t.snapshot(), nil
}
