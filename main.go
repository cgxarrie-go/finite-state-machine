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

func (fsm *Fsm) ExecuteCommand(command Command) bool {
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

func printSectionBreak() {
	fmt.Println("")
	fmt.Println("--------------------------------------")
}

func printCommandExecuted(executed bool, command Command, from State) {
	if executed {
		msg := fmt.Sprintf("Transition found for Command %s From State %s", command, from)
		fmt.Println(msg)
		fmt.Println("Command executed", command)
	} else {
		msg := fmt.Sprintf("Transition NOT found for Command %s From State %s", command, from)
		fmt.Println(msg)
	}
}

func executeFsmCommand(fsm *Fsm, command Command, expectedState State) {
	initialState := fsm.State
	msg := fmt.Sprintf("Execute command %s from %s -> expected state : %s", command, fsm.State, expectedState)
	fmt.Println(msg)

	executed := fsm.ExecuteCommand(command)
	printCommandExecuted(executed, command, initialState)
}

func printFsmTransitions(fsm Fsm) {
	fmt.Println("Fsm transitions count ", len(fsm.Transitions))
	for i := 0; i < len(fsm.Transitions); i++ {
		msg := fmt.Sprintf("When %s from %s goes to %s", fsm.Transitions[i].Command, fsm.Transitions[i].From, fsm.Transitions[i].To)
		fmt.Println(msg)
	}
}

func main() {
	printSectionBreak()
	fmt.Println("Initialize state machine")
	fsm := New()

	fmt.Println("Fsm Current State ", fsm.State)
	printFsmTransitions(fsm)

	printSectionBreak()
	executeFsmCommand(&fsm, InsertCoin, Unlocked)
	fmt.Println("Fsm Current State ", fsm.State)

	printSectionBreak()
	executeFsmCommand(&fsm, InsertCoin, Unlocked)
	fmt.Println("Fsm Current State ", fsm.State)

	printSectionBreak()
	executeFsmCommand(&fsm, PushButton, Locked)
	fmt.Println("Fsm Current State ", fsm.State)

	printSectionBreak()
	executeFsmCommand(&fsm, PushButton, Locked)
	fmt.Println("Fsm Current State ", fsm.State)
}
