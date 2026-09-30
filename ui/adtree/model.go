package adtree

import (
	"iter"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/gdamore/tcell/v3"
)

// Model represents an adapter-device tree.
type Model struct {
	root *rootNode

	selectionState ctree.SelectionState

	graphicsSet tview.BorderSet
	markers     ctree.Markers

	rv ui.RootView
}

// New returns a new adapter-device tree.
func New(rv ui.RootView) Model {
	root := newRootNode(rv)

	state := ctree.SelectionState{}
	state.SetCurrentNode(root.Node)

	graphicsSet := tview.BorderSetRound()
	markers := ctree.Markers{
		Expanded:  theme.Icons().TriangleDown.String() + " ",
		Collapsed: theme.Icons().TriangleRight.String() + " ",
		Leaf:      " ",
	}

	return Model{
		root: root,

		selectionState: state,

		graphicsSet: graphicsSet,
		markers:     markers,

		rv: rv,
	}
}

// ViewID returns the view's ID.
func (m *Model) ViewID() ui.ViewID {
	return ui.ViewIDAdTree
}

// HandleRouterMsg handles the routed message.
func (m *Model) HandleRouterMsg(msg ui.RouterMsg) tview.Cmd {
	return m.updateModel(msg)
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return m.populate
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(focused bool) tview.Element {
	tree := ctree.New(m.root.Node, &m, &m.selectionState).
		Focused(focused).
		Graphics(true).
		GraphicsSet(m.graphicsSet).
		GraphicsStyle(theme.Current().ADTree.Bg).
		Markers(m.markers)

	return box.New(tree).
		Background(theme.Current().Global.GetBackground()).
		Padding(2, 1, 1, 1)
}

// Keybind converts keybindings to actions.
func (m *Model) Keybind(msg tview.KeyMsg) (ctree.Action, bool) {
	switch {
	case kb().NavigateUp.Matches(msg):
		return ctree.ActionUp, true

	case kb().NavigateDown.Matches(msg):
		return ctree.ActionDown, true

	case kb().ADTree.ToggleNodes.Matches(msg):
		return ctree.ActionSelect, true
	}

	return 0, false
}

// OnChange processes a change in the tree's state and returns a [tview.Msg]
func (m *Model) OnChange(chg ctree.Change) tview.Msg {
	return ui.ViewIDAdTree.RouterMessage(selectionChange(chg))
}

// OnSelect returns a message after processing the selected node.
func (m *Model) OnSelect(node *ctree.Node) tview.Msg {
	return ui.ViewIDAdTree.RouterMessage(selectedMsg{node})
}

// MarkerStyle returns the styles for the markers.
func (m *Model) MarkerStyle(*ctree.Node) tcell.Style {
	return theme.Current().ADTree.Bg
}

// StyledLabels returns a sequence of text and their associated styles.
func (m *Model) StyledLabels(node *ctree.Node, selected bool) iter.Seq2[string, tcell.Style] {
	return func(yield func(string, tcell.Style) bool) {
		data := node.Data[*adNode]()

		style := data.GetStyle()
		if selected {
			style = theme.Current().ADTree.Selection
		}

		for _, segment := range [3]string{data.GetIcon(), " ", data.Name()} {
			if !yield(segment, style) {
				return
			}
		}
	}
}

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case selectionChange:
		m.selectionState.Apply(ctree.Change(ms))

	case rootNodeMsg:
		m.root = ms.node
		m.selectionState.SetCurrentNode(m.root.Node)

	case selectedMsg:
		data := ms.node.Data[*adNode]()
		if data.nodeType != nodeTypeAction {
			ms.node.SetExpanded(!ms.node.Expanded())
			m.selectionState.SetCurrentNode(ms.node)
			return nil
		}

		m.selectionState.SetCurrentNode(ms.node)
		return acStateToOpMsg(data.ID(), data.actionState).SendRoutedMsg(m.rv)

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

type selectionChange ctree.Change

type selectedMsg struct {
	node *ctree.Node
}

type rootNodeMsg struct {
	node *rootNode
}

func kb() *keybindings.Keybindings {
	return keybindings.Current
}

var (
	_ ui.Model[Model] = Model{}
	_ ui.View         = (*Model)(nil)
)
