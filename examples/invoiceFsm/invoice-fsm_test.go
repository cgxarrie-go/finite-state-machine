package invoiceFsm

import (
	"fmt"
	"testing"

	"github.com/cgxarrie/fsm-go/fsm"
)

func TestNewInvoiceStateMachineShouldInitialize(t *testing.T) {
	sm := NewInvoiceStateMachine()

	if sm.State != fsm.State(draft) {
		t.Errorf("Unexpected initial state: Got %d , expected %d", sm.State, draft)
	}
}

var commandTests = []struct {
	command       InvoiceCommand
	fromState     InvoiceState
	toState       InvoiceState
	expectedError bool
}{
	{confirm, draft, waitingForApproval, false},
	{confirm, waitingForApproval, waitingForApproval, true},
	{confirm, rejected, rejected, true},
	{confirm, completed, completed, true},
	{confirm, waitingForPayment, waitingForPayment, true},
	{reject, draft, draft, true},
	{reject, waitingForApproval, rejected, false},
	{reject, rejected, rejected, true},
	{reject, waitingForPayment, waitingForPayment, true},
	{reject, completed, completed, true},
	{approve, draft, draft, true},
	{approve, waitingForApproval, waitingForPayment, false},
	{approve, rejected, rejected, true},
	{approve, waitingForPayment, waitingForPayment, true},
	{approve, completed, completed, true},
	{pay, draft, draft, true},
	{pay, waitingForApproval, waitingForApproval, true},
	{pay, rejected, rejected, true},
	{pay, waitingForPayment, completed, false},
	{pay, completed, completed, true},
}

func TestExecuteCommandShouldBehaveAsExpected(t *testing.T) {
	for _, data := range commandTests {
		sm := NewInvoiceStateMachine()
		sm.State = fsm.State(data.fromState)

		_, err := sm.ExecuteCommand(fsm.Command(data.command))

		if data.expectedError {
			if err == nil {
				t.Errorf("Command %s from State %s should throw error", fmt.Sprint(data.command), fmt.Sprint(data.fromState))
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Command %s from State %s thrown error : %v", fmt.Sprint(data.command), fmt.Sprint(data.fromState), err)
				continue
			}

			if sm.State != fsm.State(data.toState) {
				t.Errorf("Command %s from State %s should change state to %s", fmt.Sprint(data.command), fmt.Sprint(data.fromState), fmt.Sprint(data.toState))
			}
		}
	}
}
