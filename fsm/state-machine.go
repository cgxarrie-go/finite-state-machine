package fsm

type State uint32

type Command uint32

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

	transition := fsm.findTransition(command)
	if transition != nil {
		return false
	}

	fsm.Transitions = append(fsm.Transitions, Transition{From: from, Command: command, To: to})
	return true
}

func (fsm *StateMachine) ExecuteCommand(command Command) *CommandNotAvailableError {
	transition := fsm.findTransition(command)
	if transition == nil {
		err := &CommandNotAvailableError{
			FromState: fsm.State,
			Command:   command,
		}
		return err
	}

	fsm.State = transition.To
	return nil
}

func (fsm StateMachine) findTransition(command Command) *Transition {

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
