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
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/app/components/treeview"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
	"github.com/bluetuith-org/bluetuith/ui/theme"
	"github.com/gdamore/tcell/v3"
)

// adTreeModel represents a tree of adapters and its associated devices.
type adTreeModel struct {
	*treeview.Model[adTreeNode]

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

	model := treeview.NewModel(newAdTreeProvider())
	model.SetPrefixes(prefixes)
	model.SetRoot(a.rootNode.treeNode)
	model.SetCurrentNode(model.Root())
	model.SetMarkers(treeview.Markers{
		Expanded:  theme.Icons().TriangleDown.String() + " ",
		Collapsed: theme.Icons().TriangleRight.String() + " ",
		Leaf:      " ",
	})

	a.Model = model

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
}

// RefreshContent refreshes the content of the view.
func (a *adTreeModel) RefreshContent() {
}

// Update receives messages when this model has focus.
func (a *adTreeModel) Update(msg tview.Msg) tview.Cmd {
	switch m := msg.(type) {
	case tview.InitMsg:
		return a.populate()

	case tview.KeyMsg:
		switch {
		case kb().NavigateUp.Matches(m):
			a.Move(-1)

		case kb().NavigateDown.Matches(m):
			a.Move(1)

		case kb().ADTree.ToggleNodes.Matches(m):
			a.CurrentNode().SetExpanded(!a.CurrentNode().Expanded())
		}

	default:
	}

	return a.Model.Update(msg)
}

// View draws this model onto the screen.
func (a *adTreeModel) View(screen tview.Screen) {
	a.Model.View(screen)
}

func (a *adTreeModel) populate() tview.Cmd {
	return func() tview.Msg {
		adapters, err := a.rv.Session().Adapters()
		if err != nil {
			// TODO: Log error messages
			return nil
		}

		for _, adapter := range adapters {
			if err := a.addAdapter(adapter); err != nil {
				_ = err
				return nil
			}
		}

		return msgAdTreeUpdate()
	}
}

func (a *adTreeModel) addAdapter(adapter bluetooth.AdapterData) error {
	devices, err := a.rv.Session().Adapter(adapter.AdapterAddress).Devices()
	if err != nil {
		// TODO:Log error messages
		_ = err
		return err
	}

	a.rootNode.addAdapter(adapter, devices)

	return nil
}

func (a *adTreeModel) addDevice(device bluetooth.DeviceData) {
	a.rootNode.addDevice(device)
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

type adTreeProvider struct{}

func newAdTreeProvider() *adTreeProvider {
	return &adTreeProvider{}
}

func (a *adTreeProvider) Icon(node *treeview.Node[adTreeNode]) string {
	return node.Data().getIcon()
}

func (a *adTreeProvider) Label(node *treeview.Node[adTreeNode]) string {
	return node.Name()
}

func (a *adTreeProvider) NodeStyles(node *treeview.Node[adTreeNode], isFocused bool) (prefixStyle tcell.Style, markerStyle tcell.Style, textStyle tcell.Style) {
	if isFocused {
		selectionStyle := theme.Current().ADTree.Selection
		return selectionStyle, selectionStyle, selectionStyle
	}

	markerStyle = theme.Current().ADTree.Indicators
	nodeStyle := node.Data().getStyle()

	return markerStyle, markerStyle, nodeStyle
}

func (a *adTreeProvider) TreeStyle() (bgStyle tcell.Style, graphicsStyle tcell.Style) {
	return theme.Current().ADTree.Bg, theme.Current().ADTree.Heading
}

var _ treeview.Provider[adTreeNode] = (*adTreeProvider)(nil)

type adTreeNodeType uint8

const (
	nodeTypeRoot adTreeNodeType = iota
	nodeTypeAdapter
	nodeTypeDevice
	nodeTypeDevicesList
	nodeTypeAction
	nodeTypeActionsList
	_nodeTypeMax
)

type adSubNodePos uint8

const (
	relPosActionsListNode adSubNodePos = iota
	relPosDevicesListNode
	_relPosMax
)

type adTreeNode struct {
	id nodeID

	nodeType adTreeNodeType
	treeNode *treeview.Node[adTreeNode]

	noder       adNoder
	actionState *adActionState

	tree *adTreeModel
}

func newAdNode(ntype adTreeNodeType, name string, expanded bool, id nodeID, noder adNoder, tree *adTreeModel) *adTreeNode {
	node := treeview.NewNode(id.String(), "", adTreeNode{})
	node.SetName(name)
	node.SetExpanded(expanded)

	data := node.Data()

	data.id = id
	data.treeNode = node
	data.nodeType = ntype
	data.noder = noder
	data.tree = tree

	return data
}

func (a *adTreeNode) getStyle() tcell.Style {
	switch a.nodeType {
	case nodeTypeAction:
		return theme.Current().ADTree.ActionsList.Nodes

	case nodeTypeActionsList:
		return theme.Current().ADTree.ActionsList.Style

	case nodeTypeDevicesList:
		return theme.Current().ADTree.DevicesList
	}

	if a.noder == nil {
		return tcell.Style{}
	}

	return a.noder.nodeStyle()
}

func (a *adTreeNode) getIcon() string {
	switch a.nodeType {
	case nodeTypeAction:
		return theme.Icons().Actions.String()

	case nodeTypeActionsList:
		return theme.Icons().Actions.String()

	case nodeTypeDevicesList:
		return theme.Icons().Devices.String()
	}

	if a.noder == nil {
		return ""
	}

	return a.noder.icon()
}

func (a *adTreeNode) addDevicesList(devices []bluetooth.DeviceData) {
	dlAdNode := newAdNode(
		nodeTypeDevicesList, "Devices",
		true, a.id.appendSubNodeNib(nibDevicesList),
		a.noder, a.tree,
	)

	for _, device := range devices {
		newDeviceAdNode(a.tree, dlAdNode.treeNode, device)
	}

	a.treeNode.AddChild(dlAdNode.treeNode)
}

func (a *adTreeNode) addActionsList() {
	alAdNode := newAdNode(
		nodeTypeActionsList, "Actions",
		false, a.id.appendSubNodeNib(nibActionsList),
		a.noder, a.tree,
	)

	alAdNode.treeNode.SetChildren(make([]*treeview.Node[adTreeNode], 0, 5))

	a.treeNode.AddChild(alAdNode.treeNode)
}

func (a *adTreeNode) addAction(key keybindings.Keybinding, state actionStateSpec, isToggleable bool, invoker actionInvoker) {
	actionNode := newAdNode(
		nodeTypeAction, "",
		false, a.id.appendSubNodeTextNib(nibAction, strconv.Itoa(int(key.ID))),
		a.noder, a.tree,
	)

	actionNode.actionState = newAdActionState(key, state, isToggleable, invoker)

	a.treeNode.AddChild(actionNode.treeNode)
}

type adNoder interface {
	handleKeys(p tview.KeyMsg) (routerMsg, bool)
	nodeStyle() tcell.Style
	icon() string
	populateActions()
	updateActionNode(actionNode *treeview.Node[adTreeNode], updateMsg actionUpdateMsg)
	setAdapterEventData(ev bluetooth.AdapterEventData)
	setDeviceEventData(ev bluetooth.DeviceEventData)
}

type rootAdNode struct {
	*adTreeNode
	emptyNoder

	sync.RWMutex
}

func newRootAdNode(tree *adTreeModel) *rootAdNode {
	rn := &rootAdNode{}

	rootAdNode := newAdNode(
		nodeTypeRoot, "Adapters",
		true, newRootNodeID(),
		rn, tree,
	)

	rootAdNode.treeNode.SetChildren(make([]*treeview.Node[adTreeNode], 0, 10))
	rn.adTreeNode = rootAdNode

	return rn
}

func (r *rootAdNode) nodeStyle() tcell.Style {
	return theme.Current().ADTree.Heading
}

func (r *rootAdNode) icon() string {
	return theme.Icons().Adapters.String()
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
		parentNode := r.treeNode
		if parentNode != nil {
			parentNode.SetChildren(slices.Delete(parentNode.Children(), pos, pos+1))
		}

		return
	}

	adapterNode.Data().noder.setAdapterEventData(adapterEvent)
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
		parentNode := deviceNode.Parent()
		if parentNode != nil {
			parentNode.SetChildren(slices.Delete(parentNode.Children(), pos, pos+1))
		}

		return
	}

	deviceNode.Data().noder.setDeviceEventData(deviceEvent)
}

func (r *rootAdNode) updateAction(updateMsg actionUpdateMsg) {
	r.Lock()
	defer r.Unlock()

	node, _, ok := findTreeNodeByID(r, updateMsg.id, nodeTypeAction)
	if !ok {
		return
	}

	node.Data().noder.updateActionNode(node, updateMsg)
}

var _ adNoder = (*rootAdNode)(nil)

type adapterAdNode struct {
	*adTreeNode
	emptyNoder

	adapter bluetooth.AdapterData
}

func newAdapterAdNode(tree *adTreeModel, rootNode *rootAdNode, adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) *adTreeNode {
	an := &adapterAdNode{}

	adapterNode := newAdNode(
		nodeTypeAdapter, getAdapterDisplayName(adapter),
		true, newAdapterNodeID(adapter.AdapterAddress),
		an, tree,
	)

	an.adTreeNode = adapterNode
	an.adapter = adapter

	an.addActionsList()
	an.addDevicesList(devices)

	an.populateActions()

	rootNode.treeNode.AddChild(adapterNode.treeNode)

	return adapterNode
}

func (a *adapterAdNode) nodeStyle() tcell.Style {
	return theme.Current().ADTree.Adapter
}

func (a *adapterAdNode) icon() string {
	return theme.Icons().Adapter.String()
}

func (a *adapterAdNode) handleKeys(p tview.KeyMsg) (routerMsg, bool) {
	return actionHandleKeyMsg(p, a.treeNode)
}

func (a *adapterAdNode) populateActions() {
	actionsListNode := a.treeNode.Children()[relPosActionsListNode]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data()
	if node == nil {
		return
	}

	node.addAction(kb().Adapter.TogglePower, boolToActionState(a.adapter.Powered.Value()), true, a.actionPowered)
	node.addAction(kb().Adapter.ToggleDiscoverable, boolToActionState(a.adapter.Discovering.Value()), true, a.actionDiscoverable)
	node.addAction(kb().Adapter.TogglePairable, boolToActionState(a.adapter.Pairable.Value()), true, a.actionPairable)
	node.addAction(kb().Adapter.ToggleScan, boolToActionState(a.adapter.Discovering.Value()), true, a.actionScan)

	for _, actionNode := range actionsListNode.Children() {
		a.updateActionNode(actionNode, emptyActionUpdateMsg())
	}
}

func (a *adapterAdNode) updateActionNode(actionNode *treeview.Node[adTreeNode], updateMsg actionUpdateMsg) {
	adnode := actionNode.Data()
	state := adnode.actionState

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

var _ adNoder = (*adapterAdNode)(nil)

type deviceAdNode struct {
	*adTreeNode
	emptyNoder

	device bluetooth.DeviceData
}

func newDeviceAdNode(tree *adTreeModel, parentNode *treeview.Node[adTreeNode], device bluetooth.DeviceData) *adTreeNode {
	dn := &deviceAdNode{}

	deviceNode := newAdNode(
		nodeTypeDevice, getDeviceDisplayName(device.DeviceEventData),
		true, newDeviceNodeID(device.DeviceAddress),
		dn, tree,
	)

	dn.device = device
	dn.adTreeNode = deviceNode

	dn.addActionsList()
	dn.populateActions()

	parentNode.AddChild(dn.treeNode)

	return deviceNode
}

func (d *deviceAdNode) nodeStyle() tcell.Style {
	paired, ok := d.device.Paired.Get()
	if !ok {
		return theme.Current().ADTree.Device.Unknown
	}
	if paired {
		return theme.Current().ADTree.Device.Paired
	}

	return theme.Current().ADTree.Device.Discovered
}

func (d *deviceAdNode) icon() string {
	return theme.DeviceClassToIcon(theme.Icons(), d.device.Class).String()
}

func (d *deviceAdNode) handleKeys(p tview.KeyMsg) (routerMsg, bool) {
	return actionHandleKeyMsg(p, d.treeNode)
}

func (d *deviceAdNode) populateActions() {
	actionsListNode := d.treeNode.Children()[relPosActionsListNode]
	actionsListNode.SetChildren(nil)

	node := actionsListNode.Data()
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
		d.updateActionNode(actionNode, emptyActionUpdateMsg())
	}
}

func (d *deviceAdNode) updateActionNode(actionNode *treeview.Node[adTreeNode], updateMsg actionUpdateMsg) {
	adnode := actionNode.Data()
	state := adnode.actionState

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

	adnode.treeNode.SetName(actionText)
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

var _ adNoder = (*deviceAdNode)(nil)

type emptyNoder struct{}

func (e *emptyNoder) handleKeys(_ tview.KeyMsg) (routerMsg, bool) {
	return emptyRouterMsg(), false
}

func (e *emptyNoder) populateActions() {
}

func (e *emptyNoder) updateActionNode(_ *treeview.Node[adTreeNode], _ actionUpdateMsg) {
}

func (e *emptyNoder) setAdapterEventData(_ bluetooth.AdapterEventData) {
}

func (e *emptyNoder) setDeviceEventData(_ bluetooth.DeviceEventData) {
}

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

func (n nodeID) appendSubNodeTextNib(nib nodeIDNib, text string) nodeID {
	length := len(n.id) + lenNibPlusColon + len(text)

	n.id = useStringBuffer(length, func(b *strings.Builder) {
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

func findTreeNodeByID(rootNode *rootAdNode, id string, nodeType adTreeNodeType) (*treeview.Node[adTreeNode], int, bool) {
	currNode := rootNode.treeNode
	currNodePos := -1

	for nodeID := range iterNodeID(id) {
		found := false

		for pos, node := range currNode.Children() {
			if node.ID() == nodeID {
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

	return currNode, currNodePos, currNode != nil && currNode.Data().nodeType == nodeType
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

var _noStyle = tcell.Style{}

func overrideStyleByType(node *adTreeNode) (tcell.Style, bool) {
	switch node.nodeType {
	case nodeTypeAction:
		return theme.Current().ADTree.ActionsList.Nodes, true

	case nodeTypeActionsList:
		return theme.Current().ADTree.ActionsList.Style, true

	case nodeTypeDevicesList:
		return theme.Current().ADTree.DevicesList, true
	}

	return _noStyle, false
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
