package adtree

import (
	"iter"
	"log/slog"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/text"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/popup"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/gdamore/tcell/v3"
)

// Model represents an adapter-device tree.
type Model struct {
	root        *rootNode
	currentNode *ctree.Node

	scrollState viewport.ScrollState

	graphicsSet tview.BorderSet
	markers     ctree.Markers

	rv ui.RootView
}

// ViewID returns the view's ID.
func (m *Model) ViewID() ui.ViewID {
	return ui.ViewIDAdTree
}

// New returns a new adapter-device tree.
func New(rv ui.RootView) *Model {
	root := newRootNode(rv)

	graphicsSet := tview.BorderSetRound()
	markers := ctree.Markers{
		Expanded:  theme.Icons().TriangleDown.String() + " ",
		Collapsed: theme.Icons().TriangleRight.String() + " ",
		Leaf:      " ",
	}

	m := &Model{
		root:        root,
		currentNode: root.Node,

		graphicsSet: graphicsSet,
		markers:     markers,

		rv: rv,
	}

	return m
}

// Init returns a command to run when the model starts, or nil.
func (m *Model) Init() tview.Cmd {
	return m.populate
}

// Update receives messages when this model has focus.
func (m *Model) Update(msg tview.Msg, rsp *ui.Response) {
	if m.handleEvents(msg, rsp) {
		return
	}

	switch ms := msg.(type) {
	case ui.R[scrollMsg]:
		m.scrollState = viewport.ScrollState(ms.Value.change)
		if ms.Value.msg != nil {
			m.Update(ms.Value.msg, rsp)
		}

	case ui.R[selectionChange]:
		m.currentNode = ms.Value
		m.scrollState.ScrollToTarget()

	case ui.R[rootNodeMsg]:
		m.root = ms.Value.node
		if ch := m.root.Children(); len(ch) > 0 {
			m.setCurrentNode(ch[0])
		}

	case ui.R[selectedMsg]:
		m.rv.SetFocused(m.ViewID())
		m.setCurrentNode(ms.Value.node)

		acnode, ok := ncheck[*acNode](ms.Value.node)
		if !ok {
			ms.Value.node.SetExpanded(!ms.Value.node.Expanded())
			return
		}

		if actionMsg, ok := acStateToOpMsg(acnode); ok {
			rsp.PushRouted(actionMsg)
		}

	case tview.KeyMsg:
		node := m.currentNode
		if node == nil {
			return
		}

		noder := nget(node)

		if ui.Kb().ADTree.Info.Matches(ms) {
			id := noder.AdNode().ID()
			info := noder.Information()
			if info == nil {
				return
			}

			pop := popup.NewTextView(id, info)
			rsp.PushMsg(popup.CreateMsg(pop))

			return
		}

		if rmsg, ok := noder.HandleKeys(ms); ok {
			rsp.PushRouted(rmsg)
		}
	}
}

// View draws this model onto the screen.
func (m *Model) View(focused bool) tview.Widget {
	state := m.currentNode
	if !focused {
		state = nil
	}

	tree := ctree.New(m, m.root.Node, state).
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
		Keybind(m.ViewportKeybinds).
		OnChange(m.OnViewportChange).
		OnChildMsg(m.OnViewportChildMsg)

	return box.New(vp).
		Background(theme.Current().Global.GetBackground()).
		Padding(2, 0, 2, 1)
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

// OnChange processes a change in the tree's state and returns a [tview.Msg]
func (m *Model) OnChange(chg *ctree.Node) tview.Msg {
	return m.ViewID().R(selectionChange(chg))
}

// OnSelect returns a message after processing the selected node.
func (m *Model) OnSelect(node *ctree.Node) tview.Msg {
	return m.ViewID().R(selectedMsg{node})
}

// ViewportKeybinds implements keybindings for the viewport.
func (m *Model) ViewportKeybinds(msg tview.KeyMsg) viewport.Action {
	switch {
	case ui.Kb().ScrollTop.Matches(msg):
		return viewport.ActionTop

	case ui.Kb().ScrollBottom.Matches(msg):
		return viewport.ActionBottom

	case ui.Kb().ScrollLeft.Matches(msg):
		return viewport.ActionLeft

	case ui.Kb().ScrollRight.Matches(msg):
		return viewport.ActionRight
	}

	return viewport.ActionNone
}

// OnViewportChange processes a change in the viewport and returns a message.
func (m *Model) OnViewportChange(c *viewport.ScrollState) tview.Msg {
	return m.ViewID().R(scrollMsg{change: *c})
}

// OnViewportChildMsg processes a message that was passed to its underlying view in the viewport and returns a message.
func (m *Model) OnViewportChildMsg(c *viewport.ScrollState, msg tview.Msg) tview.Msg {
	return m.ViewID().R(scrollMsg{change: *c, msg: msg})
}

// MarkerStyle returns the styles for the markers.
func (m *Model) MarkerStyle(*ctree.Node) tcell.Style {
	return theme.Current().Global
}

// LabelWidth implements [ctree.Provider].
func (m *Model) LabelWidth(node *ctree.Node) int {
	adn := nget(node).AdNode()

	return adn.labelWidth()
}

// Spacing implements [ctree.Provider].
func (m *Model) Spacing(node *ctree.Node) int {
	if _, ok := ncheck[*acNode](node); ok {
		return 0
	}

	return 1
}

// StyledLabels returns a sequence of text and their associated styles.
func (m *Model) StyledLabels(node *ctree.Node, selected bool) iter.Seq2[text.Span, tcell.Style] {
	return func(yield func(text.Span, tcell.Style) bool) {
		adn := nget(node).AdNode()

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

func (m *Model) populate() tview.Msg {
	defer func() {
		go m.adapterEvents()
		go m.deviceEvents()
	}()

	root := newRootNode(m.rv)

	adapters, err := m.rv.Session().Adapters()
	if err != nil {
		slog.Error("An error occurred while fetching the list of adapters", "internal error", err)
		return nil
	}

	for _, adapter := range adapters {
		root.addAdapter(adapter)
	}

	return ui.ViewIDAdTree.R(rootNodeMsg{root})
}

func (m *Model) setCurrentNode(n *ctree.Node) {
	m.currentNode = n
}

func (m *Model) handleEvents(msg tview.Msg, _ *ui.Response) bool {
	switch ms := msg.(type) {
	case ui.R[adMsg[bluetooth.AdapterData]]:
		m.root.addAdapter(ms.Value.data)

	case ui.R[adMsg[bluetooth.AdapterEventData]]:
		m.root.updateAdapter(ms.Value.data, ms.Value.action == bluetooth.EventActionRemoved)

	case ui.R[adMsg[bluetooth.DeviceData]]:
		m.root.addDevice(ms.Value.data)

	case ui.R[adMsg[bluetooth.DeviceEventData]]:
		m.root.updateDevice(ms.Value.data, ms.Value.action == bluetooth.EventActionRemoved)

	default:
		return false
	}

	return true
}

func (m *Model) adapterEvents() {
	sub, ok := bluetooth.AdapterEvents().Subscribe()
	if !ok {
		slog.Error("Tree: cannot subscribe to adapter events: unknown error")
		return
	}

	slog.Info("Listening for adapter events...")

	for {
		select {
		case <-sub.Done:
			return

		case added := <-sub.AddedEvents:
			m.rv.SendMsg(newAdMsg(added, bluetooth.EventActionAdded))

		case updated := <-sub.UpdatedEvents:
			m.rv.SendMsg(newAdMsg(updated, bluetooth.EventActionUpdated))

		case removed := <-sub.RemovedEvents:
			m.rv.SendMsg(newAdMsg(removed, bluetooth.EventActionRemoved))
		}
	}
}

func (m *Model) deviceEvents() {
	sub, ok := bluetooth.DeviceEvents().Subscribe()
	if !ok {
		slog.Error("Tree: cannot subscribe to device events: unknown error")
		return
	}

	slog.Info("Listening for device events...")

	for {
		select {
		case <-sub.Done:
			return

		case added := <-sub.AddedEvents:
			m.rv.SendMsg(newAdMsg(added, bluetooth.EventActionAdded))

		case updated := <-sub.UpdatedEvents:
			m.rv.SendMsg(newAdMsg(updated, bluetooth.EventActionUpdated))

		case removed := <-sub.RemovedEvents:
			m.rv.SendMsg(newAdMsg(removed, bluetooth.EventActionRemoved))
		}
	}
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
)

type adEvent interface {
	bluetooth.AdapterData | bluetooth.AdapterEventData |
		bluetooth.DeviceData | bluetooth.DeviceEventData
}

type adMsg[T adEvent] struct {
	action bluetooth.EventAction
	data   T
}

func newAdMsg[T adEvent](data T, action bluetooth.EventAction) ui.Router {
	return ui.ViewIDAdTree.R(adMsg[T]{action, data})
}

var (
	_ ui.Model       = (*Model)(nil)
	_ ctree.Provider = (*Model)(nil)
)
