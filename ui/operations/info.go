package operations

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

type operationStage uint8

const (
	operationInProgress operationStage = iota
	operationCompleted
	operationError
)

func (o operationStage) Format() (string, tcell.Style) {
	switch o {
	case operationCompleted:
		return theme.Icons().CheckMark.String() + " Completed", theme.Current().Operations.Completed

	case operationError:
		return theme.Icons().Error.String() + " Error", theme.Current().Operations.Error
	}

	return theme.Icons().Progress.String() + " In Progress", theme.Current().Operations.InProgress
}

// Invoker describes a function that is executed when the operation is running.
type Invoker func(ov *RunningInfo) tview.Msg

// CreationInfo holds the information for the operation to be executed.
type CreationInfo struct {
	id string

	header, message string
	stage           operationStage
	err             error

	description stext.TextPairs
}

// NewCreationInfo creates an initial operation context, to be scheduled for execution.
func NewCreationInfo(desc stext.TextPairs) CreationInfo {
	return CreationInfo{description: desc}
}

// UpdateID updates the ID of the operation being created.
func (o *CreationInfo) UpdateID(id string) {
	o.id = id
}

func (o *CreationInfo) fromUpdateMsg(msg opUpdateMsg) {
	if o.id != msg.id || msg.message == "" && msg.err == nil {
		return
	}

	o.message = msg.message
	o.stage = msg.stage

	o.err = msg.err
}

func (o *CreationInfo) buildSegments() (header, msg, status richtext.Segment) {
	header = richtext.NewSegment(o.header, theme.Current().Operations.OperationNumber)

	switch o.stage {
	case operationError:
		msg = richtext.NewSegment(o.err.Error(), theme.Current().Operations.Error)

	default:
		msg = richtext.NewSegment(o.message, theme.Current().Operations.Message)
	}

	status = richtext.NewSegment(o.stage.Format())

	return
}

func (o CreationInfo) updateHeader(header string) CreationInfo {
	o.header = header
	return o
}

// RunningInfo holds the current state of the running operation.
type RunningInfo struct {
	CreationInfo

	v ui.RootView
}

// newRunningInfo creates a new executing operation context.
func newRunningInfo(v ui.RootView, creationInfo CreationInfo) *RunningInfo {
	return &RunningInfo{CreationInfo: creationInfo, v: v}
}

// Info sets the information for the current operation.
func (o *RunningInfo) Info(msg string) {
	o.v.SendMsg(msgOpUpdate(o.CreationInfo.id, msg, operationInProgress))
}

// Ok sets the information for the completed operation.
func (o *RunningInfo) Ok(msg string) ui.RouterMsg {
	return msgOpUpdate(o.CreationInfo.id, msg, operationCompleted)
}

// Error sends an error to the operations manager.
func (o *RunningInfo) Error(msg string, err error) ui.RouterMsg {
	return msgOpUpdateErr(o.id, msg, err)
}
