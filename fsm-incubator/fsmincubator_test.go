package fsmincubator

import (
	"testing"

	"github.com/dc0d/errgo/sentinel"
	"github.com/dc0d/fsm"
	"github.com/stretchr/testify/assert"
)

var (
	errInvalid = sentinel.Errorf("INVALID")
	errPanic   = sentinel.Errorf("PANIC")
)

type sample struct {
	state    int
	previous fsm.State
}

func (s *sample) start() (next fsm.State, funcErr error) {
	if s.previous != nil {
		next = s.previous
		return
	}
	next = fsm.StateFunc(s.validate)
	return
}

func (s *sample) validate() (next fsm.State, funcErr error) {
	if s.state == -100 {
		panic(errInvalid)
	}
	if s.state < 0 {
		return
	}
	if s.state%23 != 0 {
		funcErr = errInvalid
		return
	}
	next = fsm.StateFunc(s.increase)
	return
}

func (s *sample) increase() (next fsm.State, funcErr error) {
	s.state++
	next = fsm.StateFunc(s.validate)
	return
}

func (s *sample) final() (next fsm.State, funcErr error) {
	s.state = -1
	return
}

func (s *sample) onError(err error) fsm.State {
	return fsm.StateFunc(func() (next fsm.State, funcErr error) {
		// err handled
		s.state = -2
		return
	})
}

var panicFlag int

func justPanic() (next fsm.State, funcErr error) {
	panicFlag = 100
	panic(errPanic)
}

func TestActivate(t *testing.T) {
	assert := assert.New(t)

	m := new(sample)
	err := Activate(fsm.StateFunc(m.start))
	assert.Equal(errInvalid, err)
}

func TestActivate2(t *testing.T) {
	assert := assert.New(t)

	m := new(sample)
	m.state = -1000
	err := Activate(fsm.StateFunc(m.start))
	assert.NoError(err)
}

func TestFinal(t *testing.T) {
	assert := assert.New(t)

	m := new(sample)
	err := Activate(fsm.StateFunc(m.start), Final(fsm.StateFunc(m.final)))
	assert.Equal(errInvalid, err)
	assert.Equal(-1, m.state)
}

func TestOnError(t *testing.T) {
	assert := assert.New(t)

	m := new(sample)
	err := Activate(
		fsm.StateFunc(m.start),
		Final(fsm.StateFunc(m.final)),
		OnError(m.onError))
	assert.Equal(errInvalid, err)
	assert.Equal(-2, m.state)
}

func TestOnErrorWithPanic(t *testing.T) {
	assert := assert.New(t)

	m := new(sample)
	m.state = -100
	err := Activate(
		fsm.StateFunc(m.start),
		Final(fsm.StateFunc(m.final)),
		OnError(m.onError))
	assert.Equal(errInvalid, err)
	assert.Equal(-2, m.state)
}

func TestOnErrorWithPanic2(t *testing.T) {
	assert := assert.New(t)

	m := new(sample)
	m.previous = fsm.StateFunc(justPanic)
	err := Activate(
		fsm.StateFunc(m.start),
		OnError(m.onError))
	assert.Error(err)
	assert.Equal(errPanic, err)
	assert.Equal(-2, m.state)
	assert.Equal(100, panicFlag)
}

func ExampleActivate() {
	// in this sample a struct is used for preserving the state
	// but it's not a requirement since fsm.State is just an interface
	// and it's even possible to go to states from another struct, etc, etc.
	// Final and OnError states are optional.
	m := new(sample)
	err := Activate(
		fsm.StateFunc(m.start),
		Final(fsm.StateFunc(m.final)),
		OnError(m.onError))
	// handle err
	_ = err
}
