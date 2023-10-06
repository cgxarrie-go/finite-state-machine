package invoiceFsm

import "github.com/cgxarrie-go/fsm"

type Invoice struct {
	state             InvoiceState
	isSignatureReceived bool
	isApproved          bool
	needsSignature    bool
}

func (i *Invoice) SetState(state fsm.State) {
	i.state = InvoiceState(state)
}

func (i *Invoice) State() fsm.State {
	return fsm.State(i.state)
}

func NewInvoice(needsSignature bool) Invoice {
	return Invoice{
		state:             draft,
		isSignatureReceived: false,
		isApproved:          false,
		needsSignature:    needsSignature,
	}
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

type InvoiceCommand fsm.CommandID
const (
	abandon InvoiceCommand = iota
	confirm
	approve
	receiveSignature
	reject
	pay
)


func NewInvoiceStateMachine(invoice *Invoice) fsm.StateMachine {

	sm := fsm.New(invoice)
	sm.
		WithCommand(fsm.CommandID(abandon), invoice.Abandon).
		WithCommand(fsm.CommandID(confirm), invoice.Confirm).
		WithCommand(fsm.CommandID(approve), invoice.Approve).
		WithCommand(fsm.CommandID(receiveSignature), invoice.ReceiveSignature).
		WithCommand(fsm.CommandID(reject), invoice.Reject).
		WithCommand(fsm.CommandID(pay), invoice.Pay)


	sm.From(fsm.State(draft)).
		On(fsm.CommandID(abandon)).To(fsm.State(abandoned)).Add().
		On(fsm.CommandID(confirm)).To(fsm.State(waitingForApproval)).Add()
	
	needsSignature :=func() bool {
		return func(i Invoice) bool {
			return i.needsSignature && !i.isSignatureReceived
		}(*invoice)
	}

	sm.From(fsm.State(waitingForApproval)).
		On(fsm.CommandID(abandon)).To(fsm.State(abandoned)).Add().
		On(fsm.CommandID(approve)).If(needsSignature).To(fsm.State(waitingForsignature)).Add().
		On(fsm.CommandID(approve)).To(fsm.State(waitingForPayment)).Add().
		On(fsm.CommandID(receiveSignature)).To(fsm.State(waitingForApproval)).Add().
		On(fsm.CommandID(reject)).To(fsm.State(rejected)).Add()

	sm.From(fsm.State(waitingForsignature)).
		On(fsm.CommandID(abandon)).To(fsm.State(abandoned)).Add().
		On(fsm.CommandID(receiveSignature)).To(fsm.State(waitingForPayment)).Add()

	sm.From(fsm.State(waitingForPayment)).
		On(fsm.CommandID(abandon)).To(fsm.State(abandoned)).Add().
		On(fsm.CommandID(pay)).To(fsm.State(completed)).Add()

	return sm
}
