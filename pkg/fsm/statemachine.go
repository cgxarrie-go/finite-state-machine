package fsm

import (
	"fmt"
)

type State uint32

type Command uint32

type StateMachine struct {
	state   State
	actions map[Command]*Action
}

type Action struct {
	command     Command
	fn          func() error
	transitions map[State]*Transition
}

type Transition struct {
	From    State
	Targets map[State]*Target
}

type Target struct {
	To        State
	Condition func() bool
}

func New(initialState State) StateMachine {
	fsm := &StateMachine{
		state: initialState,
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
		command: command,
		fn:      fn,
	}

	action := fsm.actions[command]
	return action
}

func (a *Action) WithTransition(from, to State, condition func() bool) *Action {
	if a.transitions == nil {
		a.transitions = make(map[State]*Transition)
	}

	if _, ok := a.transitions[from]; !ok {
		a.transitions[from] = &Transition{
			From: from,
			Targets: map[State]*Target{
				to: {
					To:        to,
					Condition: condition,
				},
			},
		}
		return a
	}

	if _, ok := a.transitions[from].Targets[to]; !ok {
		a.transitions[from].Targets[to] = &Target{
			To:        to,
			Condition: condition,
		}
		return a
	}

	return a
}

func (fsm *StateMachine) ExecuteCommand(command Command) (bool, error) {

	if fsm.actions == nil {
		return false, fmt.Errorf("cannot execute requested command %v from state %v",
			command, fsm.state)
	}

	action, ok := fsm.actions[command]
	if !ok {
		return false, fmt.Errorf("cannot find action for command %v", command)
	}

	if action.transitions == nil {
		return false, fmt.Errorf("cannot find transitions for command %v",
			command)
	}

	transition, ok := action.transitions[fsm.state]
	if !ok {
		return false,
			fmt.Errorf("cannot find transitions for command %v and state %v",
				command, fsm.state)
	}

	for _, target := range transition.Targets {
		if target.Condition != nil && !target.Condition() {
			continue
		}
		err := action.fn()
		if err != nil {
			return false,
				fmt.Errorf("executing action for command %v and state %v failed: %v",
					command, fsm.state, err)
		}
		fsm.state = target.To
		return true, nil
	}

	return false,
		fmt.Errorf("cannot find executable transition for command %v and state %v",
			command, fsm.state)
}
