package fsm

import "fmt"

type CommandNotAvailableError struct {
	FromState State
	Command   Command
}

func (m *CommandNotAvailableError) Error() string {
	msg := fmt.Sprintf("Cannot execute command %s from status %s", fmt.Sprint(m.Command), fmt.Sprint(m.FromState))
	return msg
}
