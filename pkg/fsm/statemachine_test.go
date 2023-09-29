package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testElement struct {
	state State
}

func (e *testElement) SetState(state State) {
	e.state = state
}

func (e *testElement) State() State {
	return e.state
}

func (e *testElement) Command1() error {
	return nil
}

func (e *testElement) Command2() error {
	return fmt.Errorf("command2 error")
}

func Test_Action_WithTransition_AddFirstTransition_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	action := Action{}

	// Act
	action.WithTransition(state1, state2)

	// Assert
	assert.Len(t, action.transitions, 1)
	assert.Equal(t, state1, action.transitions[state1].From)
	assert.Len(t, action.transitions[state1].Targets, 1)
	assert.Equal(t, state2, action.transitions[state1].Targets[state2].To)
	assert.Equal(t, "", action.transitions[state1].Targets[state2].funcName)
}

func Test_Action_WithConditionedTransition_AddFirstTransition_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	action := Action{}

	// Act
	action.WithConditionedTransition(state1, state2, "condition-func")

	// Assert
	assert.Len(t, action.transitions, 1)
	assert.Equal(t, state1, action.transitions[state1].From)
	assert.Len(t, action.transitions[state1].Targets, 1)
	assert.Equal(t, state2, action.transitions[state1].Targets[state2].To)
	assert.Equal(t, "condition-func", action.transitions[state1].Targets[state2].funcName)
}

func Test_Action_WithTransition_AddTargetToExistingTransition_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	action := Action{
		transitions: map[State]*Transition{
			state1: {
				From: state1,
				Targets: map[State]*Target{
					state2: {
						To:       state2,
						funcName: "condFunc1",
					},
				},
			},
		},
	}

	// Act
	action.WithTransition(state1, state3)

	// Assert
	assert.Len(t, action.transitions, 1)
	assert.Equal(t, state1, action.transitions[state1].From)
	assert.Len(t, action.transitions[state1].Targets, 2)
	assert.Equal(t, state2, action.transitions[state1].Targets[state2].To)
	assert.Equal(t, "condFunc1", action.transitions[state1].Targets[state2].funcName)
	assert.Equal(t, state3, action.transitions[state1].Targets[state3].To)
	assert.Equal(t, "", action.transitions[state1].Targets[state3].funcName)
}

func Test_Action_WithTransition_AddExistingTargetToExistingTransition_ShouldNotAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	action := Action{
		transitions: map[State]*Transition{
			state1: {
				From: state1,
				Targets: map[State]*Target{
					state2: {
						To:       state2,
						funcName: "",
					},
					state3: {
						To:       state3,
						funcName: "",
					},
				},
			},
		},
	}

	// Act
	action.WithTransition(state1, state3)

	// Assert
	assert.Len(t, action.transitions, 1)
	assert.Equal(t, state1, action.transitions[state1].From)
	assert.Len(t, action.transitions[state1].Targets, 2)
	assert.Equal(t, state2, action.transitions[state1].Targets[state2].To)
	assert.Equal(t, state3, action.transitions[state1].Targets[state3].To)
}

func Test_StateMachine_WithCommand_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	const (
		cmd1 Command = iota
		cmd2
		cmd3
	)
	element := testElement{state: state1}

	fsm := New(&element)

	// Act
	fsm.WithCommand(cmd1, "Command1")

	// Assert
	assert.Len(t, fsm.actions, 1)
	assert.Equal(t, cmd1, fsm.actions[cmd1].command)
}

func Test_StateMachine_WithCommand_WithTransition_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	const (
		cmd1 Command = iota
		cmd2
		cmd3
	)

	element := testElement{state: state1}

	fsm := New(&element)

	// Act
	fsm.WithCommand(cmd1, "Command1").
		WithTransition(state1, state2).
		WithTransition(state1, state3)

	// Assert
	assert.Len(t, fsm.actions, 1)
	assert.Equal(t, cmd1, fsm.actions[cmd1].command)
	assert.Len(t, fsm.actions[cmd1].transitions, 1)
	assert.Equal(t, state1, fsm.actions[cmd1].transitions[state1].From)
	assert.Len(t, fsm.actions[cmd1].transitions[state1].Targets, 2)
	assert.Equal(t, state2, fsm.actions[cmd1].transitions[state1].Targets[state2].To)
	assert.Equal(t, state3, fsm.actions[cmd1].transitions[state1].Targets[state3].To)
}
