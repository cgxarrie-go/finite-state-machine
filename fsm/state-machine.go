package fsm

import (
	"github.com/cgxarrie/fsm-go/fsm/commands"
	"github.com/cgxarrie/fsm-go/fsm/states"
)

type StateMachine struct {
	State       states.State
	Transitions []Transition
}

func New(initialState states.State) StateMachine {
	fsm := &StateMachine{}
	fsm.State = initialState

	return *fsm
}

func (fsm *StateMachine) AddTransition(from states.State, command commands.Command, to states.State) bool {

	transition := fsm.findTransition(command)
	if transition != nil {
		return false
	}

	fsm.Transitions = append(fsm.Transitions, Transition{From: from, Command: command, To: to})
	return true
}

func (fsm *StateMachine) ExecuteCommand(command commands.Command) bool {

	transition := fsm.findTransition(command)
	if transition == nil {
		return false
	}

	fsm.State = transition.To
	return true
}

func (fsm StateMachine) findTransition(command commands.Command) *Transition {

	for _, transition := range fsm.Transitions {

		if transition.From != fsm.State || transition.Command != command {
			continue
		}

		if transition.From == fsm.State && transition.Command == command {
			return &transition
		}
	}

	return nil
}
