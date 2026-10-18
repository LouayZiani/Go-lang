package task

// State is a step in a task's lifecycle.
type State int

const (
	Pending    State = iota // created but not processed yet
	Processing              // processor's running
	Completed               // finished, result's available
	Failed                  // processor returned an err
)

func (s State) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Processing:
		return "Processing"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

var allowed = map[State][]State{
	Pending:    {Processing},
	Processing: {Completed, Failed},
}

func (s State) CanTransitionTo(next State) bool {
	for _, n := range allowed[s] {
		if n == next {
			return true
		}
	}
	return false
}
