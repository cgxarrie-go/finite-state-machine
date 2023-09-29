package fsm

import (
	"fmt"
	"reflect"
)

type State uint32

type Command uint32

type HandledElement interface {
	SetState(State)
	State() State
}

type StateMachine struct {
	element HandledElement
	actions map[Command]*Action
}

type Action struct {
	command     Command
	funcName    string
	transitions map[State]*Transition
}

type Transition struct {
	From    State
	Targets map[State]*Target
}

type Target struct {
	To       State
	funcName string
}

func New(element HandledElement) StateMachine {
	fsm := &StateMachine{
		element: element,
	}

	return *fsm
}

func (fsm *StateMachine) WithCommand(command Command, methodName string) *Action {
	if fsm.actions == nil {
		fsm.actions = make(map[Command]*Action)
	}

	if action, ok := fsm.actions[command]; ok {
		return action
	}

	fsm.actions[command] = &Action{
		funcName: methodName,
		command:  command,
	}

	action := fsm.actions[command]
	return action
}

func (a *Action) WithTransition(from, to State) *Action {
	return a.WithConditionedTransition(from, to, "")
}

func (a *Action) WithConditionedTransition(from, to State, conditionFuncName string) *Action {
	if a.transitions == nil {
		a.transitions = make(map[State]*Transition)
	}

	if _, ok := a.transitions[from]; !ok {
		a.transitions[from] = &Transition{
			From: from,
			Targets: map[State]*Target{
				to: {
					To:       to,
					funcName: conditionFuncName,
				},
			},
		}
		return a
	}

	if _, ok := a.transitions[from].Targets[to]; !ok {
		a.transitions[from].Targets[to] = &Target{
			To:       to,
			funcName: conditionFuncName,
		}
		return a
	}

	return a
}

func (fsm *StateMachine) ExecuteCommand(command Command) error {

	if fsm.actions == nil {
		return fmt.Errorf("cannot execute requested command %v from state %v",
			command, fsm.element.State())
	}

	action, ok := fsm.actions[command]
	if !ok {
		return fmt.Errorf("cannot find action for command %v", command)
	}

	if action.transitions == nil {
		return fmt.Errorf("cannot find transitions for command %v",
			command)
	}

	transition, ok := action.transitions[fsm.element.State()]
	if !ok {
		return fmt.Errorf("cannot find transitions for command %v "+
			"and state %v", command, fsm.element)
	}

	meth := reflect.ValueOf(fsm.element).MethodByName(action.funcName)
	if !fsm.isValidCommandFunc(meth.Interface()) {
		return fmt.Errorf("method %v for command %v has wrong signature. "+
			"It should be func() error", action.funcName, command)
	}
	err := fsm.runCommandFunc(meth)
	if err != nil {
		return fmt.Errorf("method %v for command %v returned error: %v",
			action.funcName, command, err)
	}

	for _, target := range transition.Targets {
		if target.funcName != "" {
			meth := reflect.ValueOf(fsm.element).MethodByName(target.funcName)
			if !fsm.isValidTransitionConditionFunc(meth.Interface()) {
				return fmt.Errorf("method %v for command %v has wrong "+
					"signature. It should be func() bool", target.funcName,
					command)
			}
			if !fsm.runTransitionConditionFunc(meth) {
				continue
			}

			fsm.element.SetState(target.To)
			return nil
		}

		fsm.element.SetState(target.To)
		return nil
	}

	return fmt.Errorf("cannot find executable transition for command %v "+
		"and state %v", command, fsm.element.State())
}

func (fsm StateMachine) isValidCommandFunc(fn interface{}) bool {
	fnType := reflect.TypeOf(fn)

	if fnType.NumIn() != 0 {
		return false
	}

	if fnType.NumOut() != 1 {
		return false
	}

	return fnType.Out(0).AssignableTo(reflect.TypeOf((*error)(nil)).Elem())
}

func (fsm StateMachine) runCommandFunc(meth reflect.Value) error {
	methResult := meth.Call(nil)
	errorElement := reflect.TypeOf((*error)(nil)).Elem()
	if len(methResult) > 0 && !methResult[0].IsNil() &&
		methResult[0].Type().Implements(errorElement) {
		return methResult[0].Interface().(error)
	}
	return nil
}

func (fsm StateMachine) isValidTransitionConditionFunc(fn interface{}) bool {
	fnType := reflect.TypeOf(fn)

	if fnType.NumIn() != 0 {
		return false
	}

	if fnType.NumOut() != 1 {
		return false
	}

	return fnType.Out(0).AssignableTo(reflect.TypeOf((*bool)(nil)).Elem())
}

func (fsm StateMachine) runTransitionConditionFunc(meth reflect.Value) bool {
	methResult := meth.Call(nil)
	if len(methResult) > 0 &&
		methResult[0].Kind() == reflect.Bool {
		return methResult[0].Bool()
	}
	return false
}
