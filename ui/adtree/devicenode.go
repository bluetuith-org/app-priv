package adtree

import (
	"strconv"
	"strings"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetooth-classic/api/optional"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

// deviceNode describes a device node.
type deviceNode struct {
	*adNode

	alnode *alNode
	dlnode *dlNode

	device bluetooth.DeviceData
}

// newDeviceNode creates a new device node.
func newDeviceNode(rv ui.RootView, parentNode *dlNode, device bluetooth.DeviceData) *adNode {
	dn := &deviceNode{}

	deviceNode := newAdNode(
		newDeviceNodeID(device.DeviceAddress),
		dn, rv,
		true, 2,
	)

	dn.device = device
	dn.adNode = deviceNode

	dn.dlnode = parentNode
	dn.alnode = newActionsListNode(dn.adNode, dn)

	parentNode.AdNode().AddChild(dn.Node)

	dn.Refresh()
	dn.PopulateActions()

	return deviceNode
}

// AdNode returns the base node data.
func (d *deviceNode) AdNode() *adNode {
	return d.adNode
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (d *deviceNode) HandleKeys(p tview.KeyMsg) (ui.Router, bool) {
	alNode := d.alnode

	if msg, ok := actionHandleKeyMsg(p, alNode); ok {
		return msg, ok
	}

	return d.dlnode.HandleKeys(p)
}

// Refresh refreshes the content of the node.
func (d *deviceNode) Refresh() {
	d.SetMarkers(d.StatusIcons(true))
	d.SetContent(
		buildLabel(
			d.Icon(), getDeviceDisplayName(
				d.rv, d.device.DeviceEventData, false,
			),
			d.NodeStyle(),
		),
	)
}

// NodeStyle returns the style to be applied for this node's label.
func (d *deviceNode) NodeStyle() tcell.Style {
	paired, ok := d.device.Paired.Get()
	if !ok {
		return theme.Current().ADTree.Device.Unknown
	}
	if paired {
		return theme.Current().ADTree.Device.Paired
	}

	return theme.Current().ADTree.Device.Discovered
}

// Icon returns the Icon associated with this node.
func (d *deviceNode) Icon() string {
	return theme.DeviceClassToIcon(theme.Icons(), d.device.Class).String()
}

// RootID returns the root ID for the subnodes of an adapter or device node.
func (d *deviceNode) RootID() string {
	return d.id.String()
}

// Information returns a message to print adapter/device information.
func (d *deviceNode) Information() []text.Span {
	bgStyle := theme.Current().Global

	builder := stext.NewBuilder(bgStyle)

	infoPropStyle := theme.Current().Info.PropertyName
	infoPropValStyle := theme.Current().Info.PropertyValue

	builder.ResetKVStyle(infoPropStyle, infoPropValStyle)

	deviceNameStyle := theme.Current().ADTree.Device.Discovered
	if paired, ok := d.device.Paired.Get(); ok && paired {
		deviceNameStyle = theme.Current().ADTree.Device.Paired
	}

	builder.Appendln(" Device Information ", theme.Current().Info.Title)
	builder.Newline()

	builder.AddKVStyled(
		"Device:", infoPropStyle,
		getDeviceDisplayName(d.rv, d.device.DeviceEventData, false), deviceNameStyle,
	)
	builder.AddKV("Address:", printAddress(d.rv, d.device.Address, true))
	if alias, ok := d.device.Alias.Get(); ok && alias != "" {
		builder.AddKV("Alias:", alias)
	}

	var assocAdapter strings.Builder
	adapter := d.dlnode.AdapterNode().adapter

	assocAdapter.WriteString(getAdapterDisplayName(d.rv, adapter, true))

	builder.AddKVStyled("On Adapter:", infoPropStyle, assocAdapter.String(), theme.Current().ADTree.Adapter.Name)
	builder.Newline()

	builder.AddKV("Type:", printPlaceHolder(d.device.Type))
	builder.AddKV("Class:", strconv.FormatUint(uint64(d.device.Class), 10))
	builder.Newline()

	builder.AddKVFunc("Status:", func(bb *stext.Builder) {
		bb.Newline()
		bb.AppendSpans(
			ui.RenderMergedPill(d.StatusIcons(false), bgStyle.GetBackground())...,
		)
	})
	builder.Newline()

	printServiceUUIDs(builder, d.device.UUIDs)

	return builder.Text()
}

// StatusIcons returns the status icon pills for this node.
func (d *deviceNode) StatusIcons(compact bool) []text.Span {
	var line []text.Span

	for i, state := range [4]struct {
		Name  dvStatusIndicator
		Opt   optional.Optional[bool]
		Style tcell.Style
	}{
		{dvStatusUnpaired, d.device.Paired, theme.Current().ADTree.Device.NewDevicePill},
		{dvStatusConnected, d.device.Connected, theme.Current().ADTree.Device.ConnectedPill},
		{dvStatusTrusted, d.device.Trusted, theme.Current().ADTree.Device.TrustedPill},

		{dvStatusBlocked, d.device.Blocked, theme.Current().ADTree.Device.BlockedPill},
	} {
		v, ok := state.Opt.Get()
		if !ok {
			continue
		}

		switch {
		case i != 0 && !v:
			continue

		case i == 0 && v:
			state.Name = dvStatusPaired
			state.Style = theme.Current().ADTree.Device.PairedPill
		}

		line = append(line, stext.NewSpan(state.Name.String(compact), state.Style))
	}

	return line
}

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (d *deviceNode) PopulateActions() {
	node := d.alnode
	node.ClearActions()

	if v, ok := d.device.Connected.Get(); ok {
		newActionNode(node, ui.Kb().Device.ToggleConnection, boolToActionState(v), true, d.actionConnect)
	}
	if v, ok := d.device.Paired.Get(); ok {
		newActionNode(node, ui.Kb().Device.TogglePairedState, boolToActionState(v), true, d.actionPair)
	}

	if d.rv.Features().Has(appfeatures.FeatureSendFile, appfeatures.FeatureReceiveFile) &&
		d.device.HaveService(bluetooth.ObexObjpushServiceClass) {
		newActionNode(node, ui.Kb().Device.SendFiles, actionStateNone, false, d.actionSend)
	}

	if v, ok := d.device.Trusted.Get(); ok {
		newActionNode(node, ui.Kb().Device.Trust, boolToActionState(v), true, d.actionTrust)
	}
	if v, ok := d.device.Blocked.Get(); ok {
		newActionNode(node, ui.Kb().Device.Block, boolToActionState(v), true, d.actionBlock)
	}

	if d.rv.Features().Has(appfeatures.FeatureMediaPlayer) {
		if d.device.HaveService(bluetooth.AudioSourceServiceClass) ||
			d.device.HaveService(bluetooth.AudioSinkServiceClass) {
			newActionNode(node, ui.Kb().Device.AudioProfiles, actionStateNone, false, d.actionAudioProfiles)
		}

		if d.device.HaveService(bluetooth.AudioSourceServiceClass) &&
			d.device.HaveService(bluetooth.AvRemoteServiceClass) &&
			d.device.HaveService(bluetooth.AvRemoteTargetServiceClass) {
			newActionNode(node, ui.Kb().Player.ToggleDisplay, actionStateDisabled, true, d.actionMediaPlayer)
		}
	}

	if d.rv.Features().Has(appfeatures.FeatureNetwork) &&
		d.device.HaveService(bluetooth.NapServiceClass) &&
		(d.device.HaveService(bluetooth.PanuServiceClass) ||
			d.device.HaveService(bluetooth.DialupNetServiceClass)) {
		newActionNode(node, ui.Kb().Device.NetworkOptions, actionStateNone, false, d.actionNetwork)
	}

	for _, ac := range node.Actions() {
		d.UpdateActionNode(ac, emptyActionUpdateMsg())
	}
}

// UpdateActionNode updates the node of type [nodeTypeAction] with the message.
func (d *deviceNode) UpdateActionNode(actionNode *acNode, updateMsg actionUpdateMsg) {
	state := actionNode.actionState
	updateMsg.updateState(actionNode.AdNode().ID(), state)

	var enabledText, disabledText string

	switch state.key {
	case ui.Kb().Device.ToggleConnection:
		enabledText = "Connect"
		disabledText = "Disconnect"

	case ui.Kb().Device.TogglePairedState:
		enabledText = "Pair"
		disabledText = "Unpair/Remove"

	case ui.Kb().Device.Trust:
		enabledText = "Trust"
		disabledText = "Untrust"

	case ui.Kb().Device.Block:
		enabledText = "Block"
		disabledText = "Unblock"

	case ui.Kb().Device.SendFiles:
		enabledText = "Send file(s)"

	case ui.Kb().Device.NetworkOptions:
		enabledText = "List network profiles (Bluetooth tethering)"

	case ui.Kb().Device.AudioProfiles:
		enabledText = "List audio profiles"

	//TODO: Show/hide
	case ui.Kb().Player.ToggleDisplay:
		enabledText = "Show media player"
		disabledText = "Hide media player"

	default:
		return
	}

	actionText := enabledText
	if state.isToggleable && state.currentState == actionStateEnabled && disabledText != "" {
		actionText = disabledText
	}

	actionNode.AdNode().SetContent(buildLabel(actionNode.Icon(), actionText, actionNode.NodeStyle()))
}

// SetAdapterEventData sets the adapter event data for the node.
func (d *deviceNode) SetAdapterEventData(bluetooth.AdapterEventData) {
}

// SetDeviceEventData sets the device event data for the node.
func (d *deviceNode) SetDeviceEventData(ev bluetooth.DeviceEventData) {
	d.device.DeviceEventData = ev

	d.Refresh()
	d.PopulateActions()
}

func (d *deviceNode) actionConnect(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionPair(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionTrust(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionBlock(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionSend(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionNetwork(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionAudioProfiles(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionMediaPlayer(_ *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

var _ adNoder = (*deviceNode)(nil)

type dvStatusIndicator uint8

const (
	dvStatusNone dvStatusIndicator = iota
	dvStatusUnpaired
	dvStatusTrusted
	dvStatusPaired
	dvStatusConnected
	dvStatusBlocked
)

func (d dvStatusIndicator) String(compact bool) string {
	switch d {
	case dvStatusConnected:
		if compact {
			return " C "
		}

		return " Connected "

	case dvStatusBlocked:
		if compact {
			return " Bl "
		}

		return " Blocked "

	case dvStatusPaired:
		if compact {
			return " P "
		}

		return " Paired "

	case dvStatusTrusted:
		if compact {
			return " Tr "
		}

		return " Trusted "

	case dvStatusUnpaired:
		if compact {
			return " New "
		}

		return " Unpaired "
	}

	return ""
}
