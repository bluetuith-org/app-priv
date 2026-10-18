package root

import (
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/column"
	"github.com/ayn2op/tview/row"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/adtree"
	"github.com/bluetuith-org/bluetuith/ui/tabpane"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

// Model holds and renders a set of views for the application.
type Model struct {
	ui.AppBinder

	headerContent string

	adTreeModel adtree.Model
	tabsModel   tabpane.Model

	state *focusState
	ev    *evState
}

// New returns a new root model.
func New(binder ui.AppBinder) Model {
	m := Model{
		AppBinder:     binder,
		headerContent: fmt.Sprintf("bluetuith %s (%s)", binder.Configuration().Version, binder.Configuration().Revision),

		ev: newEvState(binder.Screen()),
	}

	m.adTreeModel = adtree.New(&m)
	m.tabsModel = tabpane.New(&m)
	m.state = newFocusState(ui.ViewIDAdTree)

	return m
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	ui.StartLogHandler()

	return tview.Batch(m.adTreeModel.Init(), m.tabsModel.Init())
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View() tview.Widget {
	panelFocused := m.PanelFocused()

	header := stext.New().
		Background(theme.Current().App.TitleBar).
		PushSpace(tview.AlignmentLeft).
		Push(m.adTreeModel.Header(), tview.AlignmentLeft).
		PushSpace(tview.AlignmentLeft)

	for _, h := range m.tabsModel.Header() {
		header = header.Push(h, tview.AlignmentLeft).PushSpace(tview.AlignmentLeft)
	}

	header = header.
		Push(stext.NewLine(
			stext.NewSpan(m.headerContent, theme.Current().App.TitleBar),
		),
			tview.AlignmentRight,
		).
		PushSpace(tview.AlignmentRight)

	cvStyle := theme.Current().App.CurrentViewName
	line := stext.NewLine(
		stext.NewSpan(" ", cvStyle),
		stext.NewSpan(
			m.state.Focused().String(),
			cvStyle,
		),
		stext.NewSpan(" ", cvStyle),
	)

	status := stext.New().
		Background(theme.Current().StatusBar.Style).
		Push(m.tabsModel.Status(), tview.AlignmentLeft).
		PushSpace(tview.AlignmentLeft).
		Push(line, tview.AlignmentRight)

	vflex := row.New(
		m.adTreeModel.View(!panelFocused),
		m.tabsModel.View(panelFocused),
	)

	hflex := column.New(
		header,
		vflex,
		status,
	)

	return box.New(hflex).Background(theme.Current().Global.GetBackground())
}

// Focused returns the currently focused view.
func (m *Model) Focused() ui.ViewID {
	return m.state.Focused()
}

// SetFocused stores the currently focused view.
func (m *Model) SetFocused(id ui.ViewID) {
	m.state.SetFocus(id)
}

// PanelFocused returns whether the panel is focused.
func (m *Model) PanelFocused() bool {
	return m.Focused() != ui.ViewIDAdTree
}

// SendMsg sends a message to the program.
func (m *Model) SendMsg(msg tcell.Event) {
	m.ev.sendMsg(msg)
}

// SendRoutedUpdateMsg routes the message to the specified view.
// Should only be called from the view's [Update] function.
func (m *Model) SendRoutedUpdateMsg(msg ui.Router) tview.Cmd {
	if msg == nil || !msg.IsValid() {
		return nil
	}

	return m.updateModel(msg)
}

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case tview.KeyMsg:
		if m.tabsModel.SetViewKeybind(ms) {
			return nil
		}

		switch {
		case ui.Kb().Quit.Matches(ms):
			return m.quitCmd()

		case ui.Kb().ADTree.TabPane.Matches(ms):
			if m.Focused() == ui.ViewIDAdTree {
				m.state.SetFocus(m.tabsModel.GetActiveID())
			}

		case ui.Kb().TabPane.ADTree.Matches(ms):
			if m.Focused() != ui.ViewIDAdTree {
				m.state.SetFocus(ui.ViewIDAdTree)
			}
		}

		return m.routeMessageToView(m.Focused(), ms)

	case ui.Router:
		if ms == nil || !ms.IsValid() {
			return nil
		}

		return m.routeMessageToView(ms.ViewID(), ms)
	}

	return m.updateAllViews(msg)
}

func (m *Model) updateAllViews(msg tview.Msg) tview.Cmd {
	cmds := make([]tview.Cmd, 0, 2)
	var cmd tview.Cmd

	m.adTreeModel, cmd = m.adTreeModel.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.tabsModel, cmd = m.tabsModel.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tview.Batch(cmds...)
}

func (m *Model) routeMessageToView(id ui.ViewID, routedMsg tview.Msg) tview.Cmd {
	var cmd tview.Cmd

	if id == ui.ViewIDAdTree {
		m.adTreeModel, cmd = m.adTreeModel.Update(routedMsg)
		return cmd
	}

	return m.tabsModel.HandleRouterMsg(id, routedMsg)
}

func (m *Model) quitCmd() tview.Cmd {
	m.ev.close()

	return tview.Quit()
}

var (
	_ ui.RootView        = (&Model{})
	_ tview.Model[Model] = (Model{})
)
