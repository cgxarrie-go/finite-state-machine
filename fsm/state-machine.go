package fsm

import "fmt"

type State uint32

type Command uint32

type Transition struct {
	From    State
	Command Command
	To      State
}

type StateMachine struct {
	State       State
	Transitions []Transition
}

func New(initialState State) StateMachine {
	fsm := &StateMachine{}
	fsm.State = initialState

	return *fsm
}

func (fsm *StateMachine) AddTransition(from State, command Command, to State) bool {

	_, err := fsm.findTransition(command)
	if err == nil {
		return false
	}

	fsm.Transitions = append(fsm.Transitions, Transition{From: from, Command: command, To: to})
	return true
}

func (fsm *StateMachine) ExecuteCommand(command Command) (bool, error) {
	transition, err := fsm.findTransition(command)

	if err != nil {
		return false, err
	}

	fsm.State = transition.To
	return true, nil
}

func (fsm StateMachine) findTransition(command Command) (*Transition, error) {

	for _, transition := range fsm.Transitions {

		if transition.From != fsm.State || transition.Command != command {
			continue
		}

		if transition.From == fsm.State && transition.Command == command {
			return &transition, nil
		}
	}

	return nil, fmt.Errorf("Transition not found : Command%d from status %d", command, fsm.State)
}
