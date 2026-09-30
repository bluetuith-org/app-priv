package operations

import (
	"github.com/bluetuith-org/bluetuith/ui"
)

// OpCreateMsg holds a operation creation message.
type OpCreateMsg struct {
	CreationInfo

	opAction Invoker
}

// MsgCreate creates a message to be sent to the operations manager, to create an operation.
func MsgCreate(creationInfo CreationInfo, action Invoker) ui.RouterMsg {
	return ui.ViewIDOperations.RouterMessage(OpCreateMsg{CreationInfo: creationInfo, opAction: action})
}

type opUpdateMsg struct {
	CreationInfo
}

func msgOpUpdate(creationInfo CreationInfo, desc, msg string) ui.RouterMsg {
	creationInfo.description = desc
	creationInfo.message = msg

	return ui.ViewIDOperations.RouterMessage(opUpdateMsg{CreationInfo: creationInfo})
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
