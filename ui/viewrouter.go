package ui

import (
	"time"

	"github.com/gdamore/tcell/v3"
)

// ViewID is a unique identifier for a [View].
type ViewID int

// The different kinds of [ViewID].
const (
	ViewIDNone ViewID = iota
	ViewIDAdTree
	ViewIDTabs
	ViewIDOperations
	ViewIDLog
	ViewIDTransfers
)

// IsValid checks if the [ViewID] is valid.
func (v ViewID) IsValid() bool {
	return v != ViewIDNone
}

// R creates a router message for a [ViewID].
func (v ViewID) R[T any](t T) *RouterMsg[T] {
	return NewRouterMsg(v, t)
}

func (v ViewID) String() string {
	switch v {
	case ViewIDAdTree:
		return "Tree"

	case ViewIDTabs:
		return "Tabs"

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
	tcell.Event

	// IsValid returns if the message is a valid routed message.
	IsValid() bool

	// ViewID returns the view's ID that the message has to be routed to.
	ViewID() ViewID
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

// When implements [tcell.Event].
func (r *RouterMsg[T]) When() time.Time {
	return time.Time{}
}

var _ Router = (*RouterMsg[struct{}])(nil)
