package adtree

import (
	"errors"
	"strings"
	"time"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetooth-classic/api/optional"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
	"github.com/gdamore/tcell/v3"
)

// adapterNode describes an adapter node.
type adapterNode struct {
	*adNode

	stn     *stNode
	adapter bluetooth.AdapterData
}

// newAdapterNode creates a new adapter node.
func newAdapterNode(rv ui.RootView, rootNode *rootNode, adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) *adNode {
	an := &adapterNode{}

	adapterNode := newAdNode(
		nodeTypeAdapter,
		true, 5,
		newAdapterNodeID(adapter.AdapterAddress),
		an, rv,
	)

	an.adNode = adapterNode
	an.adapter = adapter
	an.stn = newStatusNode(adapterNode, true)

	newActionsListNode(an.adNode)
	newDeviceListNode(adapterNode, devices)

	rootNode.AddChild(adapterNode.Node)

	an.Refresh()
	an.PopulateActions()

	return adapterNode
}

// Refresh refreshes the content of the node.
func (a *adapterNode) Refresh() {
	a.stn.SetContent(a.StatusIcons())

	a.SetContent(buildLabel(a.Icon(), getAdapterDisplayName(a.adapter), a.NodeStyle()))
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (a *adapterNode) HandleKeys(p tview.KeyMsg) (ui.RouterMsg, bool) {
	return actionHandleKeyMsg(p, a.Node)
}

// NodeStyle returns the style to be applied for this node's label.
func (a *adapterNode) NodeStyle() tcell.Style {
	return theme.Current().ADTree.Adapter.Name
}

// Icon returns the Icon associated with this node.
func (a *adapterNode) Icon() string {
	return theme.Icons().Adapter.String()
}

// RootID returns the root ID for the subnodes of an adapter or device node.
func (a *adapterNode) RootID() string {
	return a.id.String()
}

// SubnodePosition returns the position of the subnode within the parent node.
func (a *adapterNode) SubnodePosition(relpos adSubNodePos) (int, bool) {
	switch relpos {
	case relPosStatusNode:
		return 0, true

	case relPosActionsListNode:
		return 1, true

	case relPosDevicesListNode:
		return 2, true
	}

	return -1, false
}

// Information returns a message to print adapter/device information.
func (a *adapterNode) Information() (string, richtext.Text) {
	return a.id.String(), nil
}

// StatusIcons returns the status icon pills for this node.
func (a *adapterNode) StatusIcons() richtext.Line {
	var line richtext.Line

	if powered, ok := a.adapter.Powered.Get(); ok {
		en, enStyle := " Enabled ", theme.Current().ADTree.Adapter.PoweredOnPill
		if !powered {
			en, enStyle = " Disabled ", theme.Current().ADTree.Adapter.PoweredOffPill
		}

		line = append(line, richtext.NewSegment(en, enStyle))
	}

	for _, state := range [3]struct {
		Name  string
		Opt   optional.Optional[bool]
		Style tcell.Style
	}{
		{" Scanning ", a.adapter.Discovering, theme.Current().ADTree.Adapter.ScanningPill},
		{" Discoverable ", a.adapter.Discoverable, theme.Current().ADTree.Adapter.DiscoverablePill},
		{" Pairable ", a.adapter.Pairable, theme.Current().ADTree.Adapter.PairablePill},
	} {
		if v, ok := state.Opt.Get(); !ok || !v {
			continue
		}

		line = append(line, richtext.NewSegment(state.Name, state.Style))
	}

	return ui.RenderMergedPill(line)
}

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (a *adapterNode) PopulateActions() {
	pos, _ := a.SubnodePosition(relPosActionsListNode)

	actionsListNode := a.Children()[pos]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data[*adNode]()
	if node == nil {
		return
	}

	if v, ok := a.adapter.Powered.Get(); ok {
		newActionNode(node, ui.Kb().Adapter.TogglePower, boolToActionState(v), true, a.actionPowered)
	}
	if v, ok := a.adapter.Discovering.Get(); ok {
		newActionNode(node, ui.Kb().Adapter.ToggleDiscoverable, boolToActionState(v), true, a.actionDiscoverable)
	}
	if v, ok := a.adapter.Pairable.Get(); ok {
		newActionNode(node, ui.Kb().Adapter.TogglePairable, boolToActionState(v), true, a.actionPairable)
	}
	if v, ok := a.adapter.Discovering.Get(); ok {
		newActionNode(node, ui.Kb().Adapter.ToggleScan, boolToActionState(v), true, a.actionScan)
	}

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
	case ui.Kb().Adapter.TogglePower:
		text = "Power"

	case ui.Kb().Adapter.ToggleDiscoverable:
		text = "Discoverable"

	case ui.Kb().Adapter.TogglePairable:
		text = "Pairable"

	case ui.Kb().Adapter.ToggleScan:
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

	actionNode.SetContent(buildLabel(actionNode.noder.Icon(), displayName, actionNode.noder.NodeStyle()))
}

// SetAdapterEventData sets the adapter event data for the node.
func (a *adapterNode) SetAdapterEventData(ev bluetooth.AdapterEventData) {
	a.adapter.AdapterEventData = ev

	a.Refresh()
	a.PopulateActions()
}

// SetDeviceEventData sets the device event data for the node.
func (a *adapterNode) SetDeviceEventData(bluetooth.DeviceEventData) {
}

func (a *adapterNode) actionPowered(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.NewCreationInfo(newAdapterStatus(a, "Turning adapter", state)), func(ov *operations.RunningInfo) tview.Msg {
		ov.Info("Powering on...")
		time.Sleep(2 * time.Second)

		return ov.Ok("Done")
	}
}

func (a *adapterNode) actionDiscoverable(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.NewCreationInfo(newAdapterStatus(a, "Setting adapter's discoverable", state)), func(ov *operations.RunningInfo) tview.Msg {
		ov.Info("Starting Disc")
		time.Sleep(2 * time.Second)

		return ov.Error("An error: ", errors.New("error"))
	}
}

func (a *adapterNode) actionPairable(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (a *adapterNode) actionScan(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

var _ adNoder = (*adapterNode)(nil)

type adapterStatus struct {
	name, icon, adapterName string
	currentState            actionStateSpec
}

func newAdapterStatus(node *adapterNode, name string, state *adActionState) *adapterStatus {
	return &adapterStatus{
		name:        name,
		adapterName: getAdapterDisplayName(node.adapter),
		icon:        node.Icon(),

		currentState: state.currentState,
	}
}

func (a *adapterStatus) Content() richtext.Text {
	builder := ui.NewBuilder(
		theme.Current().Global,
	)

	style := theme.Current().Operations.PropertyName
	adapterStyle := theme.Current().ADTree.Adapter.Name
	msgStyle := theme.Current().Operations.Message

	builder.AddKVFunc("Adapter:", style, func(bb *ui.Builder) {
		bb.Append(a.icon, adapterStyle)
		bb.Space()
		bb.Append(a.adapterName, adapterStyle)
	})

	builder.AddKVFunc("Action:", style, func(bb *ui.Builder) {
		bb.Append(a.name, msgStyle)
		bb.Space()
		bb.Append(a.currentState.Format("on", "off"), msgStyle)
	})

	return builder.Text()
}

var _ card.ContentProvider = (*adapterStatus)(nil)
