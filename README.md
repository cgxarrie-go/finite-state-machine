# Go - Finite State Machine
A finite state machine which provides an easy to configure and use state mahine

## Features
- Fluent addition of transitions
- Transitions based on current status and requested action
- Execute non available transition returns error

## How to use
- Declare the collection of possible status to be handled by the machine, of type fsm.State
- Declare the collection of possible commands to be executed by the machine, of type fsm.Command
- Declare the machine of type fsm.StateMachine
- Declare transitions in constructor of the machine

### Declaration of transitions
Valid transitions are added to the state machine on constructor via the following fluent command
```go
        WithTransition().
            On(command). // On executing this command
            From(state). // when the machine is in this state
            To(state). // then change to tis state
            Build() //Build the transition and add it to the state machine
```


## Example : State-Machine
We will simulate an Invoice workflow, declaring the following statuses
- Draft
- Waiting for approval
- Waiting for payment
- Rejected
- Completed

the folloing actions
- Confirm
- Approve
- Reject
- Pay

and the following use cases
- In satus Draft, On Confirm, change to Waiting for approval
- In status Waiting for approval, On Approve, change to Waiting for payment
- In status Waiting for approval, On Reject, change to Rejected
- In status Waiting for payment, On Pay, change to Completed


```go

// import fsm
import (
	"github.com/cgxarrie/fsm-go/fsm"
)

// declare available states
type InvoiceState fsm.State

const (
	draft InvoiceState = iota
	waitingForApproval
	waitingForPayment
	rejected
	completed
)

// declare available commands
type InvoiceCommand fsm.Command

const (
	confirm InvoiceCommand = iota
	reject
	approve
	pay
)

// declare state machine constructor
func NewInvoiceStateMachine() fsm.StateMachine {

    // Create state machine with starting state
	sm := fsm.New(fsm.State(draft))

	sm.WithTransition().                    // add transition
		On(fsm.Command(confirm)).           // On executing command Confirm
		From(fsm.State(draft)).             // from state draft
		To(fsm.State(waitingForApproval)).  // change to waitingForApproval
		Build()

	sm.WithTransition().
		On(fsm.Command(reject)).
		From(fsm.State(waitingForApproval)).
		To(fsm.State(rejected)).
		Build()

	sm.WithTransition().
		On(fsm.Command(approve)).
		From(fsm.State(waitingForApproval)).
		To(fsm.State(waitingForPayment)).
		Build()

	sm.WithTransition().
		On(fsm.Command(pay)).
		From(fsm.State(waitingForPayment)).
		To(fsm.State(completed)).
		Build()

	return sm
}

func main (){

    sm := NewInvoiceStateMachine()

    _, err := sm.ExecuteCommand(fsm.Command(confirm))

    if  err != nil{
        // the transition does not exist.
        // cannot execute specified command in the current state.
        //State machine state not changed
    }
}

```
