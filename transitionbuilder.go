package fsm

type TransitionBuilder struct {
	sm        *StateMachine
	from      State
	to        State
	cmdID     CommandID
	condition Condition
}

func (t *TransitionBuilder) To(s State) *TransitionBuilder {
	t.to = s
	return t
}

func (t *TransitionBuilder) On(cmd CommandID) *TransitionBuilder {
	t.cmdID = cmd
	return t
}

func (t *TransitionBuilder) If(cond Condition) *TransitionBuilder {
	t.condition = cond
	return t
}

func (t *TransitionBuilder) Add() *TransitionBuilder {

	if _, ok := t.sm.transitions[t.from]; !ok {
		t.sm.transitions[t.from] = map[CommandID]Targets{}
	}

	if _, ok := t.sm.transitions[t.from][t.cmdID]; !ok {
		t.sm.transitions[t.from][t.cmdID] = Targets{}
	}

	t.sm.transitions[t.from][t.cmdID][t.to] = t.condition

	return &TransitionBuilder{
		sm:   t.sm,
		from: t.from,
	}
}
