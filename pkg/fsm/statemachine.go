package fsm

import (
	"fmt"
)

type State uint32

type Command uint32

type HandledElement interface {
	SetState(State)
	State() State
}

type StateMachine struct {
	element HandledElement
	actions map[Command]*Action
}

type Action struct {
	command     Command
	fn    		func() error
	transitions map[State]*Transition
}

type Transition struct {
	From    State
	Targets map[State]*Target
}

type Target struct {
	To       State
	condition func() bool
}

func New(element HandledElement) StateMachine {
	fsm := &StateMachine{
		element: element,
	}

	return *fsm
}

func (fsm *StateMachine) WithCommand(command Command, fn func() error) *Action {
	if fsm.actions == nil {
		fsm.actions = make(map[Command]*Action)
	}

	if action, ok := fsm.actions[command]; ok {
		return action
	}

	fsm.actions[command] = &Action{
		fn: fn,
		command:  command,
	}

	action := fsm.actions[command]
	return action
}

func (a *Action) WithTransition(from, to State) *Action {
	return a.WithConditionedTransition(from, to, nil)
}

func (a *Action) WithConditionedTransition(from, to State, condition func() bool) *Action {
	if a.transitions == nil {
		a.transitions = make(map[State]*Transition)
	}

	if _, ok := a.transitions[from]; !ok {
		a.transitions[from] = &Transition{
			From: from,
			Targets: map[State]*Target{
				to: {
					To:       to,
					condition: condition,
				},
			},
		}
		return a
	}

	if _, ok := a.transitions[from].Targets[to]; !ok {
		a.transitions[from].Targets[to] = &Target{
			To:       to,
			condition: condition,
		}
		return a
	}

	return a
}

func (fsm *StateMachine) ExecuteCommand(command Command) error {

	if fsm.actions == nil {
		return fmt.Errorf("cannot execute requested command %v from state %v",
			command, fsm.element.State())
	}

	action, ok := fsm.actions[command]
	if !ok {
		return fmt.Errorf("cannot find action for command %v", command)
	}

	if action.transitions == nil {
		return fmt.Errorf("cannot find transitions for command %v",
			command)
	}

	transition, ok := action.transitions[fsm.element.State()]
	if !ok {
		return fmt.Errorf("cannot find transitions for command %v "+
			"and state %v", command, fsm.element)
	}

	err := action.fn()
	if err != nil {
		return fmt.Errorf("command %v returned error: %v", command, err)
	}

	for _, target := range transition.Targets {
		if target.condition != nil {
			if !target.condition() {
				continue
			}

			fsm.element.SetState(target.To)
			return nil
		}

		fsm.element.SetState(target.To)
		return nil
	}

	return fmt.Errorf("cannot find executable transition for command %v "+
		"and state %v", command, fsm.element.State())
}
