package fsm

import (
	"errors"
	"testing"
)

func TestStateMachine(t *testing.T) {
	// Define some states and commands
	const (
		StateA State = iota
		StateB
		StateC
		CommandX Command = iota
		CommandY
		CommandZ
	)

	// Define some test actions
	action1 := func() error { return nil }
	action2 := func() error { return errors.New("error") }

	// Create a new state machine with an initial state of StateA
	fsm := New(StateA)

	// Add some actions to the state machine
	fsm.WithAction(CommandX, action1).
		WithTransition(StateA, StateB, func() bool { return true }).
		WithTransition(StateB, StateC, func() bool { return true })

	fsm.WithAction(CommandY, action2).
		WithTransition(StateA, StateC, func() bool { return true })

	// Test that executing a command with valid transitions returns true and no error
	ok, err := fsm.ExecuteCommand(CommandX)
	if !ok || err != nil {
		t.Errorf("Expected ExecuteCommand(CommandX) to return true and no error, but got %v and %v", ok, err)
	}

	// Test that executing a command with invalid transitions returns false and no error
	ok, err = fsm.ExecuteCommand(CommandY)
	if ok || err != nil {
		t.Errorf("Expected ExecuteCommand(CommandY) to return false and an error, but got %v and %v", ok, err)
	}

	// Test that executing a command with no transitions returns false and no error
	ok, err = fsm.ExecuteCommand(CommandZ)
	if ok || err != nil {
		t.Errorf("Expected ExecuteCommand(CommandZ) to return false and no error, but got %v and %v", ok, err)
	}
}
