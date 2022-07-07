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
	if expected, got := Locked, sm.State; expected != got {
		t.Errorf("Incorrect initial status: Got: %v, expected %v", sm.State, Locked)
	}

	if expected, got := 0, len(sm.Transitions); expected != got {
		t.Errorf("Incorrect initial Transitions lengh: Got: %v, expected %v", len(sm.Transitions), 0)
	}
}

func TestAddTransitionShouldAddTransitionWhenTransitionDoesNotExist(t *testing.T) {

	sm := New(Locked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)

	if expected, got := 1, len(sm.Transitions); expected != got {
		t.Errorf("Incorrect Transitions lengh: Got: %v, expected %v", len(sm.Transitions), 1)
	}

	if expected, got := Locked, sm.Transitions[0].From; expected != got {
		t.Errorf("Incorrect added transitions From: Got: %v, expected %v", got, expected)
	}

	if expected, got := Unlocked, sm.Transitions[0].To; expected != got {
		t.Errorf("Incorrect added transitions To: Got: %v, expected %v", got, expected)
	}

	if expected, got := InsertCoin, sm.Transitions[0].Command; expected != got {
		t.Errorf("Incorrect added transitions Command: Got: %v, expected %v", got, expected)
	}
}

func TestAddTransitionShouldNotAddTransitionWhenTransitionExists(t *testing.T) {

	sm := New(Locked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)

	if expected, got := 1, len(sm.Transitions); expected != got {
		t.Errorf("Incorrect Transitions lengh: Got: %v, expected %v", len(sm.Transitions), 1)
	}
}

func TestExecuteCommandWhenTransitionExistsShouldExecute(t *testing.T) {
	sm := New(Locked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)

	_, err := sm.ExecuteCommand(InsertCoin)

	if err != nil {
		t.Errorf("Command not executed")
	}

	if expected, got := Unlocked, sm.State; expected != got {
		t.Errorf("Incorrect status : Got: %v, expected %v", sm.State, Unlocked)
	}
}

func TestExecuteCommandWhenTransitionDoesNotExistShouldReturnCommandNotAvailableError(t *testing.T) {
	sm := New(Locked)
	sm.AddTransition(Locked, InsertCoin, Unlocked)
	sm.AddTransition(Unlocked, PushButton, Locked)

	_, err := sm.ExecuteCommand(PushButton)

	if err == nil {
		t.Errorf("Expected to receive error")
		return
	}
}

func TestExecuteCommandWhenTransitionsAreEmptyShouldNotExecute(t *testing.T) {
	sm := New(Locked)

	_, err := sm.ExecuteCommand(PushButton)

	if err == nil {
		t.Errorf("Command executed")
	}

	if expected, got := Locked, sm.State; expected != got {
		t.Errorf("Incorrect status : Got: %v, expected %v", sm.State, Locked)
	}
}
