package fsmincubator

import (
	"github.com/dc0d/fsm"
	"github.com/pkg/errors"
)

//-----------------------------------------------------------------------------

type activationOptions struct {
	final   fsm.State
	onError func(error) fsm.State
}

//-----------------------------------------------------------------------------

// Option for Activate func
type Option func(*activationOptions)

// Final sets the final state that flow will go to, the returned error will be ignored
// and just stops the activation of the final state. This state will not be activated
// if there is an error by initial state.
func Final(final fsm.State) Option {
	return func(opt *activationOptions) {
		opt.final = final
	}
}

// OnError sets the state that flow will go to, even if it panics.
// If not provided, panics will not be recovered and must be handled by caller.
// This shoud not panic and errors from this state will be ignored.
func OnError(onError func(error) fsm.State) Option {
	return func(opt *activationOptions) {
		opt.onError = onError
	}
}

//-----------------------------------------------------------------------------

// Activate activates the state and it's consecutive states until the next state
// is nil or encounters an error
func Activate(initial fsm.State, options ...Option) (funcErr error) {
	// setup options
	opt := new(activationOptions)
	for _, vopt := range options {
		if vopt == nil {
			continue
		}
		vopt(opt)
	}

	// prepare calling scene for onError
	if opt.onError != nil {
		defer func() {
			if funcErr == nil {
				return
			}
			fsm.Activate(opt.onError(funcErr))
		}()

		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					funcErr = err
					return
				}
				funcErr = errors.Errorf("%+v", e)
			}
		}()
	}

	// activating the main state
	funcErr = fsm.Activate(initial)

	// activating final state
	if opt.final != nil {
		fsm.Activate(opt.final)
	}

	return
}

//-----------------------------------------------------------------------------
