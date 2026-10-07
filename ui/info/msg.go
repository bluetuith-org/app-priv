package info

import (
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/ui"
)

type infoMsg struct {
	id      string
	content richtext.Text
}

type resetMsg struct{}

// Msg returns a message to print information.
func Msg(id string, content richtext.Text) ui.RouterMsg {
	return ui.ViewIDInfo.RouterMessage(infoMsg{id, content})
}

// ResetMsg sends a message to reset the information view.
func ResetMsg() ui.RouterMsg {
	return ui.ViewIDInfo.RouterMessage(resetMsg{})
}
