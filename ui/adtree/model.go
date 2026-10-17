package adtree

import (
	"iter"
	"time"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/text"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/info"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/gdamore/tcell/v3"
)

const dbTime = 500 * time.Millisecond

// Model represents an adapter-device tree.
type Model struct {
	root        *rootNode
	currentNode *ctree.Node

	scrollState viewport.ScrollState

	graphicsSet tview.BorderSet
	markers     ctree.Markers

	tm        *time.Timer
	reset     bool
	rootident string

	rv ui.RootView
}

// New returns a new adapter-device tree.
func New(rv ui.RootView) Model {
	root := newRootNode(rv)

	graphicsSet := tview.BorderSetRound()
	markers := ctree.Markers{
		Expanded:  theme.Icons().TriangleDown.String() + " ",
		Collapsed: theme.Icons().TriangleRight.String() + " ",
		Leaf:      " ",
	}

	tm := time.NewTimer(dbTime)
	tm.Reset(dbTime)

	return Model{
		root:        root,
		currentNode: root.Node,

		graphicsSet: graphicsSet,
		markers:     markers,

		tm: tm,

		rv: rv,
	}
}

// ViewID returns the view's ID.
func (m *Model) ViewID() ui.ViewID {
	return ui.ViewIDAdTree
}

// Header returns the content to print in the title bar.
func (m *Model) Header() []text.Span {
	return m.root.Title()
}

// Keybind converts keybindings to actions.
func (m *Model) Keybind(msg tview.KeyMsg) ctree.Action {
	switch {
	case ui.Kb().NavigateUp.Matches(msg):
		return ctree.ActionUp

	case ui.Kb().NavigateDown.Matches(msg):
		return ctree.ActionDown

	case ui.Kb().ADTree.ExpandOrSelect.Matches(msg):
		return ctree.ActionSelect
	}

	return ctree.ActionNone
}

// ScrollKeybinds processes a change in the viewport's state and returns a [viewport.Action].
func (m *Model) ScrollKeybinds(msg tview.KeyMsg) viewport.Action {
	switch {
	case ui.Kb().NavigateTop.Matches(msg):
		return viewport.ActionTop

	case ui.Kb().NavigateBottom.Matches(msg):
		return viewport.ActionBottom
	}

	return viewport.ActionNone
}

// OnChange processes a change in the tree's state and returns a [tview.Msg]
func (m *Model) OnChange(chg *ctree.Node) tview.Msg {
	return m.ViewID().RouterMessage(selectionChange(chg))
}

// OnSelect returns a message after processing the selected node.
func (m *Model) OnSelect(node *ctree.Node) tview.Msg {
	return m.ViewID().RouterMessage(selectedMsg{node})
}

// ViewportKeybinds implements keybindings for the viewport.
func (m Model) ViewportKeybinds(msg tview.KeyMsg) viewport.Action {
	switch {
	case ui.Kb().NavigateTop.Matches(msg):
		return viewport.ActionTop

	case ui.Kb().NavigateBottom.Matches(msg):
		return viewport.ActionBottom
	}

	return viewport.ActionNone
}

// OnViewportChange processes a change in the viewport and returns a message.
func (m Model) OnViewportChange(c *viewport.ScrollState) tview.Msg {
	return m.ViewID().RouterMessage(scrollMsg{change: *c})
}

// OnViewportChildMsg processes a message that was passed to its underlying view in the viewport and returns a message.
func (m Model) OnViewportChildMsg(c *viewport.ScrollState, msg tview.Msg) tview.Msg {
	return m.ViewID().RouterMessage(scrollMsg{change: *c, msg: msg})
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return tview.Batch(m.populate, m.dbTimer)
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(focused bool) tview.Widget {
	state := m.currentNode
	if !focused {
		state = nil
	}

	tree := ctree.New(&m, m.root.Node, state).
		TopLevel(1).
		Focused(focused).
		Graphics(true).
		GraphicsSet(m.graphicsSet).
		GraphicsStyle(theme.Current().Global).
		Markers(m.markers)

	vp := viewport.New(tree, m.scrollState).
		Target(tree.Target).
		ContentWidth(tree.RowsWidth).
		Axes(layout.Axes{Width: true, Height: true}).
		Focused(focused).
		Keybind(m.ScrollKeybinds).
		OnChange(m.OnViewportChange).
		OnChildMsg(m.OnViewportChildMsg)

	return box.New(vp).
		Background(theme.Current().Global.GetBackground()).
		Padding(2, 0, 2, 1)
}

// MarkerStyle returns the styles for the markers.
func (m *Model) MarkerStyle(*ctree.Node) tcell.Style {
	return theme.Current().Global
}

// LabelWidth implements [ctree.Provider].
func (m *Model) LabelWidth(node *ctree.Node) int {
	adn := node.Data[*adNode]()

	return adn.labelWidth()
}

// Spacing implements [ctree.Provider].
func (m *Model) Spacing(node *ctree.Node) int {
	adn := node.Data[*adNode]()
	if adn.nodeType == nodeTypeAction {
		return 0
	}

	return 1
}

// StyledLabels returns a sequence of text and their associated styles.
func (m *Model) StyledLabels(node *ctree.Node, selected bool) iter.Seq2[text.Span, tcell.Style] {
	return func(yield func(text.Span, tcell.Style) bool) {
		adn := node.Data[*adNode]()

		for _, span := range adn.Markers() {
			if !yield(span, tcell.Style{}) {
				return
			}
		}

		for _, span := range adn.Content() {
			st := tcell.Style{}
			if selected {
				st = theme.Current().ADTree.Selection
			}

			if !yield(span, st) {
				return
			}
		}
	}
}

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case dbEvent:
		return m.dbHandler()

	case scrollMsg:
		m.scrollState = viewport.ScrollState(ms.change)
		if ms.msg != nil {
			return m.updateModel(ms.msg)
		}

	case selectionChange:
		m.currentNode = ms
		m.scrollState.ScrollToTarget()
		m.resetTimer()
		return m.resetInformation()

	case rootNodeMsg:
		m.root = ms.node
		m.setCurrentNode(m.root.Node)
		m.resetTimer()
		return m.resetInformation()

	case selectedMsg:
		m.rv.SetFocused(m.ViewID())

		data := ms.node.Data[*adNode]()
		if data.nodeType != nodeTypeAction {
			ms.node.SetExpanded(!ms.node.Expanded())
			m.setCurrentNode(ms.node)
			return nil
		}

		m.setCurrentNode(ms.node)
		return acStateToOpMsg(data.ID(), data.actionState).SendRoutedMsg(m.rv)

	case tview.KeyMsg:
		data := m.currentNode.Data[*adNode]()
		if rmsg, ok := data.Noder().HandleKeys(ms); ok {
			return func() tview.Msg {
				return rmsg
			}
		}

	default:
	}

	return nil
}

func (m *Model) populate() tview.Msg {
	root := newRootNode(m.rv)

	adapters, err := m.rv.Session().Adapters()
	if err != nil {
		// TODO: Log error messages
		return nil
	}

	for _, adapter := range adapters {
		devices, err := m.rv.Session().Adapter(adapter.AdapterAddress).Devices()
		if err != nil {
			continue
		}

		root.addAdapter(adapter, devices)
	}

	return ui.ViewIDAdTree.RouterMessage(rootNodeMsg{root})
}

func (m *Model) dbHandler() tview.Cmd {
	currNode := m.currentNode
	if currNode == nil {
		return m.dbTimer
	}

	node := currNode.Data[*adNode]()
	rootID := node.noder.RootID()
	if m.rootident == rootID && !m.reset {
		return m.dbTimer
	}

	m.rootident = rootID
	m.reset = false

	id, inf := node.noder.Information()

	return tview.Batch(func() tview.Msg {
		return info.Msg(id, inf)
	}, m.dbTimer)
}

func (m *Model) resetTimer() {
	m.tm.Reset(dbTime)
}

func (m *Model) resetInformation() tview.Cmd {
	if m.reset {
		return nil
	}

	currNode := m.currentNode
	if currNode == nil {
		return nil
	}

	node := currNode.Data[*adNode]()
	rootID := node.noder.RootID()
	if m.rootident == rootID {
		return nil
	}

	m.reset = true
	return func() tview.Msg {
		return info.ResetMsg()
	}
}

func (m *Model) setCurrentNode(n *ctree.Node) {
	m.currentNode = n
}

func (m *Model) dbTimer() tview.Msg {
	<-m.tm.C
	return m.ViewID().RouterMessage(dbEvent{})
}

type (
	selectionChange *ctree.Node

	scrollMsg struct {
		change viewport.ScrollState
		msg    tview.Msg
	}

	selectedMsg struct {
		node *ctree.Node
	}
)

type (
	rootNodeMsg struct {
		node *rootNode
	}

	dbEvent struct{}
)

var (
	_ ui.Model[Model] = Model{}
	_ ctree.Provider  = (*Model)(nil)
	_ ui.View         = (*Model)(nil)
)
