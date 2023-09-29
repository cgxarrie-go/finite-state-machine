package fsm

import (
	"fmt"
)

type State uint32

type Command uint32

type StateMachine struct {
	state   State
	actions map[Command]Action
}

type Action struct {
	Command     Command
	Func        func() error
	Transitions map[State]Transition
}

type Transition struct {
	From    State
	Targets map[State]Target
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

func (t *Transition) withTarget(to State, condition func() bool) *Transition {
	if t.Targets == nil {
		t.Targets = make(map[State]Target)
	}

	if _, ok := t.Targets[to]; ok {
		return t
	}

	t.Targets[to] = Target{
		To:        to,
		Condition: condition,
	}

	return t
}

func (fsm *StateMachine) WithAction(command Command, fn func() error) *Action {
	if fsm.actions == nil {
		fsm.actions = make(map[Command]Action)
	}

	if action, ok := fsm.actions[command]; ok {
		return &action
	}

	fsm.actions[command] = Action{
		Command: command,
		Func:    fn,
	}

	action := fsm.actions[command]
	return &action
}

func (a *Action) WithTransition(from, to State, condition func() bool) *Action {
	if a.Transitions == nil {
		a.Transitions = make(map[State]Transition)
	}

	if transition, ok := a.Transitions[from]; ok {
		if _, ok := transition.Targets[to]; ok {
			return a
		}
		transition.withTarget(to, condition)
		return a
	}

	transition := Transition{
		From:    from,
		Targets: map[State]Target{},
	}
	transition.withTarget(to, condition)
	a.Transitions[from] = transition

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

	if action.Transitions == nil {
		return false, fmt.Errorf("cannot find transitions for command %v",
			command)
	}

	transition, ok := action.Transitions[fsm.state]
	if !ok {
		return false,
			fmt.Errorf("cannot find transitions for command %v and state %v",
				command, fsm.state)
	}

	for _, target := range transition.Targets {
		if target.Condition != nil && !target.Condition() {
			continue
		}
		err := action.Func()
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
