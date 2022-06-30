package fsm

import (
	"github.com/cgxarrie/fsm-go/fsm/commands"
	"github.com/cgxarrie/fsm-go/fsm/states"
)

type StateMachine struct {
	State       states.State
	Transitions []Transition
}

func New() StateMachine {
	fsm := &StateMachine{}
	fsm.State = states.Locked

	fsm.AddTransition(states.Locked, commands.InsertCoin, states.Unlocked)
	fsm.AddTransition(states.Unlocked, commands.PushButton, states.Locked)

	return *fsm
}

func (fsm *StateMachine) AddTransition(from states.State, command commands.Command, to states.State) bool {

	_, exists := fsm.transitionExists(command)

	if exists {
		return false
	}

	fsm.Transitions = append(fsm.Transitions, Transition{From: from, Command: command, To: to})
	return true
}

func (fsm *StateMachine) ExecuteCommand(command commands.Command) bool {

	transition, exists := fsm.transitionExists(command)

	if !exists {
		return false
	}

	fsm.State = transition.To
	return true
}

func (fsm StateMachine) transitionExists(command commands.Command) (Transition, bool) {

	for _, transition := range fsm.Transitions {

		if transition.From != fsm.State || transition.Command != command {
			continue
		}

		if transition.From == fsm.State && transition.Command == command {
			return transition, true
		}
	}

	return Transition{"", "", ""}, false
}
