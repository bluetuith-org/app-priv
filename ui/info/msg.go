package info

import (
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/ui"
)

type infoMsg struct {
	id      string
	content []text.Span
}

type resetMsg struct{}

// Msg returns a message to print information.
func Msg(id string, content []text.Span) ui.Router {
	return ui.ViewIDInfo.RouterMessage(infoMsg{id, content})
}

// ResetMsg sends a message to reset the information view.
func ResetMsg() ui.Router {
	return ui.ViewIDInfo.RouterMessage(resetMsg{})
}
