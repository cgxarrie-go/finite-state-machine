package invoiceFsm

import (
	"fmt"
	"testing"

	"github.com/cgxarrie/fsm-go/fsm"
)

func TestNewInvoiceStateMachineShouldInitialize(t *testing.T) {
	sm := NewInvoiceStateMachine()

	if sm.State != fsm.State(InvoiceStateDraft) {
		t.Errorf("Unexpected initial state: Got %d , expected %d", sm.State, InvoiceStateDraft)
	}
}

var commandTests = []struct {
	command       InvoiceCommand
	fromState     InvoiceState
	toState       InvoiceState
	expectedError bool
}{
	{InvoiceCommandConfirm, InvoiceStateDraft, InvoiceStateWaitingForApproval, false},
	{InvoiceCommandConfirm, InvoiceStateWaitingForApproval, InvoiceStateWaitingForApproval, true},
	{InvoiceCommandConfirm, InvoiceStateRejected, InvoiceStateRejected, true},
	{InvoiceCommandConfirm, InvoiceStateCompleted, InvoiceStateCompleted, true},
	{InvoiceCommandConfirm, InvoiceStateWaitingForPayment, InvoiceStateWaitingForPayment, true},
	{InvoiceCommandReject, InvoiceStateDraft, InvoiceStateDraft, true},
	{InvoiceCommandReject, InvoiceStateWaitingForApproval, InvoiceStateRejected, false},
	{InvoiceCommandReject, InvoiceStateRejected, InvoiceStateRejected, true},
	{InvoiceCommandReject, InvoiceStateWaitingForPayment, InvoiceStateWaitingForPayment, true},
	{InvoiceCommandReject, InvoiceStateCompleted, InvoiceStateCompleted, true},
	{InvoiceCommandApprove, InvoiceStateDraft, InvoiceStateDraft, true},
	{InvoiceCommandApprove, InvoiceStateWaitingForApproval, InvoiceStateWaitingForPayment, false},
	{InvoiceCommandApprove, InvoiceStateRejected, InvoiceStateRejected, true},
	{InvoiceCommandApprove, InvoiceStateWaitingForPayment, InvoiceStateWaitingForPayment, true},
	{InvoiceCommandApprove, InvoiceStateCompleted, InvoiceStateCompleted, true},
	{InvoiceCommandPay, InvoiceStateDraft, InvoiceStateDraft, true},
	{InvoiceCommandPay, InvoiceStateWaitingForApproval, InvoiceStateWaitingForApproval, true},
	{InvoiceCommandPay, InvoiceStateRejected, InvoiceStateRejected, true},
	{InvoiceCommandPay, InvoiceStateWaitingForPayment, InvoiceStateCompleted, false},
	{InvoiceCommandPay, InvoiceStateCompleted, InvoiceStateCompleted, true},
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
