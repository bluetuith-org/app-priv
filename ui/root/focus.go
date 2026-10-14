package root

import "github.com/bluetuith-org/bluetuith/ui"

type focusState struct {
	focus ui.ViewID
}

func newFocusState(id ui.ViewID) *focusState {
	return &focusState{id}
}

func (f *focusState) SetFocus(id ui.ViewID) {
	f.focus = id
}

func (f *focusState) Focused() ui.ViewID {
	return f.focus
}
