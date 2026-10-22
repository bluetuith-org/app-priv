package info

import (
	"github.com/bluetuith-org/bluetuith/ui"
)

type infoMsg struct {
	content Content
}

type resetMsg struct{}

// Msg returns a message to print information.
func Msg(content Content) ui.Router {
	return ui.ViewIDInfo.R(infoMsg{content})
}

// ResetMsg sends a message to reset the information view.
func ResetMsg() ui.Router {
	return ui.ViewIDInfo.R(resetMsg{})
}
