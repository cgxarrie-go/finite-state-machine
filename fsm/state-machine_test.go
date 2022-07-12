package fsm

import (
	"fmt"
	"testing"
)

const (
	insertCoin Command = iota
	pushButton
)

var commands = map[Command]string{
	insertCoin: "InsertCoin",
	pushButton: "PushButton",
}

const (
	locked State = iota
	unlocked
)

var states = map[State]string{
	locked:   "Locked",
	unlocked: "Unlocked",
}

func TestNewShouldReturnStateMachineInInitialStatusAndNoTransitions(t *testing.T) {

	sm := New(locked)
	if expected, got := locked, sm.State; expected != got {
		t.Errorf("Incorrect initial status: Got: %v, expected %v", states[got], states[expected])
	}

	if expected, got := 0, len(sm.Transitions); expected != got {
		t.Errorf("Incorrect initial Transitions lengh: Got: %v, expected %v", got, expected)
	}
}

func TestAddTransitionShouldAddTransitionWhenTransitionDoesNotExist(t *testing.T) {

	sm := New(locked)
	sm.WithTransition().On(insertCoin).From(locked).To(unlocked).Add()

	if expected, got := 1, len(sm.Transitions); expected != got {
		t.Errorf("Incorrect Transitions lengh: Got: %v, expected %v", len(sm.Transitions), 1)
	}

	expectedTransitions := map[Command]map[State]State{
		insertCoin: {
			locked: unlocked,
		},
	}

	if expected, got := fmt.Sprint(expectedTransitions), fmt.Sprint(sm.Transitions); expected != got {
		t.Errorf("Unexpected transitions. Expected: %v but got: %v", expected, got)
	}
}

func TestAddTransitionShouldNotAddTransitionWhenTransitionExists(t *testing.T) {

	sm := New(locked)
	sm.
		WithTransition().On(insertCoin).From(locked).To(unlocked).Add().
		WithTransition().On(insertCoin).From(locked).To(unlocked).Add()

	if expected, got := 1, len(sm.Transitions); expected != got {
		t.Errorf("Incorrect Transitions lengh: Got: %v, expected %v", len(sm.Transitions), 1)
	}
}

func TestExecuteCommandWhenTransitionExistsShouldExecute(t *testing.T) {
	sm := New(locked)
	sm.WithTransition().On(insertCoin).From(locked).To(unlocked).Add()

	_, err := sm.ExecuteCommand(insertCoin)

	if err != nil {
		t.Errorf("Command not executed")
	}

	if expected, got := unlocked, sm.State; expected != got {
		t.Errorf("Incorrect status : Got: %v, expected %v", states[got], states[expected])
	}
}

func TestExecuteCommandWhenTransitionDoesNotExistShouldReturnCommandNotAvailableError(t *testing.T) {
	sm := New(locked)
	sm.WithTransition().On(insertCoin).From(locked).To(unlocked).Add().
		WithTransition().On(pushButton).From(unlocked).To(locked).Add()

	_, err := sm.ExecuteCommand(pushButton)

	if err == nil {
		t.Errorf("Expected to receive error")
		return
	}
}

func TestExecuteCommandWhenTransitionsAreEmptyShouldNotExecute(t *testing.T) {
	sm := New(locked)

	_, err := sm.ExecuteCommand(pushButton)

	if err == nil {
		t.Errorf("Command executed")
	}

	if expected, got := locked, sm.State; expected != got {
		t.Errorf("Incorrect status : Got: %v, expected %v", states[got], states[expected])
	}
}
