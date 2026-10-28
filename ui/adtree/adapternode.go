package adtree

import (
	"errors"
	"strings"
	"time"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetooth-classic/api/optional"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

// adapterNode describes an adapter node.
type adapterNode struct {
	*adNode

	alnode *alNode
	dlnode *dlNode

	adapter bluetooth.AdapterData
}

// newAdapterNode creates a new adapter node.
func newAdapterNode(rv ui.RootView, rootNode *rootNode, adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) *adNode {
	an := &adapterNode{}

	adapterNode := newAdNode(
		nodeTypeAdapter,
		true, 2,
		newAdapterNodeID(adapter.AdapterAddress),
		an, rv,
	)

	an.adNode = adapterNode
	an.adapter = adapter

	an.alnode = newActionsListNode(an.adNode, an)
	an.dlnode = newDeviceListNode(an, devices)

	rootNode.AddChild(adapterNode.Node)

	an.Refresh()
	an.PopulateActions()

	return adapterNode
}

// AdNode returns the base node data.
func (a *adapterNode) AdNode() *adNode {
	return a.adNode
}

// Refresh refreshes the content of the node.
func (a *adapterNode) Refresh() {
	a.SetMarkers(a.StatusIcons(true))
	a.SetContent(
		buildLabel(
			a.Icon(), getAdapterDisplayName(a.rv, a.adapter, false),
			a.NodeStyle(),
		),
	)
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (a *adapterNode) HandleKeys(p tview.KeyMsg) (ui.Router, bool) {
	return actionHandleKeyMsg(p, a.alnode)
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

// Information returns a message to print adapter/device information.
func (a *adapterNode) Information() []text.Span {
	bgStyle := theme.Current().Global
	builder := stext.NewBuilder(bgStyle)

	infoPropStyle := theme.Current().Info.PropertyName
	infoPropValStyle := theme.Current().Info.PropertyValue
	adapterStyle := theme.Current().ADTree.Adapter.Name

	builder.ResetKVStyle(infoPropStyle, infoPropValStyle)

	builder.Appendln(" Adapter Information ", theme.Current().Info.Title)
	builder.Newline()

	builder.AddKVFunc("Name:", func(bb *stext.Builder) {
		bb.Append(getAdapterDisplayName(a.rv, a.adapter, false), adapterStyle)
	})
	builder.AddKVFunc("Address:", func(bb *stext.Builder) {
		bb.Append(printAddress(a.rv, a.adapter.Address, true), adapterStyle)
	})
	if a.adapter.UniqueName != "" {
		builder.AddKV("Unique Name:", a.adapter.UniqueName)
	}
	builder.Newline()

	builder.AddKVFunc("Status:", func(bb *stext.Builder) {
		bb.Newline()
		pills := ui.RenderMergedPill(a.StatusIcons(false), bgStyle.GetBackground())
		bb.AppendSpans(pills...)
	})
	builder.Newline()

	builder.AddKVFunc("Devices:", func(bb *stext.Builder) {
		bb.Newline()

		for _, node := range a.dlnode.AdNode().Children() {
			dvNode := nvget[*deviceNode](node)

			bb.Append(" ", bgStyle)
			bb.Append(theme.Icons().BottomLeft.String(), infoPropValStyle)
			bb.Append(" ", bgStyle)

			bb.AppendSpans(dvNode.Content()...)

			bb.Newline()
			bb.Newline()
		}
	})

	printServiceUUIDs(builder, a.adapter.UUIDs)

	return builder.Text()
}

// StatusIcons returns the status icon pills for this node.
func (a *adapterNode) StatusIcons(compact bool) []text.Span {
	var line []text.Span

	if powered, ok := a.adapter.Powered.Get(); ok {
		en, enStyle := adStatusPoweredOn, theme.Current().ADTree.Adapter.PoweredOnPill
		if !powered {
			en, enStyle = adStatusPoweredOff, theme.Current().ADTree.Adapter.PoweredOffPill
		}

		line = append(line, stext.NewSpan(en.String(compact), enStyle))
	}

	for _, state := range [3]struct {
		Name  adStatusIndicator
		Opt   optional.Optional[bool]
		Style tcell.Style
	}{
		{adStatusScanning, a.adapter.Discovering, theme.Current().ADTree.Adapter.ScanningPill},
		{adStatusDiscoverable, a.adapter.Discoverable, theme.Current().ADTree.Adapter.DiscoverablePill},
		{adStatusPairable, a.adapter.Pairable, theme.Current().ADTree.Adapter.PairablePill},
	} {
		if v, ok := state.Opt.Get(); !ok || !v {
			continue
		}

		line = append(line, stext.NewSpan(state.Name.String(compact), state.Style))
	}

	return line
}

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (a *adapterNode) PopulateActions() {
	node := a.alnode
	node.ClearActions()

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

	for _, ac := range node.Actions() {
		a.UpdateActionNode(ac, emptyActionUpdateMsg())
	}
}

// UpdateActionNode updates the node of type [nodeTypeAction] with the message.
func (a *adapterNode) UpdateActionNode(actionNode *acNode, updateMsg actionUpdateMsg) {
	state := actionNode.actionState
	updateMsg.updateState(actionNode.AdNode().ID(), state)

	var txt string

	switch state.key {
	case ui.Kb().Adapter.TogglePower:
		txt = "Power"

	case ui.Kb().Adapter.ToggleDiscoverable:
		txt = "Discoverable"

	case ui.Kb().Adapter.TogglePairable:
		txt = "Pairable"

	case ui.Kb().Adapter.ToggleScan:
		txt = "Device Scanning"

	default:
		return
	}

	actionText := "Off"
	if state.currentState == actionStateDisabled {
		actionText = "On"
	}

	displayName := useStringBuffer(len(txt)+len(actionText)+10, func(b *strings.Builder) {
		b.WriteString("Switch ")

		b.WriteString(txt)
		b.WriteString(" ")

		b.WriteString(actionText)
	})

	actionNode.AdNode().SetContent(buildLabel(actionNode.Icon(), displayName, actionNode.NodeStyle()))
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

func (a *adapterNode) actionPairable(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (a *adapterNode) actionScan(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
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
		adapterName: getAdapterDisplayName(node.rv, node.adapter, true),
		icon:        node.Icon(),

		currentState: state.currentState,
	}
}

func (a *adapterStatus) OperationDescription() []text.Span {
	builder := stext.NewBuilder(
		theme.Current().Global,
	)

	style := theme.Current().Operations.PropertyName
	adapterStyle := theme.Current().ADTree.Adapter.Name
	msgStyle := theme.Current().Operations.Message

	builder.AddKVStyledFunc("Adapter:", style, func(bb *stext.Builder) {
		bb.Append(a.icon, adapterStyle)
		bb.Space()
		bb.Append(a.adapterName, adapterStyle)
	})

	builder.AddKVStyledFunc("Action:", style, func(bb *stext.Builder) {
		bb.Append(a.name, msgStyle)
		bb.Space()
		bb.Append(a.currentState.Format("on", "off"), msgStyle)
	})

	return builder.Text()
}

var _ operations.Content = (*adapterStatus)(nil)

type adStatusIndicator uint8

const (
	adStatusNone adStatusIndicator = iota
	adStatusPoweredOff
	adStatusPoweredOn
	adStatusScanning
	adStatusDiscoverable
	adStatusPairable
)

func (a adStatusIndicator) String(compact bool) string {
	switch a {
	case adStatusPoweredOff:
		if compact {
			return " Off "
		}

		return " Not Powered "

	case adStatusPoweredOn:
		if compact {
			return " On "
		}

		return " Powered "

	case adStatusScanning:
		if compact {
			return " Sc "
		}

		return " Scanning "

	case adStatusDiscoverable:
		if compact {
			return " Ds "
		}

		return " Discoverable "

	case adStatusPairable:
		if compact {
			return " Pbl "
		}

		return " Pairable "
	}

	return ""
}
