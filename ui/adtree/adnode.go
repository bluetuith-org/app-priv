package adtree

import (
	"strconv"

	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/gdamore/tcell/v3"
)

// adTreeNodeType specifies the type of the tree node.
type adTreeNodeType uint8

// The different types of tree nodes.
const (
	nodeTypeRoot adTreeNodeType = iota
	nodeTypeAdapter
	nodeTypeDevice
	nodeTypeDevicesList
	nodeTypeAction
	nodeTypeActionsList
)

// adSubNodePos specifies the position of sub-nodes under
// adapter, device or root nodes.
type adSubNodePos uint8

// The different types of sub nodes.
const (
	relPosActionsListNode adSubNodePos = iota
	relPosDevicesListNode
)

// adNode describes a node in an adapter-device tree.
// noder must not be nil.
type adNode struct {
	*ctree.Node

	id       nodeID
	nodeType adTreeNodeType

	noder       adNoder
	actionState *adActionState

	name string

	rv ui.RootView
}

func newAdNode(ntype adTreeNodeType, name string, expanded bool, id nodeID, noder adNoder, rv ui.RootView) *adNode {
	data := &adNode{}

	node := ctree.NewNode(data).SetIndent(5).SetExpanded(expanded)
	data.Node = node
	data.SetName(name)

	data.id = id
	data.nodeType = ntype
	data.noder = noder
	data.rv = rv

	return data
}

// ID returns the identifier for this node.
func (a *adNode) ID() string {
	return a.id.String()
}

// Name returns the name of the node.
func (a *adNode) Name() string {
	return a.name
}

// SetName sets the name for the node.
func (a *adNode) SetName(name string) {
	a.name = name
}

func (a *adNode) ParentNode() (*adNode, bool) {
	parent, ok := a.Node.Parent()
	if !ok {
		return nil, false
	}

	return parent.Data[*adNode](), true
}

// GetStyle returns a style to be applied for this node's label.
func (a *adNode) GetStyle() tcell.Style {
	switch a.nodeType {
	case nodeTypeAction:
		return theme.Current().ADTree.ActionsList.Nodes

	case nodeTypeActionsList:
		return theme.Current().ADTree.ActionsList.Style

	case nodeTypeDevicesList:
		return theme.Current().ADTree.DevicesList
	}

	return a.noder.NodeStyle()
}

// GetIcon returns an icon for this node.
func (a *adNode) GetIcon() string {
	switch a.nodeType {
	case nodeTypeAction:
		return theme.Icons().Actions.String()

	case nodeTypeActionsList:
		return theme.Icons().Actions.String()

	case nodeTypeDevicesList:
		return theme.Icons().Devices.String()
	}

	return a.noder.Icon()
}

// AddDevicesList adds a devices-list node, under the adapter node.
func (a *adNode) AddDevicesList(devices []bluetooth.DeviceData) {
	dlAdNode := newAdNode(
		nodeTypeDevicesList, "Devices",
		true, a.id.appendSubNodeNib(nibDevicesList),
		a.noder, a.rv,
	)

	for _, device := range devices {
		newDeviceNode(a.rv, dlAdNode, device)
	}

	a.AddChild(dlAdNode.Node)
}

// AddActionsList adds an actions-list node to the tree,
// under adapter or device nodes.
func (a *adNode) AddActionsList() {
	alAdNode := newAdNode(
		nodeTypeActionsList, "Actions",
		false, a.id.appendSubNodeNib(nibActionsList),
		a.noder, a.rv,
	)

	alAdNode.SetChildren(make([]*ctree.Node, 0, 5))

	a.AddChild(alAdNode.Node)
}

// AddAction adds an action node to the tree, under the actions-list node.
func (a *adNode) AddAction(key keybindings.Keybinding, state actionStateSpec, isToggleable bool, invoker actionInvoker) {
	actionNode := newAdNode(
		nodeTypeAction, "",
		false, a.id.appendSubNodeTextNib(nibAction, strconv.Itoa(int(key.ID))),
		a.noder, a.rv,
	)

	actionNode.actionState = newAdActionState(key, state, isToggleable, invoker)

	a.AddChild(actionNode.Node)
}
