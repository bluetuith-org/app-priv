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
func MsgCreate(creationInfo CreationInfo, action Invoker) ui.Router {
	return ui.ViewIDOperations.R(CreateMsg{CreationInfo: creationInfo, opAction: action})
}

type opUpdateMsg struct {
	id      uint64
	message string
	err     error
	stage   operationStage
}

func msgOpUpdate(id uint64, msg string, stage operationStage) ui.Router {
	return ui.ViewIDOperations.R(opUpdateMsg{id, msg, nil, stage})
}

func msgOpUpdateErr(id uint64, msg string, err error) ui.Router {
	return ui.ViewIDOperations.R(opUpdateMsg{id, msg, err, operationError})
}
