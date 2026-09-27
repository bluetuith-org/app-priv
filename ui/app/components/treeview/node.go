package treeview

import "slices"

// Node represents one node in a tree view.
type Node[T any] struct {
	// The reference object.
	reference T

	// This node's child nodes.
	children []*Node[T]

	// Whether or not this node can be selected.
	selectable bool

	// Whether or not this node's children should be displayed.
	expanded bool

	// Whether or not this node can be expanded, even if children are not loaded yet.
	expandable bool

	// The additional horizontal indent of this node's text.
	indent int

	id, name string
	parent   *Node[T]
}

// NewNode creates a new tree node.
func NewNode[T any](id, name string, data T) *Node[T] {
	return &Node[T]{
		id:         id,
		name:       name,
		reference:  data,
		indent:     5,
		selectable: true,
		expanded:   true,
		expandable: false,
	}
}

// ID returns the node's ID.
func (n *Node[T]) ID() string {
	return n.id
}

// Name returns the node's name.
func (n *Node[T]) Name() string {
	return n.name
}

// SetName sets the name for this node.
func (n *Node[T]) SetName(name string) *Node[T] {
	n.name = name
	return n
}

// Parent returns the parent of this node.
func (n *Node[T]) Parent() *Node[T] {
	return n.parent
}

// Walk traverses this node's subtree in depth-first, pre-order (NLR) order and
// calls the provided callback function on each traversed node (which includes
// this node) with the traversed node and its parent node (nil for this node).
// The callback returns whether traversal should continue with the traversed
// node's child nodes (true) or not recurse any deeper (false).
func (n *Node[T]) Walk(callback func(node, parent *Node[T]) bool) *Node[T] {
	type entry struct{ node, parent *Node[T] }
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
func (n *Node[T]) Data() *T {
	return &n.reference
}

// SetData allows you to store a reference of any type in this node. This
// will allow you to establish a mapping between the Model hierarchy and your
// internal tree structure.
func (n *Node[T]) SetData(reference T) *Node[T] {
	n.reference = reference

	return n
}

// Children returns this node's children.
func (n *Node[T]) Children() []*Node[T] {
	return n.children
}

// SetChildren sets this node's child nodes.
func (n *Node[T]) SetChildren(childNodes []*Node[T]) *Node[T] {
	for _, ch := range childNodes {
		ch.parent = n
	}

	n.children = childNodes
	return n
}

// ClearChildren removes all child nodes from this node.
func (n *Node[T]) ClearChildren() *Node[T] {
	n.children = nil
	return n
}

// AddChild adds a new child node to this node.
func (n *Node[T]) AddChild(node *Node[T]) *Node[T] {
	node.parent = n
	n.children = append(n.children, node)

	return n
}

// RemoveChild removes a child node from this node. If the child node cannot be
// found, nothing happens.
func (n *Node[T]) RemoveChild(node *Node[T]) *Node[T] {
	if index := slices.Index(n.children, node); index >= 0 {
		n.children = slices.Delete(n.children, index, index+1)
	}
	return n
}

// SetSelectable sets a flag indicating whether this node can be selected by
// the user.
func (n *Node[T]) SetSelectable(selectable bool) *Node[T] {
	n.selectable = selectable
	return n
}

// Expanded returns whether the child nodes of this node are visible.
func (n *Node[T]) Expanded() bool {
	return n.expanded
}

// SetExpanded sets whether or not this node's child nodes should be displayed.
func (n *Node[T]) SetExpanded(expanded bool) *Node[T] {
	n.expanded = expanded
	return n
}

// Expandable returns whether this node can be expanded even when there are
// no loaded child nodes yet.
func (n *Node[T]) Expandable() bool {
	return n.expandable
}

// SetExpandable sets whether this node can be expanded even when there are no
// loaded child nodes yet.
func (n *Node[T]) SetExpandable(expandable bool) *Node[T] {
	n.expandable = expandable
	return n
}

// Expand makes the child nodes of this node appear.
func (n *Node[T]) Expand() *Node[T] {
	n.expanded = true
	return n
}

// Collapse makes the child nodes of this node disappear.
func (n *Node[T]) Collapse() *Node[T] {
	n.expanded = false
	return n
}

// ExpandAll expands this node and all descendent nodes.
func (n *Node[T]) ExpandAll() *Node[T] {
	n.Walk(func(node, _ *Node[T]) bool {
		node.expanded = true
		return true
	})
	return n
}

// CollapseAll collapses this node and all descendent nodes.
func (n *Node[T]) CollapseAll() *Node[T] {
	n.Walk(func(node, _ *Node[T]) bool {
		node.expanded = false
		return true
	})
	return n
}

// SetIndent sets an additional indentation for this node's text. A value of 0
// keeps the text as far left as possible with a minimum of line graphics. Any
// value greater than that moves the text to the right.
func (n *Node[T]) SetIndent(indent int) *Node[T] {
	n.indent = indent
	return n
}
