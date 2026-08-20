package views

import (
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/Digital-Shane/treeview/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
)

type adTreeNode struct {
	id nodeID

	nodeType adTreeNodeType
	node     *treeview.Node[adTreeNode]

	noder       adNoder
	actionState *adActionState

	tree *adTree
}

type adTreeNodeType uint8

const (
	nodeTypeRoot adTreeNodeType = iota
	nodeTypeAdapter
	nodeTypeDevice
	nodeTypeDevicesList
	nodeTypeAction
	nodeTypeActionsList
)

type adSubNodePos uint8

const (
	actionsListNodePos adSubNodePos = iota
	devicesTreeNodePos              = 1
)

func newAdNode(ntype adTreeNodeType, name string, expanded bool, id nodeID, noder adNoder, tree *adTree) *adTreeNode {
	node := treeview.NewNode(id.String(), "", adTreeNode{})
	node.SetName(name)
	node.SetExpanded(expanded)

	data := node.Data()

	data.id = id
	data.node = node
	data.nodeType = ntype
	data.noder = noder
	data.tree = tree

	return data
}

func (a *adTreeNode) addDevicesList(devices []bluetooth.DeviceData) {
	dlAdNode := newAdNode(
		nodeTypeDevicesList, "Devices",
		true, a.id.appendSubNodeNib(nibDevicesList),
		a.noder, a.tree,
	)

	for _, device := range devices {
		newDeviceAdNode(a.tree, dlAdNode.node, device)
	}

	a.node.AddChild(dlAdNode.node)
}

func (a *adTreeNode) addActionsList() {
	alAdNode := newAdNode(
		nodeTypeActionsList, "Actions",
		false, a.id.appendSubNodeNib(nibActionsList),
		a.noder, a.tree,
	)

	alAdNode.node.SetChildren(make([]*treeview.Node[adTreeNode], 0, 5))

	a.node.AddChild(alAdNode.node)
}

func (a *adTreeNode) addAction(key keybindings.KeyID, state actionStateSpec, isToggleable bool, invoker actionInvoker) {
	actionNode := newAdNode(
		nodeTypeAction, "",
		false, a.id.appendSubNodeTextNib(nibAction, string(key)),
		a.noder, a.tree,
	)

	actionNode.actionState = newAdActionState(key, state, isToggleable, invoker)

	a.node.AddChild(actionNode.node)
}

type adNoder interface {
	handleKeys(p tea.KeyPressMsg) (tea.Msg, bool)
	populateActions()
	updateAction(action *treeview.Node[adTreeNode], updateMsg actionUpdateMsg)
}

type rootAdNode struct{}

func newRootAdNode(tree *adTree) *adTreeNode {
	rn := &rootAdNode{}
	id := newRootNodeID()

	rootAdNode := newAdNode(
		nodeTypeRoot, "Adapters",
		true, id.NodeID(),
		rn, tree, // TODO: Change
	)

	rootAdNode.node.SetChildren(make([]*treeview.Node[adTreeNode], 0, 10))

	return rootAdNode
}

func (r *rootAdNode) handleKeys(_ tea.KeyPressMsg) (tea.Msg, bool) {
	return nil, false
}

func (r *rootAdNode) populateActions() {
}

func (r *rootAdNode) updateAction(_ *treeview.Node[adTreeNode], _ actionUpdateMsg) {
}

type adapterAdNode struct {
	adapter bluetooth.AdapterData
	*adTreeNode
}

func newAdapterAdNode(tree *adTree, rootNode *treeview.Node[adTreeNode], adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) *adTreeNode {
	an := &adapterAdNode{}
	id := newAdapterNodeID(adapter.AdapterAddress)

	adapterNode := newAdNode(
		nodeTypeAdapter, getAdapterDisplayName(adapter),
		true, id.NodeID(),
		an, tree,
	)

	an.adTreeNode = adapterNode
	an.adapter = adapter

	an.addActionsList()
	an.addDevicesList(devices)

	an.populateActions()

	rootNode.AddChild(adapterNode.node)

	return adapterNode
}

func (a *adapterAdNode) handleKeys(p tea.KeyPressMsg) (tea.Msg, bool) {
	return actionHandleKeyMsg(p, a.node)
}

func (a *adapterAdNode) populateActions() {
	actionsListNode := a.node.Children()[actionsListNodePos]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data()
	if node == nil {
		return
	}

	node.addAction(keybindings.KeyAdapterTogglePower, boolToActionState(a.adapter.Powered.Value()), true, a.actionPowered)
	node.addAction(keybindings.KeyAdapterToggleDiscoverable, boolToActionState(a.adapter.Discovering.Value()), true, a.actionDiscoverable)
	node.addAction(keybindings.KeyAdapterTogglePairable, boolToActionState(a.adapter.Pairable.Value()), true, a.actionPairable)
	node.addAction(keybindings.KeyAdapterToggleScan, boolToActionState(a.adapter.Discovering.Value()), true, a.actionScan)

	for _, actionNode := range actionsListNode.Children() {
		a.updateAction(actionNode, emptyActionUpdateMsg())
	}
}

func (a *adapterAdNode) updateAction(actionNode *treeview.Node[adTreeNode], updateMsg actionUpdateMsg) {
	adnode := actionNode.Data()
	state := adnode.actionState

	updateMsg.updateState(actionNode.ID(), state)

	var text string

	switch state.key {
	case keybindings.KeyAdapterTogglePower:
		text = "Power"

	case keybindings.KeyAdapterToggleDiscoverable:
		text = "Discoverable"

	case keybindings.KeyAdapterTogglePairable:
		text = "Pairable"

	case keybindings.KeyAdapterToggleScan:
		text = "Device Scanning"

	default:
		return
	}

	actionText := "Off"
	if state.currentState == actionStateDisabled {
		actionText = "On"
	}

	displayName := useBuffer(len(text)+len(actionText)+10, func(b *strings.Builder) {
		b.WriteString("Switch ")

		b.WriteString(text)
		b.WriteString(" ")

		b.WriteString(actionText)
	})

	actionNode.SetName(displayName)
}

func (a *adapterAdNode) actionPowered() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (a *adapterAdNode) actionDiscoverable() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (a *adapterAdNode) actionPairable() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (a *adapterAdNode) actionScan() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

type deviceAdNode struct {
	device bluetooth.DeviceData
	*adTreeNode
}

func newDeviceAdNode(tree *adTree, parentNode *treeview.Node[adTreeNode], device bluetooth.DeviceData) *adTreeNode {
	dn := &deviceAdNode{}
	id := newDeviceNodeID(device.DeviceAddress)

	deviceNode := newAdNode(
		nodeTypeDevice, getDeviceDisplayName(device.DeviceEventData),
		true, id.NodeID(),
		dn, tree,
	)

	dn.device = device
	dn.adTreeNode = deviceNode

	dn.addActionsList()
	dn.populateActions()

	parentNode.AddChild(dn.node)

	return deviceNode
}

func (d *deviceAdNode) handleKeys(p tea.KeyPressMsg) (tea.Msg, bool) {
	return actionHandleKeyMsg(p, d.node)
}

func (d *deviceAdNode) populateActions() {
	device, err := d.tree.session().Device(d.device.DeviceAddress).Properties()
	if err != nil {
		return
	}

	actionsListNode := d.node.Children()[actionsListNodePos]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data()
	if node == nil {
		return
	}

	node.addAction(keybindings.KeyDeviceConnect, boolToActionState(d.device.Connected.Value()), true, d.actionConnect)
	node.addAction(keybindings.KeyDevicePair, boolToActionState(d.device.Paired.Value()), true, d.actionPair)

	if runtime.GOOS == "linux" {
		node.addAction(keybindings.KeyDeviceTrust, boolToActionState(d.device.Trusted.Value()), true, d.actionTrust)
		node.addAction(keybindings.KeyDeviceBlock, boolToActionState(d.device.Blocked.Value()), true, d.actionBlock)
	}

	if d.tree.features().Has(appfeatures.FeatureSendFile, appfeatures.FeatureReceiveFile) &&
		device.HaveService(bluetooth.ObexObjpushServiceClass) {
		node.addAction(keybindings.KeyDeviceSendFiles, actionStateNone, false, d.actionSend)
	}

	if d.tree.features().Has(appfeatures.FeatureNetwork) &&
		device.HaveService(bluetooth.NapServiceClass) &&
		(device.HaveService(bluetooth.PanuServiceClass) ||
			device.HaveService(bluetooth.DialupNetServiceClass)) {
		node.addAction(keybindings.KeyDeviceNetwork, actionStateNone, false, d.actionNetwork)
	}

	if device.HaveService(bluetooth.AudioSourceServiceClass) ||
		device.HaveService(bluetooth.AudioSinkServiceClass) {
		node.addAction(keybindings.KeyDeviceAudioProfiles, actionStateNone, false, d.actionAudioProfiles)
	}

	if device.HaveService(bluetooth.AudioSourceServiceClass) &&
		device.HaveService(bluetooth.AvRemoteServiceClass) &&
		device.HaveService(bluetooth.AvRemoteTargetServiceClass) {
		node.addAction(keybindings.KeyPlayerShow, actionStateDisabled, true, d.actionMediaPlayer)
	}

	for _, actionNode := range actionsListNode.Children() {
		d.updateAction(actionNode, emptyActionUpdateMsg())
	}
}

func (d *deviceAdNode) updateAction(actionNode *treeview.Node[adTreeNode], updateMsg actionUpdateMsg) {
	adnode := actionNode.Data()
	state := adnode.actionState

	updateMsg.updateState(actionNode.ID(), state)

	var enabledText, disabledText string

	switch state.key {
	case keybindings.KeyDeviceConnect:
		enabledText = "Connect"
		disabledText = "Disconnect"

	case keybindings.KeyDevicePair:
		enabledText = "Pair"
		disabledText = "Unpair/Remove"

	case keybindings.KeyDeviceTrust:
		enabledText = "Trust"
		disabledText = "Untrust"

	case keybindings.KeyDeviceBlock:
		enabledText = "Block"
		disabledText = "Unblock"

	case keybindings.KeyDeviceSendFiles:
		enabledText = "Send file(s)"

	case keybindings.KeyDeviceNetwork:
		enabledText = "List network profiles (Bluetooth tethering)"

	case keybindings.KeyDeviceAudioProfiles:
		enabledText = "List audio profiles"

	//TODO: Show/hide
	case keybindings.KeyPlayerShow:
		enabledText = "Show media player"
		disabledText = "Hide media player"

	default:
		return
	}

	actionText := enabledText
	if state.isToggleable && state.currentState == actionStateEnabled && disabledText != "" {
		actionText = disabledText
	}

	adnode.node.SetName(actionText)
}

func (d *deviceAdNode) actionConnect() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (d *deviceAdNode) actionPair() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (d *deviceAdNode) actionTrust() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (d *deviceAdNode) actionBlock() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (d *deviceAdNode) actionSend() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (d *deviceAdNode) actionNetwork() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (d *deviceAdNode) actionAudioProfiles() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

func (d *deviceAdNode) actionMediaPlayer() (opCreationInfo, opInvoker) {
	return opCreationInfo{}, nil
}

type nodeID struct {
	id string

	nodeNib nodeIDNib
}

type nodeIDNib = string

const (
	nibDevicesList nodeIDNib = "dl"
	nibActionsList nodeIDNib = "al"
	nibAction      nodeIDNib = "ac"
	nibAdapter     nodeIDNib = "ad"
	nibDevice      nodeIDNib = "dv"
)

func (n nodeID) appendSubNodeNib(nib nodeIDNib) nodeID {
	return n.appendSubNodeTextNib(nib, "")
}

func (n nodeID) appendSubNodeTextNib(nib nodeIDNib, text string) nodeID {
	length := len(n.id) + len(nib) + len(text)

	n.id = useBuffer(length, func(b *strings.Builder) {
		b.WriteString(n.id)
		b.WriteString("/")
		b.WriteString(nib)

		if text != "" {
			b.WriteString(":")
			b.WriteString(text)
		}
	})

	return n
}

func (n nodeID) String() string {
	return n.id
}

type rootNodeID struct {
	nodeID
}

func newRootNodeID() rootNodeID {
	return rootNodeID{nodeID: nodeID{id: "root", nodeNib: ""}}
}

func (r *rootNodeID) NodeID() nodeID {
	return r.nodeID
}

type adapterNodeID struct {
	nodeID

	adapterAddress bluetooth.AdapterAddress
}

func newAdapterNodeID(address bluetooth.AdapterAddress) adapterNodeID {
	a := adapterNodeID{adapterAddress: address}

	a.nodeNib = nibAdapter
	a.id = a.buildID()

	return a
}

func (a *adapterNodeID) NodeID() nodeID {
	return a.nodeID
}

func (a *adapterNodeID) buildID() string {
	length := 4 + len(a.nodeNib) + (bluetooth.MaxAddressStringLength + 2)

	return useBuffer(length, func(b *strings.Builder) {
		b.WriteString(a.nodeNib)
		b.WriteString(":")
		appendMacAddress(b, a.adapterAddress.Address)
	})
}

type deviceNodeID struct {
	nodeID

	deviceAddress bluetooth.DeviceAddress
}

func newDeviceNodeID(address bluetooth.DeviceAddress) deviceNodeID {
	d := deviceNodeID{deviceAddress: address}

	d.nodeNib = nibDevice
	d.id = d.buildID()

	return d
}

func (d *deviceNodeID) NodeID() nodeID {
	return d.nodeID
}

func (d *deviceNodeID) buildID() string {
	length := 4 + len(d.nodeNib) + ((bluetooth.MaxAddressStringLength + 2) * 2)

	return useBuffer(length, func(b *strings.Builder) {
		b.WriteString(d.nodeNib)
		b.WriteString(":")
		appendMacAddress(b, d.deviceAddress.Address)

		b.WriteString("/")

		b.WriteString(nibAdapter)
		b.WriteString(":")
		appendMacAddress(b, d.deviceAddress.AssociatedAdapter)
	})
}

func appendMacAddress(b *strings.Builder, mac bluetooth.MacAddress) {
	for i := 5; i >= 0; i-- {
		c := mac[i]
		nibble := c >> 4
		if nibble <= 9 {
			b.WriteByte(nibble + '0')
		} else {
			b.WriteByte(nibble + 'A' - 10)
		}

		nibble = c & 0x0f
		if nibble <= 9 {
			b.WriteByte(nibble + '0')
		} else {
			b.WriteByte(nibble + 'A' - 10)
		}
	}
}
