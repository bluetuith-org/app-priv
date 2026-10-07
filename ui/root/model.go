package root

import (
	"context"
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
)

// Model holds and renders a set of views for the application.
type Model struct {
	ui.AppBinder

	header stext.Widget

	adTreeModel adtree.Model
	tabsModel   tabpane.Model

	ctx     context.Context
	cancel  context.CancelFunc
	msgChan chan tview.Msg

	panelFocused bool
}

// New returns a new root model.
func New(binder ui.AppBinder) Model {
	header := fmt.Sprintf(" bluetuith %s (%s)", binder.Configuration().Version, binder.Configuration().Revision)
	ctx, cancel := context.WithCancel(context.Background())

	m := Model{
		AppBinder: binder,
		header:    stext.New(tview.AlignmentRight).SetContent(header),

		ctx:     ctx,
		cancel:  cancel,
		msgChan: make(chan tview.Msg),
	}

	m.adTreeModel = adtree.New(&m)
	m.tabsModel = tabpane.New(&m)

	return m
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return tview.Batch(m.listenForMsg(), m.adTreeModel.Init(), m.tabsModel.Init())
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View() tview.Widget {
	vflex := row.New(
		m.adTreeModel.View(!m.panelFocused),
		m.tabsModel.View(m.panelFocused),
	)
	hflex := column.New(
		m.header.SetStyle(theme.Current().App.TitleBar),
		vflex,
		m.header.SetStyle(theme.Current().App.TitleBar),
	)

	return box.New(hflex).Background(theme.Current().Global.GetBackground())
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

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case tview.KeyMsg:
		switch {
		case ui.Kb().Quit.Matches(ms):
			return m.quitCmd()

		case ui.Kb().ADTree.TabPane.Matches(ms):
			if !m.panelFocused {
				m.panelFocused = true
			}

			return nil

		case ui.Kb().TabPane.ADTree.Matches(ms):
			if m.panelFocused {
				m.panelFocused = false
			}

			return nil
		}

	case ui.FocusMsg:
		switch ms.ID {
		case ui.ViewIDAdTree:
			if m.panelFocused {
				m.panelFocused = false
			}

		case ui.ViewIDTabs:
			if !m.panelFocused {
				m.panelFocused = true
			}
		}

		return nil

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

	m.adTreeModel, cmd = m.adTreeModel.Update(!m.panelFocused, msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	m.tabsModel, cmd = m.tabsModel.Update(m.panelFocused, msg)
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

	switch routerMsg.ID {
	case ui.ViewIDAdTree:
		m.adTreeModel, cmd = m.adTreeModel.Update(!m.panelFocused, routerMsg.Msg)
		return cmd

	default:
	}

	return m.tabsModel.HandleRouterMsg(m.panelFocused, routerMsg)
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
