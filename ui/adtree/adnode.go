package adtree

import (
	"github.com/ayn2op/tview/richtext"
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
	nodeTypeStatus
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

	content richtext.Line

	rv ui.RootView
}

func newAdNode(ntype adTreeNodeType, expanded bool, id nodeID, noder adNoder, rv ui.RootView) *adNode {
	data := &adNode{}

	node := ctree.NewNode(data).SetIndent(5).SetExpanded(expanded)
	data.Node = node

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

// Noder returns the interface that implements [adNoder].
func (a *adNode) Noder() adNoder {
	return a.noder
}

// Content returns the name of the node.
func (a *adNode) Content() richtext.Line {
	return a.content
}

// SetContent sets the name for the node.
func (a *adNode) SetContent(content richtext.Line) {
	a.content = content
}

func (a *adNode) ParentNode() (*adNode, bool) {
	parent, ok := a.Node.Parent()
	if !ok {
		return nil, false
	}

	return parent.Data[*adNode](), true
}

func buildLabel(icon, name string, style tcell.Style) richtext.Line {
	return richtext.NewLine(
		richtext.NewSegment(icon, style),
		richtext.NewSegment(" ", style),
		richtext.NewSegment(name, style),
	)
}

func pushSeg(line richtext.Line, bgStyle tcell.Style, yield func(richtext.Segment) bool) bool {
	if !yield(richtext.NewSegment(" ", bgStyle)) {
		return false
	}

	for _, seg := range line {
		if !yield(seg) {
			return false
		}
	}

	return true
}
