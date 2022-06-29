package main

import "fmt"

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

type Fsm struct {
	State       State
	Transitions [2]Transition
}

func main() {
	fsm := &Fsm{}
	fsm.State = Locked
	fsm.Transitions[0] = Transition{From: Locked, Command: InsertCoin, To: Unlocked}
	fsm.Transitions[1] = Transition{From: Unlocked, Command: PushButton, To: Locked}

	fmt.Println("Fsm transitions count ", len(fsm.Transitions))

	for i := 0; i < len(fsm.Transitions); i++ {
		msg := fmt.Sprintf("When %s from %s goes to %s", fsm.Transitions[i].Command, fsm.Transitions[i].From, fsm.Transitions[i].To)
		fmt.Println(msg)
	}
}
