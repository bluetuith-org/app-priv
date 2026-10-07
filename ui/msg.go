package ui

import "github.com/ayn2op/tview"

// FocusMsg holds a [ViewID] that was set to the focused state.
type FocusMsg struct {
	ID ViewID
}

// FocusViewCmd returns a [tview.Cmd] to send a [FocusMsg].
func FocusViewCmd(id ViewID) tview.Cmd {
	return func() tview.Msg {
		return FocusMsg{ID: id}
	}
}
