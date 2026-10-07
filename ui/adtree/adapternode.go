package adtree

import (
	"errors"
	"iter"
	"strings"
	"time"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetooth-classic/api/optional"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

// adapterNode describes an adapter node.
type adapterNode struct {
	*adNode

	adapter bluetooth.AdapterData
}

// newAdapterNode creates a new adapter node.
func newAdapterNode(rv ui.RootView, rootNode *rootNode, adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) *adNode {
	an := &adapterNode{}

	adapterNode := newAdNode(
		nodeTypeAdapter,
		true, newAdapterNodeID(adapter.AdapterAddress),
		an, rv,
	)

	an.adNode = adapterNode
	an.adapter = adapter

	newActionsListNode(an.adNode)
	newDeviceListNode(adapterNode, devices)

	rootNode.AddChild(adapterNode.Node)
	an.Refresh()

	return adapterNode
}

// Refresh refreshes the content of the node.
func (a *adapterNode) Refresh() {
	a.SetContent(buildLabel(a.Icon(), getAdapterDisplayName(a.adapter), a.NodeStyle()))
	a.PopulateActions()
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

// RootID returns the root ID for the subnodes of an adapter or device node.
func (a *adapterNode) RootID() string {
	return a.id.String()
}

// Information returns a message to print adapter/device information.
func (a *adapterNode) Information() (string, richtext.Text) {
	return a.id.String(), nil
}

// StatusIcons returns the status icon pills for this node.
func (a *adapterNode) StatusIcons(bgStyle tcell.Style) iter.Seq[richtext.Segment] {
	return func(yield func(richtext.Segment) bool) {
		style := tcell.StyleDefault.Foreground(color.Black).Background(color.Green)

		if powered, ok := a.adapter.Powered.Get(); ok {
			en, enStyle := "Enabled", style
			if !powered {
				en, enStyle = "Disabled", style
			}

			if !pushSeg(ui.RenderPillIcon(en, enStyle), bgStyle, yield) {
				return
			}
		}

		for _, state := range [3]struct {
			Name string
			Opt  optional.Optional[bool]
		}{
			{"Scanning", a.adapter.Discovering},
			{"Discoverable", a.adapter.Discoverable},
			{"Pairable", a.adapter.Pairable},
		} {
			if v, ok := state.Opt.Get(); !ok || !v {
				continue
			}

			if !pushSeg(
				ui.RenderPillIcon(state.Name, style),
				bgStyle, yield,
			) {
				return
			}
		}
	}
}

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (a *adapterNode) PopulateActions() {
	actionsListNode := a.Children()[relPosActionsListNode]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data[*adNode]()
	if node == nil {
		return
	}

	newActionNode(node, ui.Kb().Adapter.TogglePower, boolToActionState(a.adapter.Powered.Value()), true, a.actionPowered)
	newActionNode(node, ui.Kb().Adapter.ToggleDiscoverable, boolToActionState(a.adapter.Discovering.Value()), true, a.actionDiscoverable)
	newActionNode(node, ui.Kb().Adapter.TogglePairable, boolToActionState(a.adapter.Pairable.Value()), true, a.actionPairable)
	newActionNode(node, ui.Kb().Adapter.ToggleScan, boolToActionState(a.adapter.Discovering.Value()), true, a.actionScan)

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
}

// SetDeviceEventData sets the device event data for the node.
func (a *adapterNode) SetDeviceEventData(bluetooth.DeviceEventData) {
}

func (a *adapterNode) actionPowered(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.NewCreationInfo(a.getDesc("Turning adapter", state)), func(ov *operations.RunningInfo) tview.Msg {
		ov.Info("Powering on...")
		time.Sleep(2 * time.Second)

		return ov.Ok("Done")
	}
}

func (a *adapterNode) actionDiscoverable(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.NewCreationInfo(a.getDesc("Setting adapter's discoverable", state)), func(ov *operations.RunningInfo) tview.Msg {
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

func (a *adapterNode) getDesc(name string, state *adActionState) richtext.Text {
	builder := ui.NewBuilder(
		theme.Current().Global,
	)

	style := theme.Current().Operations.PropertyName
	adapterStyle := theme.Current().ADTree.Adapter
	msgStyle := theme.Current().Operations.Message

	builder.AddKVFunc("Adapter:", style, func(bb *ui.Builder) {
		bb.Append(a.Icon(), adapterStyle)
		bb.Space()
		bb.Append(getAdapterDisplayName(a.adapter), adapterStyle)
	})

	builder.AddKVFunc("Action:", style, func(bb *ui.Builder) {
		bb.Append(name, msgStyle)
		bb.Space()
		bb.Append(state.currentState.Format("on", "off"), msgStyle)
	})

	return builder.Text()
}

var _ adNoder = (*adapterNode)(nil)
