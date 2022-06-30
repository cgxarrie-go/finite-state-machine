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

	if len(sm.Transitions) != 4 {
		t.Errorf("Unexpected initial transitions count : Got %d , expected %d", len(sm.Transitions), 4)
	}
}

var commandTests = []struct {
	fromState     InvoiceState
	command       InvoiceCommand
	toState       InvoiceState
	expectedError bool
}{
	{InvoiceStateDraft, InvoiceCommandConfirm, InvoiceStateWaitingForApproval, false},
	{InvoiceStateDraft, InvoiceCommandApprove, InvoiceStateDraft, true},
	{InvoiceStateDraft, InvoiceCommandReject, InvoiceStateDraft, true},
	{InvoiceStateDraft, InvoiceCommandPay, InvoiceStateDraft, true},
	{InvoiceStateWaitingForApproval, InvoiceCommandConfirm, InvoiceStateWaitingForApproval, true},
	{InvoiceStateWaitingForApproval, InvoiceCommandApprove, InvoiceStateWaitingForPayment, false},
	{InvoiceStateWaitingForApproval, InvoiceCommandReject, InvoiceStateRejected, false},
	{InvoiceStateWaitingForApproval, InvoiceCommandPay, InvoiceStateWaitingForApproval, true},
	{InvoiceStateRejected, InvoiceCommandConfirm, InvoiceStateRejected, true},
	{InvoiceStateRejected, InvoiceCommandApprove, InvoiceStateRejected, true},
	{InvoiceStateRejected, InvoiceCommandReject, InvoiceStateRejected, true},
	{InvoiceStateRejected, InvoiceCommandPay, InvoiceStateRejected, true},
	{InvoiceStateWaitingForPayment, InvoiceCommandConfirm, InvoiceStateWaitingForPayment, true},
	{InvoiceStateWaitingForPayment, InvoiceCommandApprove, InvoiceStateWaitingForPayment, true},
	{InvoiceStateWaitingForPayment, InvoiceCommandReject, InvoiceStateWaitingForPayment, true},
	{InvoiceStateWaitingForPayment, InvoiceCommandPay, InvoiceStateCompleted, false},
	{InvoiceStateCompleted, InvoiceCommandConfirm, InvoiceStateCompleted, true},
	{InvoiceStateCompleted, InvoiceCommandApprove, InvoiceStateCompleted, true},
	{InvoiceStateCompleted, InvoiceCommandReject, InvoiceStateCompleted, true},
	{InvoiceStateCompleted, InvoiceCommandPay, InvoiceStateCompleted, true},
}

func TestConfirmFromDraftShouldChangeToWaitingForApproval(t *testing.T) {
	for _, data := range commandTests {
		sm := NewInvoiceStateMachine()

		err := sm.ExecuteCommand(fsm.Command(data.command))

		if data.expectedError && err == nil {
			t.Errorf("Command %s from State %s should throw error", fmt.Sprint(data.command), fmt.Sprint(data.fromState))
		}

		if !data.expectedError && sm.State != fsm.State(data.toState) {
			t.Errorf("Command %s from State %s should change state to %s", fmt.Sprint(data.command), fmt.Sprint(data.fromState), fmt.Sprint(data.toState))
		}
	}
}
