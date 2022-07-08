package fsm

import (
	"fmt"
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

	expectedTransitions := map[Command]map[State]State{
		InsertCoin: {
			Locked: Unlocked,
		},
	}

	if expected, got := fmt.Sprint(expectedTransitions), fmt.Sprint(sm.Transitions); expected != got {
		t.Errorf("Unexpected transitions. Expected: %v but got: %v", expected, got)
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
