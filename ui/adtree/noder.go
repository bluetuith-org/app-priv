package adtree

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/info"
	"github.com/gdamore/tcell/v3"
)

// adNoder describes an interface to get information from nodes and update node states.
type adNoder interface {
	info.Content

	// AdNode returns the base node data.
	AdNode() *adNode

	// HandleKeys handles a [tview.KeyMsg] and returns a message.
	HandleKeys(p tview.KeyMsg) (ui.Router, bool)

	// NodeStyle returns the style to be applied for this node's label.
	NodeStyle() tcell.Style

	// Icon returns the Icon associated with this node.
	Icon() string

	// RootID returns the root ID for the subnodes of an adapter or device node.
	RootID() string

	// PopulateActions populates all actions within a node of type [nodeTypeAction].
	PopulateActions()

	// Refresh refreshes the content of the node.
	Refresh()

	// UpdateActionNode updates the node of type [nodeTypeAction] with the message.
	UpdateActionNode(actionNode *acNode, updateMsg actionUpdateMsg)

	// SetAdapterEventData sets the adapter event data for the node.
	SetAdapterEventData(ev bluetooth.AdapterEventData)

	// SetDeviceEventData sets the device event data for the node.
	SetDeviceEventData(ev bluetooth.DeviceEventData)
}

type emptyNoder struct{}

// AdNode returns the base node data.
func (e *emptyNoder) AdNode() *adNode {
	return nil
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (e *emptyNoder) HandleKeys(tview.KeyMsg) (ui.Router, bool) {
	return ui.EmptyRouterMsg(), false
}

// NodeStyle returns the style to be applied for this node's label.
func (e *emptyNoder) NodeStyle() tcell.Style {
	return tcell.Style{}
}

// Icon returns the Icon associated with this node.
func (e *emptyNoder) Icon() string {
	return ""
}

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (e *emptyNoder) PopulateActions() {
}

// RootID implements [adNoder].
func (e *emptyNoder) RootID() string {
	return ""
}

// Information implements [adNoder].
func (e *emptyNoder) Information() []text.Span {
	return nil
}

// Refresh implements [adNoder].
func (e *emptyNoder) Refresh() {
}

// UpdateActionNode updates the node of type [nodeTypeAction] with the message.
func (e *emptyNoder) UpdateActionNode(*acNode, actionUpdateMsg) {
}

// SetAdapterEventData sets the adapter event data for the node.
func (e *emptyNoder) SetAdapterEventData(bluetooth.AdapterEventData) {
}

// SetDeviceEventData sets the device event data for the node.
func (e *emptyNoder) SetDeviceEventData(bluetooth.DeviceEventData) {
}

var _ adNoder = (*emptyNoder)(nil)
