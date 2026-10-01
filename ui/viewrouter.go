package ui

import (
	"github.com/ayn2op/tview"
)

// ViewID is a unique identifier for a [View].
type ViewID int

// The different kinds of [ViewID].
const (
	ViewIDNone ViewID = iota
	ViewIDAdTree
	ViewIDStatusBar
	ViewIDTabs
	ViewIDInfo
	ViewIDOperations
	ViewIDLog
)

// IsValid checks if the [ViewID] is valid.
func (v ViewID) IsValid() bool {
	return v != ViewIDNone
}

// RouterMessage creates a router message for a [ViewID].
func (v ViewID) RouterMessage(t tview.Msg) RouterMsg {
	return NewRouterMsg(v, t)
}

// RouterMsg describes a message to be routed to a [View].
type RouterMsg struct {
	ID  ViewID
	Msg tview.Msg
}

// NewRouterMsg returns a new router message, to send a message to
// a specific view, selected by it's [ViewID].
func NewRouterMsg(id ViewID, msg tview.Msg) RouterMsg {
	return RouterMsg{id, msg}
}

// EmptyRouterMsg returns an empty router message.
func EmptyRouterMsg() RouterMsg {
	return RouterMsg{}
}

// IsValid returns if the message is a valid routed message.
func (r *RouterMsg) IsValid() bool {
	return r.ID != ViewIDNone
}

// SendRoutedMsg prefixes the message with the destination view's [ViewID], and sends the message.
func (r RouterMsg) SendRoutedMsg(rv RootView) tview.Cmd {
	return rv.SendRoutedUpdateMsg(r)
}
