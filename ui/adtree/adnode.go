package adtree

import (
	"github.com/bluetuith-org/bluetuith/theme"
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
