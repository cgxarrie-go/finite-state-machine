package invoiceFsm

import (
	"github.com/cgxarrie/fsm-go/fsm"
)

type InvoiceState fsm.State

const InvoiceStateDraft InvoiceState = 0
const InvoiceStateWaitingForApproval InvoiceState = 1
const InvoiceStateWaitingForPayment InvoiceState = 2
const InvoiceStateRejected InvoiceState = 3
const InvoiceStateCompleted InvoiceState = 4

type InvoiceCommand fsm.Command

const InvoiceCommandConfirm InvoiceCommand = 0
const InvoiceCommandReject InvoiceCommand = 1
const InvoiceCommandApprove InvoiceCommand = 2
const InvoiceCommandPay InvoiceCommand = 3

func NewInvoiceStateMachine() fsm.StateMachine {

	sm := fsm.New(fsm.State(InvoiceStateDraft))
	sm.AddTransition(fsm.State(InvoiceStateDraft), fsm.Command(InvoiceCommandConfirm), fsm.State(InvoiceStateWaitingForApproval))
	sm.AddTransition(fsm.State(InvoiceStateWaitingForApproval), fsm.Command(InvoiceCommandReject), fsm.State(InvoiceStateRejected))
	sm.AddTransition(fsm.State(InvoiceStateWaitingForApproval), fsm.Command(InvoiceCommandApprove), fsm.State(InvoiceStateWaitingForPayment))
	sm.AddTransition(fsm.State(InvoiceStateWaitingForPayment), fsm.Command(InvoiceCommandPay), fsm.State(InvoiceStateCompleted))

	return sm
}
