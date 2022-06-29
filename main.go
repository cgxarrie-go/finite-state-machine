package main

import (
	"fmt"
)

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

func New() Fsm {
	fsm := &Fsm{}
	fsm.State = Locked
	fsm.Transitions[0] = Transition{From: Locked, Command: InsertCoin, To: Unlocked}
	fsm.Transitions[1] = Transition{From: Unlocked, Command: PushButton, To: Locked}

	return *fsm
}

func (fsm *Fsm) ExecuteCommand(command Command) {
	for _, transition := range fsm.Transitions {

		if transition.From != fsm.State || transition.Command != command {
			continue
		}

		if transition.From == fsm.State && transition.Command == command {
			fmt.Println("Transition found (CFT)", transition.Command, transition.From, transition.To)
			fsm.State = transition.To
			return
		}
	}

	fmt.Println("Transition NOT found")
}

func printSectionBreak() {
	fmt.Println("")
	fmt.Println("--------------------------------------")
}

func main() {
	printSectionBreak()
	fmt.Println("Initialize state machine")
	fsm := New()

	fmt.Println("Fsm Current State ", fsm.State)
	fmt.Println("Fsm transitions count ", len(fsm.Transitions))
	for i := 0; i < len(fsm.Transitions); i++ {
		msg := fmt.Sprintf("When %s from %s goes to %s", fsm.Transitions[i].Command, fsm.Transitions[i].From, fsm.Transitions[i].To)
		fmt.Println(msg)
	}

	printSectionBreak()
	fmt.Println("Execute command InsertCoin from Locked -> expected state : ", Unlocked)

	fsm.ExecuteCommand(InsertCoin)
	fmt.Println("Fsm Current State ", fsm.State)

	printSectionBreak()
	fmt.Println("Execute command InsertCoin from Unlocked -> expected state : ", Unlocked)

	fsm.ExecuteCommand(InsertCoin)
	fmt.Println("Fsm Current State ", fsm.State)

	printSectionBreak()
	fmt.Println("Execute command PushButton from Unlocked -> expected state : ", Locked)

	fsm.ExecuteCommand(PushButton)
	fmt.Println("Fsm Current State ", fsm.State)

	printSectionBreak()
	fmt.Println("Execute command PushButton from Locked -> expected state : ", Locked)

	fsm.ExecuteCommand(PushButton)
	fmt.Println("Fsm Current State ", fsm.State)

}
