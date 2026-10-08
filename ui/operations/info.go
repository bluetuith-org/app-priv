package operations

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
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
	ct *content
}

// NewCreationInfo creates an initial operation context, to be scheduled for execution.
func NewCreationInfo(desc card.ContentProvider) CreationInfo {
	return CreationInfo{ct: newContent(desc)}
}

// UpdateID updates the ID of the operation being created.
func (o *CreationInfo) UpdateID(id string) {
	o.ct.setID(id)
}

func (o *CreationInfo) getID() string {
	return o.ct.id
}

func (o *CreationInfo) fromUpdateMsg(msg opUpdateMsg) {
	o.ct.fromUpdateMsg(msg)
}

func (o *CreationInfo) updateHeader(header string) *CreationInfo {
	o.ct.setHeader(header)

	return o
}

func (o *CreationInfo) getContent() *card.Content {
	return o.ct.cc
}

// RunningInfo holds the current state of the running operation.
type RunningInfo struct {
	*content

	v ui.RootView
}

// newRunningInfo creates a new executing operation context.
func newRunningInfo(v ui.RootView, creationInfo *CreationInfo) *RunningInfo {
	return &RunningInfo{content: creationInfo.ct, v: v}
}

// Info sets the information for the current operation.
func (o *RunningInfo) Info(msg string) {
	o.v.SendMsg(msgOpUpdate(o.id, msg, operationInProgress))
}

// Ok sets the information for the completed operation.
func (o *RunningInfo) Ok(msg string) ui.RouterMsg {
	return msgOpUpdate(o.id, msg, operationCompleted)
}

// Error sends an error to the operations manager.
func (o *RunningInfo) Error(msg string, err error) ui.RouterMsg {
	return msgOpUpdateErr(o.id, msg, err)
}
