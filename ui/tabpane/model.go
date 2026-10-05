package tabpane

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// Model holds a tabbed pane.
type Model struct {
	handler *viewHandler

	arrowLeft, arrowRight string
	arrowWidth            int

	rv ui.RootView
}

// New returns a new tabbed pane.
func New(rv ui.RootView) Model {
	leftArrow := theme.Icons().ArrowLeft.String() + " "
	rightArrow := " " + theme.Icons().ArrowRight.String()

	return Model{
		rv: rv,

		handler:    newHandler(rv),
		arrowLeft:  leftArrow,
		arrowRight: rightArrow,
		arrowWidth: uniseg.StringWidth(leftArrow),
	}
}

// ViewID returns the view's ID.
func (m *Model) ViewID() ui.ViewID {
	return ui.ViewIDTabs
}

// HandleRouterMsg handles the routed message.
func (m *Model) HandleRouterMsg(focused bool, routerMsg ui.RouterMsg) tview.Cmd {
	var cmd tview.Cmd

	switch routerMsg.ID {
	case ui.ViewIDTabs:
		*m, cmd = m.Update(focused, routerMsg.Msg)
		return cmd

	default:
	}

	return m.handler.updateRouted(routerMsg.ID, focused, routerMsg.Msg)
}

// LabelStyle provides the styles for the inactive and active labels.
func (m *Model) LabelStyle(active bool) tcell.Style {
	if active {
		return theme.Current().TabsPane.FocusedTab
	}

	return theme.Current().TabsPane.Tab
}

// ArrowConfig provides the left/right arrow icons and the style to be applied.
func (m *Model) ArrowConfig() (arrowLeft string, arrowRight string, arrowStyle tcell.Style, arrowWidth int) {
	return m.arrowLeft, m.arrowRight, theme.Current().TabsPane.Style, m.arrowWidth
}

// Keybind converts a [tview.KeyMsg] to an [Action].
func (m *Model) Keybind(msg tview.KeyMsg) tabs.Action {
	switch {
	case kb().TabPane.Next.Matches(msg):
		return tabs.ActionNext

	case kb().TabPane.Previous.Matches(msg):
		return tabs.ActionPrevious
	}

	return tabs.ActionNone
}

// OnSelect is called when a new tab item is selected.
func (m *Model) OnSelect(index int) tview.Msg {
	return ui.ViewIDTabs.RouterMessage(selectedMsg{index})
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return m.handler.init()
}

// Update returns the model changed in response to a message and a command to run, or nil.
func (m Model) Update(focused bool, msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(focused, msg)
	return m, cmd
}

// View returns the element that draws the model.
func (m Model) View(focused bool) tview.Widget {
	return tabs.New(&m, m.handler.State).
		Focused(focused).
		Content(m.handler.activeContent(focused))
}

func (m *Model) updateModel(focused bool, msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case selectedMsg:
		m.handler.SetActiveIndex(ms.index)
		return nil

	case tview.KeyMsg:
		if !focused {
			return nil
		}

		return m.handler.updateActive(focused, msg)
	}

	return m.handler.updateAll(focused, msg)
}

type selectedMsg struct {
	index int
}

func kb() *keybindings.Keybindings {
	return keybindings.Current
}

var (
	_ ui.Model[Model] = Model{}
	_ ui.View         = (*Model)(nil)
)
