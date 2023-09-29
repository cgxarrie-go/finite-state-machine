package invoiceFsm

import "github.com/cgxarrie/fsm-go/pkg/fsm"

type Invoice struct {
	state             InvoiceState
	isSignatureReceived bool
	isApproved          bool
	needsSignature    bool
}

func NewInvoice(needsSignature bool) Invoice {
	return Invoice{
		state:             draft,
		isSignatureReceived: false,
		isApproved:          false,
		needsSignature:    needsSignature,
	}
}

func (i *Invoice) SetState(state fsm.State) {
	i.state = InvoiceState(state)
}

func (i *Invoice) State() fsm.State {
	return fsm.State(i.state)
}

func (i *Invoice) Confirm() error {
	return nil
}

func (i *Invoice) ReceiveSignature() error {
	i.isSignatureReceived = true
	return nil
}

func (i *Invoice) Reject() error {
	return nil
}

func (i *Invoice) Approve() error {
	i.isApproved = true
	return nil
}

func (i *Invoice) Pay() error {
	return nil
}

func (i Invoice) Abandon() error {
	return nil
}


type InvoiceState fsm.State

const (
	draft InvoiceState = iota
	waitingForApproval
	waitingForsignature
	waitingForPayment
	rejected
	completed
	abandoned
)

type InvoiceCommand fsm.Command

const (
	confirm InvoiceCommand = iota
	receiveSignature
	reject
	approve
	pay
	abandon
)

func NewInvoiceStateMachine(invoice *Invoice) fsm.StateMachine {

	sm := fsm.New(invoice)

	sm.WithCommand(fsm.Command(abandon), invoice.Abandon).
		WithTransition(fsm.State(draft), fsm.State(abandoned)).
		WithTransition(fsm.State(waitingForApproval), fsm.State(abandoned)).
		WithTransition(fsm.State(waitingForsignature), fsm.State(abandoned)).
		WithTransition(fsm.State(waitingForPayment), fsm.State(abandoned))

	sm.WithCommand(fsm.Command(confirm), invoice.Confirm).
		WithTransition(fsm.State(draft), fsm.State(waitingForApproval))

	sm.WithCommand(fsm.Command(approve), invoice.Approve).
		WithConditionedTransition(fsm.State(waitingForApproval), 
			fsm.State(waitingForsignature), func() bool {
				return func(i Invoice) bool {
					return i.needsSignature && !i.isSignatureReceived
				}(*invoice)
			}).
		WithTransition(fsm.State(waitingForApproval), fsm.State(waitingForPayment))

	sm.WithCommand(fsm.Command(receiveSignature), invoice.ReceiveSignature).
		WithTransition(fsm.State(waitingForsignature), fsm.State(waitingForPayment)).
		WithTransition(fsm.State(waitingForApproval), fsm.State(waitingForApproval))

	sm.WithCommand(fsm.Command(reject), invoice.Reject).
		WithTransition(fsm.State(waitingForApproval), fsm.State(rejected))

	sm.WithCommand(fsm.Command(pay), invoice.Pay).
		WithTransition(fsm.State(waitingForPayment), fsm.State(completed))

	return sm
}
