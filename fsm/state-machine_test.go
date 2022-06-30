package fsm

import (
	"testing"
)

const (
	InsertCoin Command = iota
	PushButton
)

const (
	Locked State = iota
	Unlocked
)

func TestNewShouldReturnStateMachineInInitialStatusAndNoTransitions(t *testing.T) {

	sm := New(Locked)
	if sm.State != Locked {
		t.Errorf("Incorrect initial status: Got: %d, expected %d", sm.State, Locked)
	}

	if len(sm.Transitions) != 0 {
		t.Errorf("Incorrect initial Transitions lengh: Got: %d, expected %d", len(sm.Transitions), 0)
	}
}

func TestAddTransitionShouldAddTransitionWhenTransitionDoesNotExist(t *testing.T) {

	sm := New(Locked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)

	if len(sm.Transitions) != 1 {
		t.Errorf("Incorrect Transitions lengh: Got: %d, expected %d", len(sm.Transitions), 1)
	}

	if sm.Transitions[0].From != Locked {
		t.Errorf("Incorrect added transitions From: Got: %d, expected %d", sm.Transitions[0].From, Locked)
	}

	if sm.Transitions[0].To != Unlocked {
		t.Errorf("Incorrect added transitions To: Got: %d, expected %d", sm.Transitions[0].To, Unlocked)
	}

	if sm.Transitions[0].Command != InsertCoin {
		t.Errorf("Incorrect added transitions Command: Got: %d, expected %d", sm.Transitions[0].Command, InsertCoin)
	}
}

func TestAddTransitionShouldNotAddTransitionWhenTransitionExists(t *testing.T) {

	sm := New(Locked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)

	if len(sm.Transitions) != 1 {
		t.Errorf("Incorrect Transitions lengh: Got: %d, expected %d", len(sm.Transitions), 1)
	}
}

func TestExecuteCommandWhenTransitionExistsShouldExecute(t *testing.T) {
	sm := New(Locked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)

	err := sm.ExecuteCommand(InsertCoin)

	if err != nil {
		t.Errorf("Command not executed")
	}

	if sm.State != Unlocked {
		t.Errorf("Incorrect status : Got: %d, expected %d", sm.State, Unlocked)
	}
}

func TestExecuteCommandWhenTransitionDoesNotExistShouldReturnCommandNotAvailableError(t *testing.T) {
	sm := New(Locked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)

	err := sm.ExecuteCommand(PushButton)

	if err == nil {
		t.Errorf("Expected to receive error of type")
		return
	}

	if err.Command != PushButton {
		t.Errorf("Unexpected error Command: Got %d, expected %d", err.Command, PushButton)
	}

	if err.FromState != Locked {
		t.Errorf("Unexpected error FromState: Got %d, expected %d", err.FromState, Locked)
	}

}

func TestExecuteCommandWhenTransitionsAreEmptyShouldNotExecute(t *testing.T) {
	sm := New(Locked)

	err := sm.ExecuteCommand(PushButton)

	if err == nil {
		t.Errorf("Command executed")
	}

	if sm.State != Locked {
		t.Errorf("Incorrect status : Got: %d, expected %d", sm.State, Locked)
	}
}
