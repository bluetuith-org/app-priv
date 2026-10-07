package ui

import "github.com/ayn2op/tview"

// KeybindHandler represents a handle to handle keys and returns actions.
type KeybindHandler[A any] interface {
	// Keybinds returns an action for the matched keybindings.
	Keybinds(msg tview.KeyMsg) (A, bool)
}

// ActionChanger represents a handler state changes.
type ActionChanger[C any] interface {
	// OnChange processes a change in state and returns a message.
	OnChange(change C) tview.Msg
}

// ActionSelector represents a handler to handle selection events.
type ActionSelector[S any] interface {
	// OnSelect processes a selection event and returns a message.
	OnSelect(selected S) tview.Msg
}
