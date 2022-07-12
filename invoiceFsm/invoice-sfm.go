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
	sm.WithTransition(fsm.Command(InvoiceCommandConfirm), fsm.State(InvoiceStateDraft), fsm.State(InvoiceStateWaitingForApproval)).
		WithTransition(fsm.Command(InvoiceCommandReject), fsm.State(InvoiceStateWaitingForApproval), fsm.State(InvoiceStateRejected)).
		WithTransition(fsm.Command(InvoiceCommandApprove), fsm.State(InvoiceStateWaitingForApproval), fsm.State(InvoiceStateWaitingForPayment)).
		WithTransition(fsm.Command(InvoiceCommandPay), fsm.State(InvoiceStateWaitingForPayment), fsm.State(InvoiceStateCompleted))

	return sm
}
