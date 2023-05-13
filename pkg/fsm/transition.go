package fsm

type transition map[Command]map[State]State

type TransitionBuilder struct {
	onCommand Command
	from      State
	to        State
	sm        *StateMachine
}

func (tb *TransitionBuilder) On(c Command) *TransitionBuilder {
	tb.onCommand = c
	return tb
}

func (tb *TransitionBuilder) From(s State) *TransitionBuilder {
	tb.from = s
	return tb
}

func (tb *TransitionBuilder) To(s State) *TransitionBuilder {
	tb.to = s
	return tb
}

func (tb TransitionBuilder) Build() {

	if tb.sm.Transitions == nil {
		tb.sm.Transitions = transition{
			tb.onCommand: {
				tb.from: tb.to,
			},
		}
	}

	_, commandExists := tb.sm.Transitions[tb.onCommand]

	if !commandExists {
		tb.sm.Transitions[tb.onCommand] = map[State]State{tb.from: tb.to}
	}

	tb.sm.Transitions[tb.onCommand][tb.from] = tb.to
}
