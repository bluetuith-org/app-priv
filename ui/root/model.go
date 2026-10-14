package root

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/column"
	"github.com/ayn2op/tview/richtext"
	"github.com/ayn2op/tview/row"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/adtree"
	"github.com/bluetuith-org/bluetuith/ui/tabpane"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
)

// Model holds and renders a set of views for the application.
type Model struct {
	ui.AppBinder

	headerContent string

	adTreeModel adtree.Model
	tabsModel   tabpane.Model

	ctx     context.Context
	cancel  context.CancelFunc
	msgChan chan tview.Msg

	state *focusState
}

// New returns a new root model.
func New(binder ui.AppBinder) Model {
	ctx, cancel := context.WithCancel(context.Background())

	m := Model{
		AppBinder:     binder,
		headerContent: fmt.Sprintf("bluetuith %s (%s)", binder.Configuration().Version, binder.Configuration().Revision),

		ctx:     ctx,
		cancel:  cancel,
		msgChan: make(chan tview.Msg),
	}

	m.adTreeModel = adtree.New(&m)
	m.tabsModel = tabpane.New(&m)
	m.state = newFocusState(ui.ViewIDAdTree)

	return m
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	ui.GetLogger().Start()

	return tview.Batch(m.listenForMsg(), m.adTreeModel.Init(), m.tabsModel.Init())
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View() tview.Widget {
	panelFocused := m.state.Focused() != ui.ViewIDAdTree

	header := stext.New().
		Background(theme.Current().App.TitleBar).
		PushSpace(tview.AlignmentLeft).
		Push(m.adTreeModel.Header(), tview.AlignmentLeft).
		PushSpace(tview.AlignmentLeft)

	for _, h := range m.tabsModel.Header() {
		header = header.Push(h, tview.AlignmentLeft).PushSpace(tview.AlignmentLeft)
	}

	header = header.
		Push(richtext.NewLine(
			richtext.NewSegment(m.headerContent, theme.Current().App.TitleBar),
		),
			tview.AlignmentRight,
		).
		PushSpace(tview.AlignmentRight)

	status := stext.New().
		Background(theme.Current().StatusBar.Style).
		Push(m.tabsModel.Status(), tview.AlignmentLeft)

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

// SendMsg sends a message to the program.
func (m *Model) SendMsg(msg tview.Msg) {
	select {
	case <-m.ctx.Done():
	case m.msgChan <- msg:
	}
}

// SendRoutedUpdateMsg routes the message to the specified view.
// Should only be called from the view's [Update] function.
func (m *Model) SendRoutedUpdateMsg(msg ui.RouterMsg) tview.Cmd {
	if !msg.IsValid() {
		return nil
	}

	_, cmd := m.Update(msg)

	return cmd
}

var i = 0

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case tview.KeyMsg:
		switch {
		case ui.Kb().Quit.Matches(ms):
			return m.quitCmd()

		case ui.Kb().Filebrowser.ToggleHiddenFiles.Matches(ms):
			slog.Info("Hello "+strconv.Itoa(i), "message", i)
			i++
			return nil

		case ui.Kb().ADTree.TabPane.Matches(ms):
			m.state.SetFocus(m.tabsModel.GetActiveID())

		case ui.Kb().TabPane.ADTree.Matches(ms):
			m.state.SetFocus(ui.ViewIDAdTree)
		}

	case externalMsg:
		return tview.Batch(m.updateModel(ms.msg), m.listenForMsg())

	case ui.RouterMsg:
		if !ms.IsValid() {
			return nil
		}

		return m.routeMessageToView(ms)
	}

	var cmds []tview.Cmd
	var cmd tview.Cmd

	m.adTreeModel, cmd = m.adTreeModel.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.tabsModel, cmd = m.tabsModel.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if cmds == nil {
		return nil
	}

	if len(cmds) == 1 {
		return cmds[0]
	}

	return tview.Batch(cmds...)
}

func (m *Model) routeMessageToView(routerMsg ui.RouterMsg) tview.Cmd {
	var cmd tview.Cmd

	if routerMsg.ID == ui.ViewIDAdTree {
		m.adTreeModel, cmd = m.adTreeModel.Update(routerMsg.Msg)
		return cmd
	}

	return m.tabsModel.HandleRouterMsg(routerMsg)
}

type externalMsg struct {
	msg tview.Msg
}

func (m *Model) listenForMsg() tview.Cmd {
	return func() tview.Msg {
		return externalMsg{<-m.msgChan}
	}
}

func (m *Model) quitCmd() tview.Cmd {
	return tview.Batch(func() tview.Msg {
		m.cancel()
		return nil
	}, tview.Quit())
}

var (
	_ ui.RootView        = (&Model{})
	_ tview.Model[Model] = (Model{})
)
