package fsm

type State uint32

const (
	Locked State = iota
	Unlocked
)

type Command uint32

const (
	InsertCoin Command = iota
	PushButton
)

type Transition struct {
	From    State
	Command Command
	To      State
}

type StateMachine struct {
	State       State
	Transitions [2]Transition
}

func New() StateMachine {
	fsm := &StateMachine{}
	fsm.State = Locked
	fsm.Transitions[0] = Transition{From: Locked, Command: InsertCoin, To: Unlocked}
	fsm.Transitions[1] = Transition{From: Unlocked, Command: PushButton, To: Locked}

	return *fsm
}

func (fsm *StateMachine) ExecuteCommand(command Command) bool {
	for _, transition := range fsm.Transitions {

		if transition.From != fsm.State || transition.Command != command {
			continue
		}

		if transition.From == fsm.State && transition.Command == command {
			fsm.State = transition.To
			return true
		}
	}

	return false
}
