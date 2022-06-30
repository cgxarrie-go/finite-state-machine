package main

import (
	"fmt"

	"github.com/cgxarrie/fsm-go/fsm"
	"github.com/cgxarrie/fsm-go/fsm/commands"
	"github.com/cgxarrie/fsm-go/fsm/states"
)

func printSectionBreak() {
	fmt.Println("")
	fmt.Println("--------------------------------------")
}

func printCommandExecuted(executed bool, command commands.Command, from states.State) {
	if executed {
		msg := fmt.Sprintf("Transition found for Command %s From State %s", command, from)
		fmt.Println(msg)
		fmt.Println("Command executed", command)
	} else {
		msg := fmt.Sprintf("Transition NOT found for Command %s From State %s", command, from)
		fmt.Println(msg)
	}
}

func executeFsmCommand(fsm *fsm.StateMachine, command commands.Command, expectedState states.State) {
	initialState := fsm.State
	msg := fmt.Sprintf("Execute command %s from %s -> expected state : %s", command, fsm.State, expectedState)
	fmt.Println(msg)

	executed := fsm.ExecuteCommand(command)
	printCommandExecuted(executed, command, initialState)
}

func printFsmTransitions(fsm fsm.StateMachine) {
	fmt.Println("Fsm transitions count ", len(fsm.Transitions))
	for i := 0; i < len(fsm.Transitions); i++ {
		msg := fmt.Sprintf("When %s from %s goes to %s", fsm.Transitions[i].Command, fsm.Transitions[i].From, fsm.Transitions[i].To)
		fmt.Println(msg)
	}
}

func main() {
	printSectionBreak()
	fmt.Println("Initialize state machine")
	sm := fsm.New()

	fmt.Println("Fsm Current State ", sm.State)
	printFsmTransitions(sm)

	printSectionBreak()
	executeFsmCommand(&sm, commands.InsertCoin, states.Unlocked)
	fmt.Println("Fsm Current State ", sm.State)

	printSectionBreak()
	executeFsmCommand(&sm, commands.InsertCoin, states.Unlocked)
	fmt.Println("Fsm Current State ", sm.State)

	printSectionBreak()
	executeFsmCommand(&sm, commands.PushButton, states.Locked)
	fmt.Println("Fsm Current State ", sm.State)

	printSectionBreak()
	executeFsmCommand(&sm, commands.PushButton, states.Locked)
	fmt.Println("Fsm Current State ", sm.State)
}
