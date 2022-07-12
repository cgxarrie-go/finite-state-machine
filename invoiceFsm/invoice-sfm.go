package invoiceFsm

import (
	"github.com/cgxarrie/fsm-go/fsm"
)

type InvoiceState fsm.State

const (
	draft InvoiceState = iota
	waitingForApproval
	waitingForPayment
	rejected
	completed
)

type InvoiceCommand fsm.Command

const (
	confirm InvoiceCommand = iota
	reject
	approve
	pay
)

func NewInvoiceStateMachine() fsm.StateMachine {

	sm := fsm.New(fsm.State(draft))

	sm.
		WithTransition().On(fsm.Command(confirm)).From(fsm.State(draft)).To(fsm.State(waitingForApproval)).Add().
		WithTransition().On(fsm.Command(reject)).From(fsm.State(waitingForApproval)).To(fsm.State(rejected)).Add().
		WithTransition().On(fsm.Command(approve)).From(fsm.State(waitingForApproval)).To(fsm.State(waitingForPayment)).Add().
		WithTransition().On(fsm.Command(pay)).From(fsm.State(waitingForPayment)).To(fsm.State(completed)).Add()

	return sm
}
