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

	sm.WithTransition().
		On(fsm.Command(InvoiceCommandConfirm)).
		From(fsm.State(InvoiceStateDraft)).
		To(fsm.State(InvoiceStateWaitingForApproval)).
		Add()

	sm.WithTransition().
		On(fsm.Command(InvoiceCommandReject)).
		From(fsm.State(InvoiceStateWaitingForApproval)).
		To(fsm.State(InvoiceStateRejected)).
		Add()

	sm.WithTransition().
		On(fsm.Command(InvoiceCommandApprove)).
		From(fsm.State(InvoiceStateWaitingForApproval)).
		To(fsm.State(InvoiceStateWaitingForPayment)).
		Add()

	sm.WithTransition().
		On(fsm.Command(InvoiceCommandPay)).
		From(fsm.State(InvoiceStateWaitingForPayment)).
		To(fsm.State(InvoiceStateCompleted)).
		Add()

	return sm
}
