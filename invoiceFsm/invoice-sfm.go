package invoiceFsm

import (
	"github.com/cgxarrie/fsm-go/fsm"
)

type InvoiceState fsm.State

const (
	InvoiceStateDraft InvoiceState = iota
	InvoiceStateWaitingForApproval
	InvoiceStateWaitingForPayment
	InvoiceStateRejected
	InvoiceStateCompleted
)

type InvoiceCommand fsm.Command

const (
	InvoiceCommandConfirm InvoiceCommand = iota
	InvoiceCommandReject
	InvoiceCommandApprove
	InvoiceCommandPay
)

func NewInvoiceStateMachine() fsm.StateMachine {

	sm := fsm.New(fsm.State(InvoiceStateDraft))
	sm.AddTransition(fsm.State(InvoiceStateDraft), fsm.Command(InvoiceCommandConfirm), fsm.State(InvoiceStateWaitingForApproval))
	sm.AddTransition(fsm.State(InvoiceStateWaitingForApproval), fsm.Command(InvoiceCommandReject), fsm.State(InvoiceStateRejected))
	sm.AddTransition(fsm.State(InvoiceStateWaitingForApproval), fsm.Command(InvoiceCommandApprove), fsm.State(InvoiceStateWaitingForPayment))
	sm.AddTransition(fsm.State(InvoiceStateWaitingForPayment), fsm.Command(InvoiceCommandPay), fsm.State(InvoiceStateCompleted))

	return sm
}
