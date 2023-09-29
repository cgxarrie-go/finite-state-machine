package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Action_WithTransition_AddFirstTransition_ShouldAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	action := Action{}

	// Act
	action.WithTransition(state1, state2, func() bool { return true })

	// Assert
	assert.Len(t, action.transitions, 1)
	assert.Equal(t, state1, action.transitions[state1].From)
	assert.Len(t, action.transitions[state1].Targets, 1)
	assert.Equal(t, state2, action.transitions[state1].Targets[state2].To)
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
						To:        state2,
						Condition: func() bool { return true },
					},
				},
			},
		},
	}

	// Act
	action.WithTransition(state1, state3, func() bool { return true })

	// Assert
	assert.Len(t, action.transitions, 1)
	assert.Equal(t, state1, action.transitions[state1].From)
	assert.Len(t, action.transitions[state1].Targets, 2)
	assert.Equal(t, state2, action.transitions[state1].Targets[state2].To)
	assert.Equal(t, state3, action.transitions[state1].Targets[state3].To)
}

func Test_Action_WithTransition_AddExistingTargetToExistingTransition_ShouldNotAdd(t *testing.T) {
	// Arrange
	const (
		state1 State = iota
		state2
		state3
	)

	cond12 := func() bool { return true }
	cond13 := func() bool { return true }

	action := Action{
		transitions: map[State]*Transition{
			state1: {
				From: state1,
				Targets: map[State]*Target{
					state2: {
						To:        state2,
						Condition: cond12,
					},
					state3: {
						To:        state3,
						Condition: cond13,
					},
				},
			},
		},
	}

	// Act
	action.WithTransition(state1, state3, func() bool { return true })

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

	fsm := New(state1)

	// Act
	fsm.WithCommand(cmd1, func() error { return nil })

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

	fsm := New(state1)

	// Act
	fsm.WithCommand(cmd1, func() error { return nil }).
		WithTransition(state1, state2, func() bool { return true }).
		WithTransition(state1, state3, func() bool { return true })

	// Assert
	assert.Len(t, fsm.actions, 1)
	assert.Equal(t, cmd1, fsm.actions[cmd1].command)
	assert.Len(t, fsm.actions[cmd1].transitions, 1)
	assert.Equal(t, state1, fsm.actions[cmd1].transitions[state1].From)
	assert.Len(t, fsm.actions[cmd1].transitions[state1].Targets, 2)
	assert.Equal(t, state2, fsm.actions[cmd1].transitions[state1].Targets[state2].To)
	assert.Equal(t, state3, fsm.actions[cmd1].transitions[state1].Targets[state3].To)
}
