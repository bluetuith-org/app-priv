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
	ViewIDTabs
	ViewIDInfo
	ViewIDOperations
	ViewIDLog
	ViewIDTransfers
)

// IsValid checks if the [ViewID] is valid.
func (v ViewID) IsValid() bool {
	return v != ViewIDNone
}

// RouterMessage creates a router message for a [ViewID].
func (v ViewID) RouterMessage[T any](t T) *RouterMsg[T] {
	return NewRouterMsg(v, t)
}

func (v ViewID) String() string {
	switch v {
	case ViewIDAdTree:
		return "Tree"

	case ViewIDTabs:
		return "Tabs"

	case ViewIDInfo:
		return "Info"

	case ViewIDOperations:
		return "Operations"

	case ViewIDLog:
		return "Logs"

	case ViewIDTransfers:
		return "Transfers"
	}

	return ""
}

// R is an alias for [*RouterMsg].
// It is useful to reduce verbosity when matching routed messages.
type R[T any] = *RouterMsg[T]

// Router describes a message router, that routes a particular message to the view,
// based on its view ID.
type Router interface {
	// IsValid returns if the message is a valid routed message.
	IsValid() bool

	// ViewID returns the view's ID that the message has to be routed to.
	ViewID() ViewID

	// SendRoutedMsg prefixes the message with the destination view's [ViewID], and sends the message.
	SendRoutedMsg(rv RootView) tview.Cmd
}

// RouterMsg describes a message to be routed to a [View].
type RouterMsg[T any] struct {
	ID    ViewID
	Value T
}

// NewRouterMsg returns a new router message, to send a message to
// a specific view, selected by it's [ViewID].
func NewRouterMsg[T any](id ViewID, value T) *RouterMsg[T] {
	return &RouterMsg[T]{id, value}
}

// EmptyRouterMsg returns an empty router message.
func EmptyRouterMsg() Router {
	return nil
}

// IsValid returns if the message is a valid routed message.
func (r *RouterMsg[T]) IsValid() bool {
	return r.ID != ViewIDNone
}

// ViewID returns the view's ID that the message has to be routed to.
func (r *RouterMsg[T]) ViewID() ViewID {
	return r.ID
}

// SendRoutedMsg prefixes the message with the destination view's [ViewID], and sends the message.
func (r *RouterMsg[T]) SendRoutedMsg(rv RootView) tview.Cmd {
	return rv.SendRoutedUpdateMsg(r)
}

var _ Router = (*RouterMsg[struct{}])(nil)
