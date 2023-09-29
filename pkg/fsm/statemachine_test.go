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

	action := Command{}

	// Act
	action.WithTransition(state1, state2)

	// Assert
	assert.Len(t, action.transitions, 1)
	assert.Equal(t, state1, action.transitions[state1].From)
	assert.Len(t, action.transitions[state1].Targets, 1)
	assert.Equal(t, state2, action.transitions[state1].Targets[state2].To)
	assert.Nil(t, action.transitions[state1].Targets[state2].condition)
}

func Test_Action_WithConditionedTransition_AddFirstTransition_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	action := Command{}

	// Act
	action.WithConditionedTransition(state1, state2, func() bool {return true})

	// Assert
	assert.Len(t, action.transitions, 1)
	assert.Equal(t, state1, action.transitions[state1].From)
	assert.Len(t, action.transitions[state1].Targets, 1)
	assert.Equal(t, state2, action.transitions[state1].Targets[state2].To)
	assert.NotNil(t, action.transitions[state1].Targets[state2].condition)
}

func Test_Action_WithTransition_AddTargetToExistingTransition_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	action := Command{
		transitions: map[State]*Transition{
			state1: {
				From: state1,
				Targets: map[State]*Target{
					state2: {
						To:       state2,
						condition: func() bool {return true},
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
	assert.NotNil(t, action.transitions[state1].Targets[state2].condition)
	assert.Equal(t, state3, action.transitions[state1].Targets[state3].To)
	assert.Nil(t, action.transitions[state1].Targets[state3].condition)
}

func Test_Action_WithTransition_AddExistingTargetToExistingTransition_ShouldNotAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	action := Command{
		transitions: map[State]*Transition{
			state1: {
				From: state1,
				Targets: map[State]*Target{
					state2: {
						To:       state2,
						condition: nil,
					},
					state3: {
						To:       state3,
						condition: nil,
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

	element := testElement{state: state1}

	fsm := New(&element)

	// Act
	fsm.WithCommand("Command1", element.Command1)

	// Assert
	assert.Len(t, fsm.commands, 1)
	assert.Equal(t, "Command1", fsm.commands["Command1"].name)
}

func Test_StateMachine_WithCommand_WithTransition_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)


	element := testElement{state: state1}

	fsm := New(&element)

	// Act
	fsm.WithCommand("cmd1", element.Command1).
		WithTransition(state1, state2).
		WithTransition(state1, state3)

	// Assert
	assert.Len(t, fsm.commands, 1)
	assert.Equal(t, "cmd1", fsm.commands["cmd1"].name)
	assert.Len(t, fsm.commands["cmd1"].transitions, 1)
	assert.Equal(t, state1, fsm.commands["cmd1"].transitions[state1].From)
	assert.Len(t, fsm.commands["cmd1"].transitions[state1].Targets, 2)
	assert.Equal(t, state2, fsm.commands["cmd1"].transitions[state1].Targets[state2].To)
	assert.Equal(t, state3, fsm.commands["cmd1"].transitions[state1].Targets[state3].To)
}
