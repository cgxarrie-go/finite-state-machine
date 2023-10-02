# Go - Fluent State Machine
A fluent state machine which provides an easy to configure and use state mahine

## Features
- Fluent addition of transitions
- Transitions based on current status and requested action
- Execute non available transition returns error

## How to use
- Declare the object to be handled by the state machine 
	- must implementing fsm.SMObject
	- Action methods must be func() err
- Declare the collection of possible status to be handled by the machine, of type fsm.State
- Declare the collection of possible commands to be executed by the machine, of type fsm.CommandID
- Declare the machine of type fsm.StateMachine
- Declare transitions in constructor of the machine

### Declaration of States
```go
	type MyState fsm.State

	const (
		state1 MyState = iota
		state2
		state3
		state4
		state5
		state6
		state7
	)
```

### Declaration of Commands
```go
	type MyCommand fsm.State

	const (
		command1 MyCommand = iota
		command2
		command3
		command4
		command5
	)
```

### Declaration of transitions
Valid transitions are added to the state machine on constructor via the following fluent command
```go
        sm.From(state). // when the machine is in this state
            On(command). // On executing this command
            If(condition). // If after executing command this condition is met
            To(state). // then change to tis state
            Add() //Add the transition to the state machine
```


## Example : State-Machine
We will simulate an Invoice workflow, declaring the following statuses
- Draft
- Waiting for approval
- Waiting for signature
- Waiting for payment
- Rejected
- Completed
- Abandoned

the folloing actions
- Abandor
- Confirm
- Approve
- Reject
- Pay

and the following use cases
- In satus Draft: 
	- On Abandon, change to Abandoned
	- On Confirm, change to Waiting for approval
- In status Waiting for approval:
	- On Abandon, change to Abandoned
	- On Approve, If signature is needed and not received change to Waiting for signature
	- On Approve, If signature is needed and received change to Waiting for payment
	- On Approve, If signature is not needed change to Waiting for payment
	- On ReceiveSignature, stay in status Waiting ForApproval
	- On Reject, change to Rejected
- In status Waiting for Signature:
	- On Abandon, change to Abandoned
	- On ReceiveSignature, change to Waiting for payment
- In status Waiting for payment:
	- On Abandon, change to Abandoned
	- On pay, change to completed

```go

// import fsm
import (
	"github.com/cgxarrie/fsm-go/fsm"
)

// declare the object to be handled by the state machine implementing fsm.SMObject
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

// declare constructor and action methods

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

// declare available states
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

// declare available commands
type InvoiceCommand fsm.CommandID
const (
	abandon InvoiceCommand = iota
	confirm
	approve
	receiveSignature
	reject
	pay
)

// declare state machine constructor
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

// Use the state machine
func main (){

	inv := NewInvoice(false)
	inv.SetState(fsm.State(draft))
	
    sm := NewInvoiceStateMachine(&inv)
    err := sm.Do(fsm.Command(confirm))

    if  err != nil{
        // the transition does not exist.
        // cannot execute specified command in the current state.
        // State machine state not changed
    }

	// after the confirm action, invoice has changed from
	// draft to waiting for approval
	
}

```
