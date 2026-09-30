package tabpane

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/info"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// Model holds a tabbed pane.
type Model struct {
	state *tabs.State

	arrowLeft, arrowRight string
	arrowWidth            int

	infoModel info.Model

	rv ui.RootView
}

// New returns a new tabbed pane.
//
// NOTE: To add a new view within this model:
// NOTE: - First, add the view to the Model struct, and the tab state.
// NOTE: - Then, add the view's Init() method to [Model.Init].
// NOTE: - After that, add a handler for the view within [Model.HandleRouterMsg]
// NOTE:   and [Model.View].
func New(rv ui.RootView) Model {
	leftArrow := theme.Icons().ArrowLeft.String() + " "
	rightArrow := " " + theme.Icons().ArrowRight.String()

	infoModel := info.New()

	state := tabs.NewState()
	state.AddSection(infoModel.Icon(), infoModel.Title())
	state.AddSection(infoModel.Icon(), infoModel.Title())

	return Model{
		rv: rv,

		infoModel: infoModel,

		state:      state,
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
func (m *Model) HandleRouterMsg(routerMsg ui.RouterMsg) tview.Cmd {
	var cmd tview.Cmd

	switch routerMsg.ID {
	case ui.ViewIDTabs:
		*m, cmd = m.Update(routerMsg.Msg)

	case ui.ViewIDInfo:
		m.infoModel, cmd = m.infoModel.Update(routerMsg.Msg)

	default:
	}

	return cmd
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
func (m *Model) Keybind(msg tview.KeyMsg) (tabs.Action, bool) {
	switch {
	case kb().SwitchPanes.Matches(msg):
		return tabs.ActionNext, true

	default:
	}

	return 0, false
}

// OnSelect is called when a new tab item is selected.
func (m *Model) OnSelect(index int) tview.Msg {
	return ui.ViewIDTabs.RouterMessage(selectedMsg{index})
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return tview.Batch(m.infoModel.Init())
}

// Update returns the model changed in response to a message and a command to run, or nil.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View returns the element that draws the model.
func (m Model) View(focused bool) tview.Element {
	var content tview.Element

	switch m.state.ActiveIndex() {
	case 0:
		content = m.infoModel.View(focused)

	case 1:
		content = m.infoModel.View(focused)
	}

	return tabs.New(&m, m.state).Focused(focused).Content(content)
}

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case selectedMsg:
		m.state.SetActiveIndex(ms.index)

	default:
	}

	return nil
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
