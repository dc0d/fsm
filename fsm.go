package fsm

//-----------------------------------------------------------------------------

// State represents a state activity
type State interface {
	Activate() (State, error)
}

// StateFunc is a function that satisfies the State interface
type StateFunc func() (State, error)

// Activate satisfies the State interface
func (stateFn StateFunc) Activate() (State, error) { return stateFn() }

//-----------------------------------------------------------------------------

// Activate activates the state and it's consecutive states until the next state
// is nil or encounters an error
func Activate(s State) (funcErr error) {
	next := s
	for next != nil && funcErr == nil {
		next, funcErr = next.Activate()
	}
	return
}

//-----------------------------------------------------------------------------
