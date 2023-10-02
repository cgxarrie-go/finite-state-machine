package fsm

import (
	"fmt"
)

type State uint32
type CommandID uint32
type Action func() error
type Condition func() bool

type Commands map[CommandID]Action
type Targets map[State]Condition
type Transitions map[State]map[CommandID]Targets

type SMObject interface {
	SetState(State)
	State() State
}

type StateMachine struct {
	smObject SMObject
	commands Commands
	transitions Transitions
}


func New(element SMObject) StateMachine {
	fsm := &StateMachine{
		smObject: element,
		transitions: Transitions{},
		commands: Commands{},
	}

	return *fsm
}

func (fsm *StateMachine) WithCommand(id CommandID, action Action) *StateMachine {
	fsm.commands[id] = action
	return fsm
}


func (fsm *StateMachine) From(s State) *TransitionBuilder {
	t := &TransitionBuilder{
		sm:   fsm,
		from: s,
	}
	return t
}

func (fsm StateMachine) Do(cmdID CommandID) error {

	from := fsm.smObject.State()

	if _, ok := fsm.transitions[fsm.smObject.State()]; !ok {
		return fmt.Errorf("cannot execute requested command %v from state %v",
			cmdID, from)
	}

	action, ok := fsm.commands[cmdID]
	if !ok {
		return fmt.Errorf("command %v not found", cmdID)
	}

	if action == nil {
		return fmt.Errorf("no action found for command %v", cmdID)
	}

	err := action()
	if err != nil {
		return fmt.Errorf("command %v from status %v returned error: %v", 
			cmdID, fsm.smObject.State(), err)
	}
	

	for toState, condition := range fsm.transitions[from][cmdID] {
		if condition != nil {
			if !condition() {
				continue
			}

			fsm.smObject.SetState(toState)
			return nil
		}

		fsm.smObject.SetState(toState)
		return nil
	}

	return fmt.Errorf("cannot find executable transition for command %v "+
		"and state %v", cmdID, from)
}
