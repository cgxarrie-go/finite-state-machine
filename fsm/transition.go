package fsm

import (
	"github.com/cgxarrie/fsm-go/fsm/commands"
	"github.com/cgxarrie/fsm-go/fsm/states"
)

type Transition struct {
	From    states.State
	Command commands.Command
	To      states.State
}
