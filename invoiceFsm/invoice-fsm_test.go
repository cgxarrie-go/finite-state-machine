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

	for i := 0; i < len(sm.Transitions); i++ {

		transition := sm.Transitions[i]

		if transition.From == fsm.State(InvoiceStateDraft) &&
			transition.Command == fsm.Command(InvoiceCommandConfirm) &&
			transition.To == fsm.State(InvoiceStateWaitingForApproval) {
			continue
		}
		if transition.From == fsm.State(InvoiceStateWaitingForApproval) &&
			transition.Command == fsm.Command(InvoiceCommandReject) &&
			transition.To == fsm.State(InvoiceStateRejected) {
			continue
		}
		if transition.From == fsm.State(InvoiceStateWaitingForApproval) &&
			transition.Command == fsm.Command(InvoiceCommandApprove) &&
			transition.To == fsm.State(InvoiceStateWaitingForPayment) {
			continue
		}
		if transition.From == fsm.State(InvoiceStateWaitingForPayment) &&
			transition.Command == fsm.Command(InvoiceCommandPay) &&
			transition.To == fsm.State(InvoiceStateCompleted) {
			continue
		}

		t.Errorf("Unexpected transitions found From %s Command %s To %s", fmt.Sprint(transition.From), fmt.Sprint(transition.Command), fmt.Sprint(transition.To))
	}

	if len(sm.Transitions) != 4 {
		t.Errorf("Unexpected initial transitions count : Got %d , expected %d", len(sm.Transitions), 4)
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

		err := sm.ExecuteCommand(fsm.Command(data.command))

		if data.expectedError {
			if err == nil {
				t.Errorf("Command %s from State %s should throw error", fmt.Sprint(data.command), fmt.Sprint(data.fromState))
				continue
			}
		} else {
			if err != nil {
				t.Errorf("Command %s from State %s thrown error", fmt.Sprint(data.command), fmt.Sprint(data.fromState))
				continue
			}

			if sm.State != fsm.State(data.toState) {
				t.Errorf("Command %s from State %s should change state to %s", fmt.Sprint(data.command), fmt.Sprint(data.fromState), fmt.Sprint(data.toState))
			}
		}
	}
}
