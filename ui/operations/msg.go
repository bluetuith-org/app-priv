package operations

import (
	"github.com/bluetuith-org/bluetuith/ui"
)

// CreateMsg holds a operation creation message.
type CreateMsg struct {
	CreationInfo

	opAction Invoker
}

// MsgCreate creates a message to be sent to the operations manager, to create an operation.
func MsgCreate(creationInfo CreationInfo, action Invoker) ui.RouterMsg {
	return ui.ViewIDOperations.RouterMessage(CreateMsg{CreationInfo: creationInfo, opAction: action})
}

type opUpdateMsg struct {
	id, message string
}

func msgOpUpdate(id, msg string) ui.RouterMsg {
	return ui.ViewIDOperations.RouterMessage(opUpdateMsg{id, msg})
}

type opDeleteMsg struct {
	id string
}

func msgOpDelete(id string) ui.RouterMsg {
	return ui.ViewIDOperations.RouterMessage(opDeleteMsg{id})
}

type opErrorMsg struct {
	err error
}

func msgOpError(err error) ui.RouterMsg {
	return ui.ViewIDOperations.RouterMessage(opErrorMsg{err})
}
