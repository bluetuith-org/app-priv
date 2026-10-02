package operations

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/ui"
)

// Invoker describes a function that is executed when the operation is running.
type Invoker func(ov *RunningInfo) tview.Msg

// CreationInfo holds the information for the operation to be executed.
type CreationInfo struct {
	id string

	description richtext.Text
	message     string
}

// NewCreationInfo creates an initial operation context, to be scheduled for execution.
func NewCreationInfo(desc richtext.Text) CreationInfo {
	return CreationInfo{description: desc}
}

// UpdateID updates the ID of the operation being created.
func (o *CreationInfo) UpdateID(id string) {
	o.id = id
}

// RunningInfo holds the current state of the running operation.
type RunningInfo struct {
	CreationInfo

	v ui.RootView
}

// NewRunningInfo creates a new executing operation context.
func NewRunningInfo(v ui.RootView, creationInfo CreationInfo) *RunningInfo {
	return &RunningInfo{CreationInfo: creationInfo, v: v}
}

// Info sets the information for the current operation.
func (o *RunningInfo) Info(msg string) {
	o.v.SendMsg(msgOpUpdate(o.CreationInfo.id, msg))
}

// Error sends an error to the operations manager.
func (o *RunningInfo) Error(err error) ui.RouterMsg {
	return msgOpError(err)
}
