package adtree

import (
	"strings"

	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
	"github.com/gdamore/tcell/v3"
)

// adapterNode describes an adapter node.
type adapterNode struct {
	*adNode
	emptyNoder

	adapter bluetooth.AdapterData
}

// newAdapterNode creates a new adapter node.
func newAdapterNode(rv ui.RootView, rootNode *rootNode, adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) *adNode {
	an := &adapterNode{}

	adapterNode := newAdNode(
		nodeTypeAdapter, getAdapterDisplayName(adapter),
		true, newAdapterNodeID(adapter.AdapterAddress),
		an, rv,
	)

	an.adNode = adapterNode
	an.adapter = adapter

	an.AddActionsList()
	an.AddDevicesList(devices)

	an.PopulateActions()

	rootNode.AddChild(adapterNode.Node)

	return adapterNode
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (a *adapterNode) HandleKeys(p tview.KeyMsg) (ui.RouterMsg, bool) {
	return actionHandleKeyMsg(p, a.Node)
}

// NodeStyle returns the style to be applied for this node's label.
func (a *adapterNode) NodeStyle() tcell.Style {
	return theme.Current().ADTree.Adapter
}

// Icon returns the Icon associated with this node.
func (a *adapterNode) Icon() string {
	return theme.Icons().Adapter.String()
}

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (a *adapterNode) PopulateActions() {
	actionsListNode := a.Children()[relPosActionsListNode]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data[*adNode]()
	if node == nil {
		return
	}

	node.AddAction(kb().Adapter.TogglePower, boolToActionState(a.adapter.Powered.Value()), true, a.actionPowered)
	node.AddAction(kb().Adapter.ToggleDiscoverable, boolToActionState(a.adapter.Discovering.Value()), true, a.actionDiscoverable)
	node.AddAction(kb().Adapter.TogglePairable, boolToActionState(a.adapter.Pairable.Value()), true, a.actionPairable)
	node.AddAction(kb().Adapter.ToggleScan, boolToActionState(a.adapter.Discovering.Value()), true, a.actionScan)

	for _, actionNode := range actionsListNode.Children() {
		a.UpdateActionNode(actionNode.Data[*adNode](), emptyActionUpdateMsg())
	}
}

// UpdateActionNode updates the node of type [nodeTypeAction] with the message.
func (a *adapterNode) UpdateActionNode(actionNode *adNode, updateMsg actionUpdateMsg) {
	state := actionNode.actionState
	updateMsg.updateState(actionNode.ID(), state)

	var text string

	switch state.key {
	case kb().Adapter.TogglePower:
		text = "Power"

	case kb().Adapter.ToggleDiscoverable:
		text = "Discoverable"

	case kb().Adapter.TogglePairable:
		text = "Pairable"

	case kb().Adapter.ToggleScan:
		text = "Device Scanning"

	default:
		return
	}

	actionText := "Off"
	if state.currentState == actionStateDisabled {
		actionText = "On"
	}

	displayName := useStringBuffer(len(text)+len(actionText)+10, func(b *strings.Builder) {
		b.WriteString("Switch ")

		b.WriteString(text)
		b.WriteString(" ")

		b.WriteString(actionText)
	})

	actionNode.SetName(displayName)
}

// SetAdapterEventData sets the adapter event data for the node.
func (a *adapterNode) SetAdapterEventData(ev bluetooth.AdapterEventData) {
	a.adapter.AdapterEventData = ev
}

func (a *adapterNode) actionPowered() (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (a *adapterNode) actionDiscoverable() (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (a *adapterNode) actionPairable() (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (a *adapterNode) actionScan() (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

var _ adNoder = (*adapterNode)(nil)
