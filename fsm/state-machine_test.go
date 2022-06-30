package fsm

import (
	"testing"

	"github.com/cgxarrie/fsm-go/fsm/commands"
	"github.com/cgxarrie/fsm-go/fsm/states"
)

func TestNewShouldReturnStateMachineInInitialStatusAndNoTransitions(t *testing.T) {

	sm := New(states.Locked)
	if sm.State != states.Locked {
		t.Errorf("Incorrect initial status: Got: %d, expected %d", sm.State, states.Locked)
	}

	if len(sm.Transitions) != 0 {
		t.Errorf("Incorrect initial Transitions lengh: Got: %d, expected %d", len(sm.Transitions), 0)
	}
}

func TestAddTransitionWhenTransitionDoesNotExistsShouldAddTransition(t *testing.T) {

	sm := New(states.Locked)
	sm.AddTransition(states.Locked, commands.InsertCoin, states.Unlocked)

	if len(sm.Transitions) != 1 {
		t.Errorf("Incorrect Transitions lengh: Got: %d, expected %d", len(sm.Transitions), 1)
	}

	if sm.Transitions[0].From != states.Locked {
		t.Errorf("Incorrect added transitions From: Got: %d, expected %d", sm.Transitions[0].From, states.Locked)
	}

	if sm.Transitions[0].To != states.Unlocked {
		t.Errorf("Incorrect added transitions To: Got: %d, expected %d", sm.Transitions[0].To, states.Unlocked)
	}

	if sm.Transitions[0].Command != commands.InsertCoin {
		t.Errorf("Incorrect added transitions Command: Got: %d, expected %d", sm.Transitions[0].Command, commands.InsertCoin)
	}
}

func TestAddTransitionWhenTransitionDoesExistsShouldNotAddTransition(t *testing.T) {

	sm := New(states.Locked)
	sm.AddTransition(states.Locked, commands.InsertCoin, states.Unlocked)
	sm.AddTransition(states.Locked, commands.InsertCoin, states.Unlocked)

	if len(sm.Transitions) != 1 {
		t.Errorf("Incorrect Transitions lengh: Got: %d, expected %d", len(sm.Transitions), 1)
	}
}

func TestExecuteCommandShouldExecuteWhenTransitionExists(t *testing.T) {
	sm := New(states.Locked)
	sm.AddTransition(states.Locked, commands.InsertCoin, states.Unlocked)

	executed := sm.ExecuteCommand(commands.InsertCoin)

	if !executed {
		t.Errorf("Command not executed")
	}

	if sm.State != states.Unlocked {
		t.Errorf("Incorrect status : Got: %d, expected %d", sm.State, states.Unlocked)
	}
}

func TestExecuteCommandShouldNotExecuteWhenTransitionDoesNotExists(t *testing.T) {
	sm := New(states.Locked)
	sm.AddTransition(states.Locked, commands.InsertCoin, states.Unlocked)

	executed := sm.ExecuteCommand(commands.PushButton)

	if executed {
		t.Errorf("Command not executed")
	}

	if sm.State != states.Locked {
		t.Errorf("Incorrect status : Got: %d, expected %d", sm.State, states.Locked)
	}
}

func TestExecuteCommandShouldNotExecuteWhenTransitionsAreEmpty(t *testing.T) {
	sm := New(states.Locked)

	executed := sm.ExecuteCommand(commands.PushButton)

	if executed {
		t.Errorf("Command not executed")
	}

	if sm.State != states.Locked {
		t.Errorf("Incorrect status : Got: %d, expected %d", sm.State, states.Locked)
	}
}
