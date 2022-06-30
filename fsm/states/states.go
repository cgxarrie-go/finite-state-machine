package states

type State uint32

const (
	Locked State = iota
	Unlocked
)
