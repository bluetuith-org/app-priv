package views

import (
	"iter"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/ayn2op/tview/tree"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

// adTreeModel represents a tree of adapters and its associated devices.
type adTreeModel struct {
	*tree.Model

	rootNode *rootAdNode

	focused bool
	rv      rootView
}

// ViewID returns the view's ID.
func (a *adTreeModel) ViewID() viewID {
	return viewIDAdTree
}

// Initialize initializes a view.
func (a *adTreeModel) Initialize() error {
	a.rootNode = newRootAdNode(a)

	prefixes := make([]string, 0, 12)
	prefixes = append(prefixes, "")
	for range 10 {
		prefixes = append(prefixes, " ")
	}

	treeModel := tree.NewModel()
	treeModel.SetRoot(a.rootNode.Node)
	treeModel.SetCurrentNode(treeModel.Root())
	treeModel.SetPrefixes(prefixes)
	treeModel.SetBackgroundColor(color.Gray)

	a.Model = treeModel
	a.UpdateStyles()

	return nil
}

// SetRootView sets the root view upon which the view is rendered.
// This will enable the view to access app-specific functions and send
// routed messages.
func (a *adTreeModel) SetRootView(v rootView) {
	a.rv = v
}

// AttachToTabView attaches this view to the tabbed view.
func (a *adTreeModel) AttachToTabView() (tabSection, bool) {
	return nil, false
}

// HandleRouterMsg handles the routed message.
func (a *adTreeModel) HandleRouterMsg(m routerMsg) tview.Cmd {
	return handleRouterMsg(a, m)
}

// SetFocus sets whether the view is currently focused.
func (a *adTreeModel) SetFocus(focused bool) {
	a.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (a *adTreeModel) GetFocus() bool {
	return a.focused
}

// UpdateStyles updates the styles for the view.
func (a *adTreeModel) UpdateStyles() {
	a.Model.SetBackgroundColor(color.Gray)
}

// RefreshContent refreshes the content of the view.
func (a *adTreeModel) RefreshContent() {
}

// Update receives messages when this model has focus.
func (a *adTreeModel) Update(msg tview.Msg) tview.Cmd {
	return a.Model.Update(msg)
}

// View draws this model onto the screen.
func (a *adTreeModel) View(screen tview.Screen) {
	a.Model.View(screen)
}

func (a *adTreeModel) features() *appfeatures.FeatureSet {
	return a.rv.Features()
}

type treeUpdateMsg struct{}

func msgAdTreeUpdate() routerMsg {
	return viewIDAdTree.routerMessage(treeUpdateMsg{})
}

func msgAdActionUpdate(id string, stateSpec actionStateSpec) routerMsg {
	return viewIDAdTree.routerMessage(actionUpdateMsg{id, stateSpec})
}

var _ view = (*adTreeModel)(nil)

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
	relPosActionsListNode adSubNodePos = iota
	relPosDevicesListNode              = 1
)

type nodeIDNib = string

const (
	nibDevicesList nodeIDNib = "dl"
	nibActionsList nodeIDNib = "al"
	nibAction      nodeIDNib = "ac"
	nibAdapter     nodeIDNib = "ad"
	nibDevice      nodeIDNib = "dv"
)

const (
	lenNibPlusColon = 3
	lenBdAddr       = 12
	lenNodeIDTotal  = lenNibPlusColon + lenBdAddr
)

type nodeID struct {
	id string

	nodeNib nodeIDNib
}

func newRootNodeID() nodeID {
	return nodeID{id: "root", nodeNib: ""}
}

func newAdapterNodeID(address bluetooth.AdapterAddress) nodeID {
	n := nodeID{
		nodeNib: nibAdapter,
		id: useStringBuffer(lenNodeIDTotal, func(b *strings.Builder) {
			b.WriteString(nibAdapter)
			b.WriteString(":")
			appendMacAddress(b, address.Address)
		}),
	}

	return n
}

func newDeviceNodeID(address bluetooth.DeviceAddress) nodeID {
	return nodeID{
		nodeNib: nibDevice,
		id: useStringBuffer((lenNodeIDTotal*2)+2+len(nibDevicesList), func(b *strings.Builder) {
			b.WriteString(nibAdapter)
			b.WriteString(":")
			appendMacAddress(b, address.AssociatedAdapter)

			b.WriteString("/")
			b.WriteString(nibDevicesList)
			b.WriteString("/")

			b.WriteString(nibDevice)
			b.WriteString(":")
			appendMacAddress(b, address.Address)
		}),
	}
}

func (n nodeID) appendSubNodeNib(nib nodeIDNib) nodeID {
	return n.appendSubNodeTextNib(nib, "")
}

func (n nodeID) appendSubNodeTextNib(nib nodeIDNib, txt string) nodeID {
	length := len(n.id) + lenNibPlusColon + len(txt)

	n.id = useStringBuffer(length, func(b *strings.Builder) {
		b.WriteString(n.id)
		b.WriteString("/")
		b.WriteString(nib)

		if txt != "" {
			b.WriteString(":")
			b.WriteString(txt)
		}
	})

	return n
}

func (n nodeID) String() string {
	return n.id
}

type adNoder interface {
	handleKeys(p tview.KeyMsg) (routerMsg, bool)
	populateActions()
	updateActionNode(actionNode *adTreeNode, updateMsg actionUpdateMsg)
	setAdapterEventData(ev bluetooth.AdapterEventData)
	setDeviceEventData(ev bluetooth.DeviceEventData)
}

type adTreeNode struct {
	*tree.Node

	parent *adTreeNode

	id       nodeID
	nodeType adTreeNodeType

	noder       adNoder
	actionState *adActionState

	tree *adTreeModel
}

func newAdNode(ntype adTreeNodeType, name string, expanded bool, id nodeID, noder adNoder, adtree *adTreeModel) *adTreeNode {
	data := &adTreeNode{
		Node:     tree.NewNode(name).SetIndent(2).SetExpanded(expanded),
		id:       id,
		noder:    noder,
		nodeType: ntype,
		tree:     adtree,
	}

	return data
}

func (a *adTreeNode) addDevicesList(devices []bluetooth.DeviceData) {
	dlAdNode := newAdNode(
		nodeTypeDevicesList, "Devices",
		true, a.id.appendSubNodeNib(nibDevicesList),
		a.noder, a.tree,
	)

	for _, device := range devices {
		newDeviceAdNode(a.tree, dlAdNode, device)
	}

	dlAdNode.parent = a

	a.Node.AddChild(dlAdNode.Node)
}

func (a *adTreeNode) addActionsList() {
	alAdNode := newAdNode(
		nodeTypeActionsList, "Actions",
		false, a.id.appendSubNodeNib(nibActionsList),
		a.noder, a.tree,
	)

	alAdNode.Node.SetChildren(make([]*tree.Node, 0, 5))
	alAdNode.parent = a

	a.Node.AddChild(alAdNode.Node)
}

func (a *adTreeNode) addAction(key keybindings.Keybinding, state actionStateSpec, isToggleable bool, invoker actionInvoker) {
	actionNode := newAdNode(
		nodeTypeAction, "",
		false, a.id.appendSubNodeTextNib(nibAction, strconv.Itoa(int(key.ID))),
		a.noder, a.tree,
	)

	actionNode.actionState = newAdActionState(key, state, isToggleable, invoker)
	actionNode.parent = a

	a.Node.AddChild(actionNode.Node)
}

type rootAdNode struct {
	*adTreeNode
	emptyNoder

	sync.RWMutex
}

func newRootAdNode(adtree *adTreeModel) *rootAdNode {
	rn := &rootAdNode{}

	rootAdNode := newAdNode(
		nodeTypeRoot, "Adapters",
		true, newRootNodeID(),
		rn, adtree,
	)

	rootAdNode.SetChildren(make([]*tree.Node, 0, 10))
	rn.adTreeNode = rootAdNode

	return rn
}

func (r *rootAdNode) addAdapter(adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) {
	r.Lock()
	defer r.Unlock()

	newAdapterAdNode(r.tree, r, adapter, devices)
}

func (r *rootAdNode) addDevice(device bluetooth.DeviceData) {
	r.Lock()
	defer r.Unlock()

	dlNodeID := newAdapterNodeID(device.AdapterAddress()).appendSubNodeNib(nibDevicesList)

	deviceListNode, _, ok := findTreeNodeByID(r, dlNodeID.String(), nodeTypeDevicesList)
	if !ok {
		return
	}

	newDeviceAdNode(r.tree, deviceListNode, device)
}

func (r *rootAdNode) updateAdapter(adapterEvent bluetooth.AdapterEventData, remove bool) {
	r.Lock()
	defer r.Unlock()

	adNodeID := newAdapterNodeID(adapterEvent.AdapterAddress)

	adapterNode, pos, ok := findTreeNodeByID(r, adNodeID.String(), nodeTypeAdapter)
	if !ok {
		return
	}

	if remove {
		parentNode := r
		parentNode.SetChildren(slices.Delete(parentNode.Children(), pos, pos+1))

		return
	}

	adapterNode.noder.setAdapterEventData(adapterEvent)
}

func (r *rootAdNode) updateDevice(deviceEvent bluetooth.DeviceEventData, remove bool) {
	r.Lock()
	defer r.Unlock()

	dvNodeID := newDeviceNodeID(deviceEvent.DeviceAddress)

	deviceNode, pos, ok := findTreeNodeByID(r, dvNodeID.String(), nodeTypeDevice)
	if !ok {
		return
	}

	if remove {
		parentNode := deviceNode.parent
		if parentNode != nil {
			parentNode.SetChildren(slices.Delete(parentNode.Children(), pos, pos+1))
		}

		return
	}

	deviceNode.noder.setDeviceEventData(deviceEvent)
}

func (r *rootAdNode) updateAction(updateMsg actionUpdateMsg) {
	r.Lock()
	defer r.Unlock()

	node, _, ok := findTreeNodeByID(r, updateMsg.id, nodeTypeAction)
	if !ok {
		return
	}

	node.noder.updateActionNode(node, updateMsg)
}

type adapterAdNode struct {
	*adTreeNode
	emptyNoder

	adapter bluetooth.AdapterData
}

func newAdapterAdNode(adtree *adTreeModel, rootNode *rootAdNode, adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) *adTreeNode {
	an := &adapterAdNode{}

	adapterNode := newAdNode(
		nodeTypeAdapter, getAdapterDisplayName(adapter),
		true, newAdapterNodeID(adapter.AdapterAddress),
		an, adtree,
	)

	an.adTreeNode = adapterNode
	an.adapter = adapter

	an.addActionsList()
	an.addDevicesList(devices)

	an.populateActions()

	rootNode.AddChild(adapterNode.Node)

	return adapterNode
}

func (a *adapterAdNode) handleKeys(p tview.KeyMsg) (routerMsg, bool) {
	return actionHandleKeyMsg(p, a.adTreeNode)
}

func (a *adapterAdNode) populateActions() {
	actionsListNode := a.Node.Children()[relPosActionsListNode]
	actionsListNode.SetChildren(nil)

	node := getNodeReference(actionsListNode)
	if node == nil {
		return
	}

	node.addAction(kb().Adapter.TogglePower, boolToActionState(a.adapter.Powered.Value()), true, a.actionPowered)
	node.addAction(kb().Adapter.ToggleDiscoverable, boolToActionState(a.adapter.Discovering.Value()), true, a.actionDiscoverable)
	node.addAction(kb().Adapter.TogglePairable, boolToActionState(a.adapter.Pairable.Value()), true, a.actionPairable)
	node.addAction(kb().Adapter.ToggleScan, boolToActionState(a.adapter.Discovering.Value()), true, a.actionScan)

	for _, actionNode := range actionsListNode.Children() {
		a.updateActionNode(getNodeReference(actionNode), emptyActionUpdateMsg())
	}
}

func (a *adapterAdNode) updateActionNode(actionNode *adTreeNode, updateMsg actionUpdateMsg) {
	adnode := actionNode
	state := adnode.actionState

	updateMsg.updateState(actionNode.id.String(), state)

	var txt string

	switch state.key {
	case kb().Adapter.TogglePower:
		txt = "Power"

	case kb().Adapter.ToggleDiscoverable:
		txt = "Discoverable"

	case kb().Adapter.TogglePairable:
		txt = "Pairable"

	case kb().Adapter.ToggleScan:
		txt = "Device Scanning"

	default:
		return
	}

	actiontxt := "Off"
	if state.currentState == actionStateDisabled {
		actiontxt = "On"
	}

	displayName := useStringBuffer(len(txt)+len(actiontxt)+10, func(b *strings.Builder) {
		b.WriteString("Switch ")

		b.WriteString(txt)
		b.WriteString(" ")

		b.WriteString(actiontxt)
	})

	actionNode.SetLine(text.NewLine(text.NewSegment(displayName, tcell.StyleDefault)))
}

func (a *adapterAdNode) setAdapterEventData(ev bluetooth.AdapterEventData) {
	a.adapter.AdapterEventData = ev
}

func (a *adapterAdNode) actionPowered() (opCreationInfo, opInvoker) {
	return newOpCreationInfo("Changing power state", "Message"), func(ov *opRunningInfo) tview.Msg {
		time.Sleep(1 * time.Second)
		ov.info("Updated Message 1")
		time.Sleep(1 * time.Second)
		ov.info("Updated Message 2")
		time.Sleep(1 * time.Second)

		return ov.opSuccess(actionStateDisabled)
	}
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
	*adTreeNode
	emptyNoder

	device bluetooth.DeviceData
}

func newDeviceAdNode(adtree *adTreeModel, parentNode *adTreeNode, device bluetooth.DeviceData) *adTreeNode {
	dn := &deviceAdNode{}

	deviceNode := newAdNode(
		nodeTypeDevice, getDeviceDisplayName(device.DeviceEventData),
		true, newDeviceNodeID(device.DeviceAddress),
		dn, adtree,
	)

	dn.device = device
	dn.adTreeNode = deviceNode

	dn.addActionsList()
	dn.populateActions()

	parentNode.AddChild(dn.Node)

	return deviceNode
}

func (d *deviceAdNode) handleKeys(p tview.KeyMsg) (routerMsg, bool) {
	return actionHandleKeyMsg(p, d.adTreeNode)
}

func (d *deviceAdNode) populateActions() {
	actionsListNode := d.Node.Children()[relPosActionsListNode]
	actionsListNode.SetChildren(nil)

	node := getNodeReference(actionsListNode)
	if node == nil {
		return
	}

	node.addAction(kb().Device.ToggleConnection, boolToActionState(d.device.Connected.Value()), true, d.actionConnect)
	node.addAction(kb().Device.TogglePairedState, boolToActionState(d.device.Paired.Value()), true, d.actionPair)

	if d.tree.features().Has(appfeatures.FeatureSendFile, appfeatures.FeatureReceiveFile) &&
		d.device.HaveService(bluetooth.ObexObjpushServiceClass) {
		node.addAction(kb().Device.SendFiles, actionStateNone, false, d.actionSend)
	}

	if runtime.GOOS == "linux" {
		node.addAction(kb().Device.Trust, boolToActionState(d.device.Trusted.Value()), true, d.actionTrust)
		node.addAction(kb().Device.Block, boolToActionState(d.device.Blocked.Value()), true, d.actionBlock)

		if d.device.HaveService(bluetooth.AudioSourceServiceClass) ||
			d.device.HaveService(bluetooth.AudioSinkServiceClass) {
			node.addAction(kb().Device.AudioProfiles, actionStateNone, false, d.actionAudioProfiles)
		}

		if d.device.HaveService(bluetooth.AudioSourceServiceClass) &&
			d.device.HaveService(bluetooth.AvRemoteServiceClass) &&
			d.device.HaveService(bluetooth.AvRemoteTargetServiceClass) {
			node.addAction(kb().Player.ToggleDisplay, actionStateDisabled, true, d.actionMediaPlayer)
		}

		if d.tree.features().Has(appfeatures.FeatureNetwork) &&
			d.device.HaveService(bluetooth.NapServiceClass) &&
			(d.device.HaveService(bluetooth.PanuServiceClass) ||
				d.device.HaveService(bluetooth.DialupNetServiceClass)) {
			node.addAction(kb().Device.NetworkOptions, actionStateNone, false, d.actionNetwork)
		}
	}

	for _, actionNode := range actionsListNode.Children() {
		d.updateActionNode(getNodeReference(actionNode), emptyActionUpdateMsg())
	}
}

func (d *deviceAdNode) updateActionNode(actionNode *adTreeNode, updateMsg actionUpdateMsg) {
	adnode := actionNode
	state := adnode.actionState

	updateMsg.updateState(actionNode.id.String(), state)

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

	adnode.SetLine(text.NewLine(text.NewSegment(actionText, tcell.StyleDefault)))
}

func (d *deviceAdNode) setDeviceEventData(ev bluetooth.DeviceEventData) {
	d.device.DeviceEventData = ev
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

type emptyNoder struct{}

func (e *emptyNoder) handleKeys(tview.KeyMsg) (routerMsg, bool) {
	return routerMsg{}, false
}

func (e *emptyNoder) populateActions() {
}

func (e *emptyNoder) updateActionNode(*adTreeNode, actionUpdateMsg) {
}

func (e *emptyNoder) setAdapterEventData(bluetooth.AdapterEventData) {
}

func (e *emptyNoder) setDeviceEventData(bluetooth.DeviceEventData) {
}

func findTreeNodeByID(rootNode *rootAdNode, id string, nodeType adTreeNodeType) (*adTreeNode, int, bool) {
	currNode := rootNode.adTreeNode
	currNodePos := -1

	for nodeID := range iterNodeID(id) {
		found := false

		for pos, n := range currNode.Children() {
			node := getNodeReference(n)
			if node.id.String() == nodeID {
				currNode = node
				currNodePos = pos

				found = true
				break
			}
		}

		if !found {
			return nil, -1, false
		}
	}

	return currNode, currNodePos, currNode != nil && currNode.nodeType == nodeType
}

func iterNodeID(id string) iter.Seq[string] {
	return func(yield func(string) bool) {
		currIdx := 0

		for currIdx < len(id) {
			idx := strings.Index(id[currIdx:], "/")
			if idx < 0 {
				break
			}

			currIdx += idx + 1
			nodeID := id[:currIdx-1]

			if !yield(nodeID) {
				return
			}
		}

		if !yield(id) {
			return
		}
	}
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

// getAdapterDisplayName returns the display name of the adapter.
func getAdapterDisplayName(adapterData bluetooth.AdapterData) string {
	if name, ok := adapterData.Name.Get(); ok {
		return name
	}

	if adapterData.UniqueName != "" {
		return adapterData.UniqueName
	}

	return adapterData.Address.String()
}

// getDeviceDisplayName returns the display name for the device.
func getDeviceDisplayName(deviceData bluetooth.DeviceEventData) string {
	if name, ok := deviceData.Name.Get(); ok {
		return name
	}

	if alias, ok := deviceData.Alias.Get(); ok {
		return alias
	}

	return deviceData.Address.String()
}

func useStringBuffer(size int, fn func(b *strings.Builder)) string {
	var sb strings.Builder
	sb.Grow(size)

	fn(&sb)

	return sb.String()
}
