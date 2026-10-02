package adtree

import (
	"runtime"

	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
	"github.com/gdamore/tcell/v3"
)

// deviceNode describes a device node.
type deviceNode struct {
	*adNode
	emptyNoder

	device bluetooth.DeviceData
}

// newDeviceNode creates a new device node.
func newDeviceNode(rv ui.RootView, parentNode *adNode, device bluetooth.DeviceData) *adNode {
	dn := &deviceNode{}

	deviceNode := newAdNode(
		nodeTypeDevice, getDeviceDisplayName(device.DeviceEventData),
		true, newDeviceNodeID(device.DeviceAddress),
		dn, rv,
	)

	dn.device = device
	dn.adNode = deviceNode

	dn.AddActionsList()
	dn.PopulateActions()

	parentNode.AddChild(dn.Node)

	return deviceNode
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (d *deviceNode) HandleKeys(p tview.KeyMsg) (ui.RouterMsg, bool) {
	return actionHandleKeyMsg(p, d.Node)
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

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (d *deviceNode) PopulateActions() {
	actionsListNode := d.Children()[relPosActionsListNode]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data[*adNode]()
	if node == nil {
		return
	}

	node.AddAction(kb().Device.ToggleConnection, boolToActionState(d.device.Connected.Value()), true, d.actionConnect)
	node.AddAction(kb().Device.TogglePairedState, boolToActionState(d.device.Paired.Value()), true, d.actionPair)

	if d.rv.Features().Has(appfeatures.FeatureSendFile, appfeatures.FeatureReceiveFile) &&
		d.device.HaveService(bluetooth.ObexObjpushServiceClass) {
		node.AddAction(kb().Device.SendFiles, actionStateNone, false, d.actionSend)
	}

	if runtime.GOOS == "linux" {
		node.AddAction(kb().Device.Trust, boolToActionState(d.device.Trusted.Value()), true, d.actionTrust)
		node.AddAction(kb().Device.Block, boolToActionState(d.device.Blocked.Value()), true, d.actionBlock)

		if d.device.HaveService(bluetooth.AudioSourceServiceClass) ||
			d.device.HaveService(bluetooth.AudioSinkServiceClass) {
			node.AddAction(kb().Device.AudioProfiles, actionStateNone, false, d.actionAudioProfiles)
		}

		if d.device.HaveService(bluetooth.AudioSourceServiceClass) &&
			d.device.HaveService(bluetooth.AvRemoteServiceClass) &&
			d.device.HaveService(bluetooth.AvRemoteTargetServiceClass) {
			node.AddAction(kb().Player.ToggleDisplay, actionStateDisabled, true, d.actionMediaPlayer)
		}

		if d.rv.Features().Has(appfeatures.FeatureNetwork) &&
			d.device.HaveService(bluetooth.NapServiceClass) &&
			(d.device.HaveService(bluetooth.PanuServiceClass) ||
				d.device.HaveService(bluetooth.DialupNetServiceClass)) {
			node.AddAction(kb().Device.NetworkOptions, actionStateNone, false, d.actionNetwork)
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
	case kb().Device.ToggleConnection:
		enabledText = "Connect"
		disabledText = "Disconnect"

	case kb().Device.TogglePairedState:
		enabledText = "Pair"
		disabledText = "Unpair/Remove"

	case kb().Device.Trust:
		enabledText = "Trust"
		disabledText = "Untrust"

	case kb().Device.Block:
		enabledText = "Block"
		disabledText = "Unblock"

	case kb().Device.SendFiles:
		enabledText = "Send file(s)"

	case kb().Device.NetworkOptions:
		enabledText = "List network profiles (Bluetooth tethering)"

	case kb().Device.AudioProfiles:
		enabledText = "List audio profiles"

	//TODO: Show/hide
	case kb().Player.ToggleDisplay:
		enabledText = "Show media player"
		disabledText = "Hide media player"

	default:
		return
	}

	actionText := enabledText
	if state.isToggleable && state.currentState == actionStateEnabled && disabledText != "" {
		actionText = disabledText
	}

	actionNode.SetName(actionText)
}

// SetDeviceEventData sets the device event data for the node.
func (d *deviceNode) SetDeviceEventData(ev bluetooth.DeviceEventData) {
	d.device.DeviceEventData = ev
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
