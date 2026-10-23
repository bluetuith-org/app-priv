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
	"github.com/bluetuith-org/bluetuith/ui/widgets/popup"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

// Model holds and renders a set of views for the application.
type Model struct {
	ui.AppBinder

	headerContent string

	adTreeModel *adtree.Model
	tabsModel   *tabpane.Model

	pm *popup.Manager

	focused ui.ViewID
	ev      *evState
}

// New returns a new root model.
func New(binder ui.AppBinder) *Model {
	m := &Model{
		AppBinder:     binder,
		headerContent: fmt.Sprintf("bluetuith %s (%s)", binder.Configuration().Version, binder.Configuration().Revision),

		pm: popup.NewManager(),
		ev: newEvState(binder.Screen()),
	}

	m.adTreeModel = adtree.New(m)
	m.tabsModel = tabpane.New(m)
	m.focused = ui.ViewIDAdTree

	return m
}

// Init returns a command to run when the model starts, or nil.
func (m *Model) Init() tview.Cmd {
	ui.StartLogHandler()

	return tview.Batch(m.adTreeModel.Init(), m.tabsModel.Init())
}

// Update receives messages when this model has focus.
func (m *Model) Update(msg tview.Msg) (*Model, tview.Cmd) {
	rsp := ui.NewResponse()
	depth := 10

	m.updateModel(msg, rsp)
	for range depth {
		msgs, ok := rsp.Msgs()
		if !ok {
			break
		}

		for _, msg := range msgs {
			m.updateModel(msg, rsp)
		}
	}

	return m, rsp.Cmds()
}

// View draws this model onto the screen.
func (m *Model) View() tview.Widget {
	panelFocused := m.PanelFocused()

	header := stext.New().
		Background(theme.Current().App.TitleBar).
		PushSpace(tview.HorizontalAlignmentLeft).
		Push(m.adTreeModel.Header(), tview.HorizontalAlignmentLeft).
		PushSpace(tview.HorizontalAlignmentLeft)

	for line := range m.tabsModel.Header() {
		header = header.Push(line, tview.HorizontalAlignmentLeft).PushSpace(tview.HorizontalAlignmentLeft)
	}

	header = header.
		Push(stext.NewLine(
			stext.NewSpan(m.headerContent, theme.Current().App.TitleBar),
		),
			tview.HorizontalAlignmentRight,
		).
		PushSpace(tview.HorizontalAlignmentRight)

	cvStyle := theme.Current().App.CurrentViewName
	line := stext.NewLine(
		stext.NewSpan(" ", cvStyle),
		stext.NewSpan(
			m.focused.String(),
			cvStyle,
		),
		stext.NewSpan(" ", cvStyle),
	)

	status := stext.New().
		Background(theme.Current().StatusBar.Style)

	for line := range m.tabsModel.Status() {
		status = status.Push(line, tview.HorizontalAlignmentLeft)
	}

	status = status.
		PushSpace(tview.HorizontalAlignmentLeft).
		Push(line, tview.HorizontalAlignmentRight)

	vflex := row.New(
		m.adTreeModel.View(!panelFocused),
		m.tabsModel.View(panelFocused),
	)

	hflex := column.New(
		header,
		vflex,
		status,
	)

	bx := box.New(hflex).Background(theme.Current().Global.GetBackground())

	if m.pm.HasPopups() {
		return popup.New(bx, m.pm)
	}

	return bx
}

// Focused returns the currently focused view.
func (m *Model) Focused() ui.ViewID {
	return m.focused
}

// SetFocused stores the currently focused view.
func (m *Model) SetFocused(id ui.ViewID) {
	m.focused = id
}

// PanelFocused returns whether the panel is focused.
func (m *Model) PanelFocused() bool {
	return m.Focused() != ui.ViewIDAdTree
}

// SendMsg sends a message to the program.
func (m *Model) SendMsg(msg tcell.Event) {
	m.ev.sendMsg(msg)
}

func (m *Model) updateModel(msg tview.Msg, rsp *ui.Response) {
	switch ms := msg.(type) {
	case tview.KeyMsg:
		if m.pm.HandleKey(ms, rsp) {
			return
		}

		if m.tabsModel.SetViewKeybind(ms) {
			return
		}

		switch {
		case ui.Kb().Quit.Matches(ms):
			rsp.PushCmd(m.quitCmd())

		case ui.Kb().ADTree.TabPane.Matches(ms):
			if m.Focused() == ui.ViewIDAdTree {
				m.focused = m.tabsModel.GetActiveID()
			}

		case ui.Kb().TabPane.ADTree.Matches(ms):
			if m.Focused() != ui.ViewIDAdTree {
				m.focused = ui.ViewIDAdTree
			}

		default:
			m.routeMessageToView(m.Focused(), ms, rsp)
		}

	case popup.Router:
		m.pm.HandleMsg(ms, rsp)

	case ui.Router:
		if ms == nil || !ms.IsValid() {
			return
		}

		m.routeMessageToView(ms.ViewID(), ms, rsp)

	default:
		m.updateAllViews(msg, rsp)
	}
}

func (m *Model) updateAllViews(msg tview.Msg, rsp *ui.Response) {
	m.adTreeModel.Update(msg, rsp)

	m.tabsModel.Update(msg, rsp)
}

func (m *Model) routeMessageToView(id ui.ViewID, routedMsg tview.Msg, rsp *ui.Response) {
	if id == ui.ViewIDAdTree {
		m.adTreeModel.Update(routedMsg, rsp)
		return
	}

	m.tabsModel.HandleRouterMsg(id, routedMsg, rsp)
}

func (m *Model) quitCmd() tview.Cmd {
	m.ev.close()

	return tview.Quit()
}

var (
	_ ui.RootView         = (*Model)(nil)
	_ tview.Model[*Model] = (*Model)(nil)
)
