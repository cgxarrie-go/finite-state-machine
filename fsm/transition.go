package fsm

type Transition struct {
	From    State
	Command Command
	To      State
}
