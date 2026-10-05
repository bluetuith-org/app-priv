package tabs

// Action is what a key makes tabs do.
type Action int

// The different types of actions.
const (
	ActionNone Action = iota
	ActionPrevious
	ActionNext
)
