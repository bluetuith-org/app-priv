// Package tree shows a tree of nodes that can be expanded, collapsed, and selected.
package ctree

import (
	"slices"

	"github.com/ayn2op/tview/layout"

	"github.com/ayn2op/tview"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// Markers are drawn before a node's text depending on whether it is expanded.
type Markers struct {
	Expanded  string
	Collapsed string
	Leaf      string
}

// Widget draws the nodes under a root and turns keys and the mouse into Changes.
type Widget struct {
	Provider

	root           *Node
	selectionState SelectionState
	width, height  layout.Length
	topLevel       int
	markers        Markers
	graphics       bool
	graphicsSet    tview.BorderSet
	graphicsStyle  tview.Style
	focused        bool
}

var _ tview.Widget = Widget{}

// New returns a tree of the nodes under root, with selectionState as its current node and scroll position, that draws lines between nodes and fills its parent. It is interactive only once OnChange is set.
func New(root *Node, provider Provider, selectionState SelectionState) Widget {
	return Widget{
		Provider:       provider,
		root:           root,
		selectionState: selectionState,
		width:          layout.Fill,
		height:         layout.Fill,
		graphics:       true,
	}
}

// Width sets the width of the tree.
func (w Widget) Width(width layout.Length) Widget {
	w.width = width
	return w
}

// Height sets the height of the tree.
func (w Widget) Height(height layout.Length) Widget {
	w.height = height
	return w
}

// TopLevel sets the first level shown, where the root is level 0, its children level 1, and so on.
func (w Widget) TopLevel(level int) Widget {
	w.topLevel = level
	return w
}

// Markers sets the markers drawn before nodes.
func (w Widget) Markers(markers Markers) Widget {
	w.markers = markers
	return w
}

// Graphics sets whether lines are drawn between nodes.
func (w Widget) Graphics(graphics bool) Widget {
	w.graphics = graphics
	return w
}

// GraphicsSet sets the characters the lines are drawn with.
func (w Widget) GraphicsSet(set tview.BorderSet) Widget {
	w.graphicsSet = set
	return w
}

// GraphicsStyle sets the style of the lines.
func (w Widget) GraphicsStyle(style tview.Style) Widget {
	w.graphicsStyle = style
	return w
}

// Focused sets whether the tree receives keys.
func (w Widget) Focused(focused bool) Widget {
	w.focused = focused
	return w
}

// Size returns the width and height of the tree.
func (w Widget) Size() (width, height layout.Length) {
	return w.width, w.height
}

// view is the rows of the tree and the one of the current node, or -1.
type view struct {
	rows    []row
	current int
}

// resolve returns the rows and the row of the current node, falling back to the first selectable one.
func (w Widget) resolve() view {
	v := view{rows: w.rows(), current: -1}
	if c := w.selectionState.current; c != nil {
		v.current = slices.IndexFunc(v.rows, func(r row) bool { return r.node == c && c.selectable })
		if v.current < 0 {
			v.current = slices.IndexFunc(v.rows, func(r row) bool { return r.node.selectable })
		}
	}
	return v
}

// widest returns the width of the widest of rows.
func (w Widget) widest(rows []row) (widest int) {
	for _, r := range rows {
		widest = max(widest, r.tx+uniseg.StringWidth(w.marker(r.node))+r.node.lineWidth)
	}
	return widest
}

// marker returns the marker drawn before node's text.
func (w Widget) marker(node *Node) string {
	if node.expandable || len(node.children) > 0 {
		if node.expanded {
			return w.markers.Expanded
		}
		return w.markers.Collapsed
	}
	return w.markers.Leaf
}

// node returns the node of row index, or nil if there is none.
func (v view) node(index int) *Node {
	if index < 0 || index >= len(v.rows) {
		return nil
	}
	return v.rows[index].node
}

// row is a node shown on one line.
type row struct {
	node   *Node
	parent int
	gx, tx int
	y      int
}

// rows returns the shown nodes.
func (w Widget) rows() (rows []row) {
	y := 0
	var walk func(node *Node, parent, level, parentX int)
	walk = func(node *Node, parent, level, parentX int) {
		gx, tx := parentX, parentX+node.indent
		if w.graphics {
			tx++
		}
		if level == w.topLevel || level == 0 {
			gx, tx = 0, 0
		}
		if level >= w.topLevel {
			contentY := y + node.vGapTop
			rows = append(rows, row{node: node, parent: parent, gx: gx, tx: tx, y: contentY})
			y = contentY + 1 + node.vGapBottom
			parent = len(rows) - 1
		}
		if node.expanded {
			for _, child := range node.children {
				walk(child, parent, level+1, tx)
			}
		}
	}
	if w.root != nil {
		walk(w.root, -1, 0, 0)
	}
	return rows
}

// vheight returns the total Y height of the tree including all vertical gaps.
func (w Widget) vheight(rows []row) int {
	if len(rows) == 0 {
		return 0
	}
	last := rows[len(rows)-1]
	return last.y + 1 + last.node.vGapBottom
}

// Layout returns the size of the tree within limits.
func (w Widget) Layout(limits layout.Limits) layout.Size {
	return layout.Sized(limits, w.width, w.height, func(layout.Limits) layout.Size {
		rows := w.rows()
		return layout.Size{Width: w.widest(rows), Height: w.vheight(rows)}
	})
}

// Target returns the Y offset and height of the current node's block including its gaps.
func (w Widget) Target(int) (top, height int) {
	v := w.resolve()
	if v.current >= 0 {
		r := v.rows[v.current]
		return r.y - r.node.vGapTop, 1 + r.node.vGapTop + r.node.vGapBottom
	}
	return 0, 0
}

// RowsWidth returns the width of the widest row within Y range [top, top+height).
func (w Widget) RowsWidth(top, height int) int {
	rows := w.rows()
	if len(rows) == 0 || height <= 0 {
		return 0
	}
	bottom := top + height
	var visible []row
	for _, r := range rows {
		rTop := r.y - r.node.vGapTop
		rBottom := r.y + 1 + r.node.vGapBottom
		if rBottom > top && rTop < bottom {
			visible = append(visible, r)
		}
	}
	return w.widest(visible)
}

// nodeAtY returns the node occupying the given relative Y coordinate, or nil if none.
func (v view) nodeAtY(y int) *Node {
	for _, r := range v.rows {
		rTop := r.y - r.node.vGapTop
		rBottom := r.y + 1 + r.node.vGapBottom
		if y >= rTop && y < rBottom {
			return r.node
		}
	}
	return nil
}

// Draw draws the rows with their lines, markers, text, and vertical gap lines.
func (w Widget) Draw(screen tview.Screen, area tview.Rectangle) {
	v := w.resolve()
	x, width, set := area.X, area.Width, w.graphicsSet
	_, bottom := screen.Size()

	for index, current := range v.rows {
		node := current.node
		rTop := current.y - node.vGapTop
		rBottom := current.y + 1 + node.vGapBottom

		if area.Y+rBottom <= 0 || area.Y+rTop >= min(area.Height, bottom-area.Y) {
			continue
		}

		for lineY := rTop; lineY < rBottom; lineY++ {
			y := area.Y + lineY
			if y < 0 || lineY >= area.Height || y >= bottom {
				continue
			}

			if w.graphics {
				// Branches of ancestors that continue past this node
				for ancestor := current.parent; ancestor >= 0 && v.rows[ancestor].parent >= 0; ancestor = v.rows[ancestor].parent {
					a := v.rows[ancestor]
					parent := v.rows[a.parent].node
					if a.gx < width && parent.children[len(parent.children)-1] != a.node {
						screen.Put(x+a.gx, y, set.Right, w.graphicsStyle)
					}
				}

				if lineY < current.y {
					// Top gap: continue vertical branch from parent/siblings above
					if current.parent >= 0 && current.gx < width {
						screen.Put(x+current.gx, y, set.Right, w.graphicsStyle)
					}
				} else if lineY == current.y {
					// Content line: draw tree branch connector and horizontal line
					if current.tx > current.gx && current.gx < width {
						connector := set.BottomLeft
						if current.parent >= 0 {
							if siblings := v.rows[current.parent].node.children; siblings[len(siblings)-1] != node {
								connector = set.LeftT
							}
						}
						screen.Put(x+current.gx, y, connector, w.graphicsStyle)
						for pos := current.gx + 1; pos < current.tx && pos < width; pos++ {
							screen.Put(x+pos, y, set.Top, w.graphicsStyle)
						}
					}
				} else {
					// Bottom gap: continue branch down to next sibling or expanded children
					if current.parent >= 0 && current.gx < width {
						if siblings := v.rows[current.parent].node.children; siblings[len(siblings)-1] != node {
							screen.Put(x+current.gx, y, set.Right, w.graphicsStyle)
						}
					}
					if node.expanded && len(node.children) > 0 && current.tx < width {
						screen.Put(x+current.tx, y, set.Right, w.graphicsStyle)
					}
				}
			}

			// Render marker and node content on the content line
			if lineY == current.y {
				if current.tx >= width {
					continue
				}

				markerStyle := w.Provider.MarkerStyle(node)

				textX := current.tx
				textX += tview.Print(screen, w.marker(node), x+textX, y, width-textX, tview.AlignmentLeft, markerStyle)
				for line, style := range w.StyledLabels(node, index == v.current) {
					for _, segment := range line {
						if textX >= width {
							break
						}

						st := segment.Style
						if style != (tcell.Style{}) {
							st = style
						}

						textX += tview.Print(screen, segment.Text, x+textX, y, width-textX, tview.AlignmentLeft, st)
					}
				}
			}
		}
	}
}

// Handle processes key actions and maps mouse clicks using relative Y coordinates.
func (w Widget) Handle(msg tview.Msg, area tview.Rectangle) tview.Msg {
	if key, ok := msg.(tview.KeyMsg); ok && w.focused {
		if action := w.Provider.Keybind(key); action != ActionNone {
			msg = ActionMsg(action)
		}
	}
	switch m := msg.(type) {
	case ActionMsg:
		if !w.focused {
			return msg
		}
		v := w.resolve()
		current := v.current
		switch Action(m) {
		case ActionDown:
			current = v.step(current, 1)
		case ActionUp:
			current = v.step(current, -1)
		case ActionTop:
			current = v.step(-1, 1)
		case ActionBottom:
			current = v.step(len(v.rows), -1)
		case ActionMoveToParent:
			if current >= 0 {
				if parent := v.rows[current].parent; parent >= 0 && v.rows[parent].node.selectable {
					current = parent
				}
			}
		case ActionSelect:
			return w.Provider.OnSelect(v.node(current))
		default:
			return msg
		}
		return w.Provider.OnChange(Change{current: v.node(current)})
	case tview.MouseMsg:
		x, y := m.Position()
		if m.Action != tview.MouseLeftClick || !area.Contains(x, y) {
			return msg
		}
		if node := w.resolve().nodeAtY(y - area.Y); node != nil && node.selectable {
			return w.Provider.OnSelect(node)
		}
		return nil
	}
	return msg
}

// step returns the next selectable row after index in direction, or index if there is none.
func (v view) step(index, direction int) int {
	for next := index + direction; next >= 0 && next < len(v.rows); next += direction {
		if v.rows[next].node.selectable {
			return next
		}
	}
	return index
}
