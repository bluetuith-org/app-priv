package ctree

import (
	"slices"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
)

// Node represents one node in a tree view.
type Node struct {
	// The reference object.
	reference any

	// This node's child nodes.
	children []*Node

	// The parent node of this node.
	parent *Node

	// The item's text.
	line richtext.Line
	// lineWidth caches the width of line.
	lineWidth int

	// The style of selected text.
	selectedTextStyle tview.Style

	// Whether or not this node can be selected.
	selectable bool

	// Whether or not this node's children should be displayed.
	expanded bool

	// Whether or not this node can be expanded, even if children are not loaded yet.
	expandable bool

	// The additional horizontal indent of this node's text.
	indent int

	vGapTop    int
	vGapBottom int
}

// NewNode returns a new tree node.
func NewNode(ref any) *Node {
	return &Node{
		reference:  ref,
		indent:     2,
		expanded:   true,
		selectable: true,
	}
}

// Walk traverses this node's subtree in depth-first, pre-order (NLR) order and calls the provided callback function on each traversed node (which includes this node) with the traversed node and its parent node (nil for this node). The callback returns whether traversal should continue with the traversed node's child nodes (true) or not recurse any deeper (false).
func (n *Node) Walk(callback func(node, parent *Node) bool) *Node {
	type entry struct{ node, parent *Node }
	stack := []entry{{node: n}}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		if !callback(current.node, current.parent) {
			continue
		}
		for _, v := range slices.Backward(current.node.children) {
			stack = append(stack, entry{v, current.node})
		}
	}
	return n
}

// Data returns this node's reference object.
//
//revive:disable
func (n *Node) Data[T any]() T {
	return n.reference.(T)
}

// Parent returns this node's parent node.
func (n *Node) Parent() (*Node, bool) {
	return n.parent, n.parent != nil
}

// Children returns this node's children.
func (n *Node) Children() []*Node {
	return n.children
}

// SetChildren sets this node's child nodes.
func (n *Node) SetChildren(childNodes []*Node) *Node {
	for _, ch := range childNodes {
		ch.parent = n
	}

	n.children = childNodes
	return n
}

// ClearChildren removes all child nodes from this node.
func (n *Node) ClearChildren() *Node {
	n.children = nil
	return n
}

// AddChild adds a new child node to this node.
func (n *Node) AddChild(node *Node) *Node {
	node.parent = n
	n.children = append(n.children, node)

	return n
}

// SetVGap sets the top and bottom vertical spacing (in rows) for the node.
func (n *Node) SetVGap(top, bottom int) *Node {
	n.vGapTop = top
	n.vGapBottom = bottom
	return n
}

// VGapTop returns the top vertical gap.
func (n *Node) VGapTop() int {
	return n.vGapTop
}

// VGapBottom returns the bottom vertical gap.
func (n *Node) VGapBottom() int {
	return n.vGapBottom
}

// Reference returns this node's reference object.
func (n *Node) Reference() any {
	return n.reference
}

// SetReference allows you to store a reference of any type in this node. This will allow you to establish a mapping between the Model hierarchy and your internal tree structure.
func (n *Node) SetReference(reference any) *Node {
	n.reference = reference

	return n
}

// Line returns the node's styled text line.
func (n *Node) Line() richtext.Line {
	return n.line
}

// SetLine sets the node's styled text line.
func (n *Node) SetLine(line richtext.Line) *Node {
	n.line, n.lineWidth = line, line.Width()

	return n
}

// RemoveChild removes a child node from this node. If the child node cannot be found, nothing happens.
func (n *Node) RemoveChild(node *Node) *Node {
	if index := slices.Index(n.children, node); index >= 0 {
		n.children = slices.Delete(n.children, index, index+1)
	}
	return n
}

// SetSelectable sets a flag indicating whether this node can be selected by the user.
func (n *Node) SetSelectable(selectable bool) *Node {
	n.selectable = selectable
	return n
}

// Expanded returns whether the child nodes of this node are visible.
func (n *Node) Expanded() bool {
	return n.expanded
}

// SetExpanded sets whether or not this node's child nodes should be displayed.
func (n *Node) SetExpanded(expanded bool) *Node {
	n.expanded = expanded
	return n
}

// Expandable returns whether this node can be expanded even when there are no loaded child nodes yet.
func (n *Node) Expandable() bool {
	return n.expandable
}

// SetExpandable sets whether this node can be expanded even when there are no loaded child nodes yet.
func (n *Node) SetExpandable(expandable bool) *Node {
	n.expandable = expandable
	return n
}

// Expand makes the child nodes of this node appear.
func (n *Node) Expand() *Node {
	n.expanded = true
	return n
}

// Collapse makes the child nodes of this node disappear.
func (n *Node) Collapse() *Node {
	n.expanded = false
	return n
}

// ExpandAll expands this node and all descendent nodes.
func (n *Node) ExpandAll() *Node {
	n.Walk(func(node, parent *Node) bool {
		node.expanded = true
		return true
	})
	return n
}

// CollapseAll collapses this node and all descendent nodes.
func (n *Node) CollapseAll() *Node {
	n.Walk(func(node, parent *Node) bool {
		node.expanded = false
		return true
	})
	return n
}

// SelectedTextStyle returns the text style for this node when it is selected.
func (n *Node) SelectedTextStyle() tview.Style {
	return n.selectedTextStyle
}

// SetSelectedTextStyle sets the text style for this node when it is selected.
func (n *Node) SetSelectedTextStyle(style tview.Style) *Node {
	n.selectedTextStyle = style
	return n
}

// SetIndent sets an additional indentation for this node's text. A value of 0 keeps the text as far left as possible with a minimum of line graphics. Any value greater than that moves the text to the right.
func (n *Node) SetIndent(indent int) *Node {
	n.indent = indent
	return n
}

// PathTo returns the nodes from n to node, both included, or nil if node is not under n.
func (n *Node) PathTo(node *Node) []*Node {
	if n == node {
		return []*Node{n}
	}
	for _, child := range n.children {
		if path := child.PathTo(node); path != nil {
			return append([]*Node{n}, path...)
		}
	}
	return nil
}
