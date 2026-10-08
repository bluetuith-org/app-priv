package adtree

import (
	"runtime"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetooth-classic/api/optional"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
	"github.com/gdamore/tcell/v3"
)

// deviceNode describes a device node.
type deviceNode struct {
	*adNode

	stn         *stNode
	parentNoder adNoder
	device      bluetooth.DeviceData
}

// newDeviceNode creates a new device node.
func newDeviceNode(rv ui.RootView, parentNode *adNode, device bluetooth.DeviceData) *adNode {
	dn := &deviceNode{}

	deviceNode := newAdNode(
		nodeTypeDevice,
		true, 10,
		newDeviceNodeID(device.DeviceAddress),
		dn, rv,
	)

	dn.device = device
	dn.adNode = deviceNode
	dn.parentNoder = parentNode.noder
	dn.stn = newStatusNode(deviceNode, false)

	newActionsListNode(dn.adNode)
	parentNode.AddChild(dn.Node)

	dn.Refresh()
	dn.PopulateActions()

	return deviceNode
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (d *deviceNode) HandleKeys(p tview.KeyMsg) (ui.RouterMsg, bool) {
	if msg, ok := actionHandleKeyMsg(p, d.Node); ok {
		return msg, ok
	}

	return d.parentNoder.HandleKeys(p)
}

// Refresh refreshes the content of the node.
func (d *deviceNode) Refresh() {
	d.stn.SetContent(d.StatusIcons())

	d.SetContent(buildLabel(d.Icon(), getDeviceDisplayName(d.device.DeviceEventData), d.NodeStyle()))
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

// SubnodePosition returns the position of the subnode within the parent node.
func (d *deviceNode) SubnodePosition(relpos adSubNodePos) (int, bool) {
	switch relpos {
	case relPosStatusNode:
		return 0, true

	case relPosActionsListNode:
		return 1, true
	}

	return -1, false
}

// Information returns a message to print adapter/device information.
func (d *deviceNode) Information() (string, richtext.Text) {
	return d.id.String(), nil
}

// StatusIcons returns the status icon pills for this node.
func (d *deviceNode) StatusIcons() richtext.Line {
	var line richtext.Line

	for i, state := range [4]struct {
		Name  string
		Opt   optional.Optional[bool]
		Style tcell.Style
	}{
		{" New ", d.device.Paired, theme.Current().ADTree.Device.NewDevicePill},
		{" Connected ", d.device.Connected, theme.Current().ADTree.Device.ConnectedPill},
		{" Trusted ", d.device.Trusted, theme.Current().ADTree.Device.TrustedPill},

		{" Blocked ", d.device.Blocked, theme.Current().ADTree.Device.BlockedPill},
	} {
		v, ok := state.Opt.Get()
		if !ok {
			continue
		}

		switch {
		case i != 0 && !v:
			continue

		case i == 0 && v:
			state.Name = " Paired "
			state.Style = theme.Current().ADTree.Device.PairedPill
		}

		line = append(line, richtext.NewSegment(state.Name, state.Style))
	}

	return ui.RenderMergedPill(line)
}

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (d *deviceNode) PopulateActions() {
	pos, _ := d.SubnodePosition(relPosActionsListNode)

	actionsListNode := d.Children()[pos]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data[*adNode]()
	if node == nil {
		return
	}

	newActionNode(node, ui.Kb().Device.ToggleConnection, boolToActionState(d.device.Connected.Value()), true, d.actionConnect)
	newActionNode(node, ui.Kb().Device.TogglePairedState, boolToActionState(d.device.Paired.Value()), true, d.actionPair)

	if d.rv.Features().Has(appfeatures.FeatureSendFile, appfeatures.FeatureReceiveFile) &&
		d.device.HaveService(bluetooth.ObexObjpushServiceClass) {
		newActionNode(node, ui.Kb().Device.SendFiles, actionStateNone, false, d.actionSend)
	}

	if runtime.GOOS == "linux" {
		newActionNode(node, ui.Kb().Device.Trust, boolToActionState(d.device.Trusted.Value()), true, d.actionTrust)
		newActionNode(node, ui.Kb().Device.Block, boolToActionState(d.device.Blocked.Value()), true, d.actionBlock)

		if d.device.HaveService(bluetooth.AudioSourceServiceClass) ||
			d.device.HaveService(bluetooth.AudioSinkServiceClass) {
			newActionNode(node, ui.Kb().Device.AudioProfiles, actionStateNone, false, d.actionAudioProfiles)
		}

		if d.device.HaveService(bluetooth.AudioSourceServiceClass) &&
			d.device.HaveService(bluetooth.AvRemoteServiceClass) &&
			d.device.HaveService(bluetooth.AvRemoteTargetServiceClass) {
			newActionNode(node, ui.Kb().Player.ToggleDisplay, actionStateDisabled, true, d.actionMediaPlayer)
		}

		if d.rv.Features().Has(appfeatures.FeatureNetwork) &&
			d.device.HaveService(bluetooth.NapServiceClass) &&
			(d.device.HaveService(bluetooth.PanuServiceClass) ||
				d.device.HaveService(bluetooth.DialupNetServiceClass)) {
			newActionNode(node, ui.Kb().Device.NetworkOptions, actionStateNone, false, d.actionNetwork)
		}
	}

	for _, actionNode := range actionsListNode.Children() {
		d.UpdateActionNode(actionNode.Data[*adNode](), emptyActionUpdateMsg())
	}
}

// UpdateActionNode updates the node of type [nodeTypeAction] with the message.
func (d *deviceNode) UpdateActionNode(actionNode *adNode, updateMsg actionUpdateMsg) {
	state := actionNode.actionState
	updateMsg.updateState(actionNode.ID(), state)

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

	actionNode.SetContent(buildLabel(actionNode.noder.Icon(), actionText, actionNode.noder.NodeStyle()))
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

func (d *deviceNode) actionConnect(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionPair(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionTrust(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionBlock(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionSend(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionNetwork(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionAudioProfiles(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

func (d *deviceNode) actionMediaPlayer(state *adActionState) (operations.CreationInfo, operations.Invoker) {
	return operations.CreationInfo{}, nil
}

var _ adNoder = (*deviceNode)(nil)
