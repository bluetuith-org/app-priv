package popup

import (
	"time"

	"github.com/gdamore/tcell/v3"
)

// R is an alias for [*RouterMsg].
type R[T any] = *RouterMsg[T]

// Router describes a popup message router, that routes a message to a popup,
// based on it's popup ID.
type Router interface {
	tcell.Event

	// IsValid returns if the message is a valid routed message.
	IsValid() bool

	// PopupID returns the ID of the popup the message should be routed to.
	PopupID() string
}

// RouterMsg describes a message that is to be routed to a popup.
type RouterMsg[T any] struct {
	id string

	Value T
}

// NewRouterMsg returns a new routed popup message.
func NewRouterMsg[T any](id string, v T) *RouterMsg[T] {
	return &RouterMsg[T]{id, v}
}

// CreateMsg returns a message to create a new popup.
func CreateMsg(popup Overlay) *RouterMsg[popupModifyMsg] {
	return NewRouterMsg(popup.PopupID(), popupModifyMsg{popup})
}

// RemoveMsg returns a message to close a popup based on it's ID.
func RemoveMsg(id string) *RouterMsg[popupRemoveMsg] {
	return NewRouterMsg(id, popupRemoveMsg{})
}

// PopupID implements [Router].
func (p *RouterMsg[T]) PopupID() string {
	return p.id
}

// IsValid implements [Router].
func (p *RouterMsg[T]) IsValid() bool {
	return p.id != ""
}

// When implements [Router].
func (p *RouterMsg[T]) When() time.Time {
	return time.Time{}
}

func newDeleteMsg() *RouterMsg[popupDeleteMsg] {
	return NewRouterMsg("", popupDeleteMsg{})
}

type (
	popupModifyMsg struct {
		popup Overlay
	}

	popupDeleteMsg struct{}
	popupRemoveMsg struct{}
)

var _ Router = (*RouterMsg[struct{}])(nil)
