package fsm

//-----------------------------------------------------------------------------

// State represents a state activity
type State interface {
	Activate() (State, error)
}

//-----------------------------------------------------------------------------

// StateFunc is a function that satisfies the State interface
type StateFunc func() (State, error)

// Activate satisfies the State interface
func (stateFn StateFunc) Activate() (State, error) { return stateFn() }

//-----------------------------------------------------------------------------
