package fsm

import (
	"fmt"
)

type State uint32

type Command uint32

type transition map[Command]map[State]State

type StateMachine struct {
	State       State
	Transitions transition
}

func New(initialState State) StateMachine {
	fsm := &StateMachine{}
	fsm.State = initialState

	return *fsm
}

func (fsm *StateMachine) WithTransition(command Command, from State, to State) *StateMachine {

	if fsm.Transitions == nil {
		fsm.Transitions = transition{
			command: {
				from: to,
			},
		}
		return fsm
	}

	_, commandExists := fsm.Transitions[command]

	if !commandExists {
		fsm.Transitions[command] = map[State]State{from: to}
		return fsm
	}

	fsm.Transitions[command][from] = to
	return fsm
}

func (fsm *StateMachine) ExecuteCommand(command Command) (bool, error) {

	canExecute := fsm.canExecuteCommand(command)

	if !canExecute {
		return false, fmt.Errorf("cannot execute requested command %v from state %v", command, fsm.State)
	}

	fsm.State = fsm.Transitions[command][fsm.State]
	return true, nil
}

func (fsm StateMachine) canExecuteCommand(command Command) bool {

	commandStates, commandExist := fsm.Transitions[command]

	if !commandExist {
		return false
	}

	_, fromSatateExist := commandStates[fsm.State]

	return fromSatateExist
}
