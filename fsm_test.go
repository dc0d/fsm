package fsm

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	errNegative = errors.Errorf("NEGATIVE")
)

type sample struct {
	state int
}

func (s *sample) Start() State {
	return StateFunc(s.dispatcher)
}

func (s *sample) odd() (State, error) {
	s.state++
	return StateFunc(s.dispatcher), nil
}

func (s *sample) even() (State, error) {
	s.state++
	return StateFunc(s.dispatcher), nil
}

func (s *sample) dispatcher() (State, error) {
	if s.state > 8 {
		return nil, nil
	}
	if s.state < 0 {
		return nil, errNegative
	}
	if s.state%2 == 0 {
		return StateFunc(s.even), nil
	}
	return StateFunc(s.odd), nil
}

func Test01(t *testing.T) {
	assert := assert.New(t)

	fsm := &sample{}
	err := Activate(fsm.Start())

	assert.Equal(9, fsm.state)
	assert.Equal(nil, err)

	fsm.state = -10
	err = Activate(fsm.Start())

	assert.Equal(-10, fsm.state)
	assert.Equal(errNegative, err)
}
