package root

import (
	"context"
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/column"
	"github.com/ayn2op/tview/row"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/tabpane"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
)

// Model holds and renders a set of views for the application.
type Model struct {
	ui.AppBinder

	header stext.Widget
	layout box.Widget

	initedViews map[ui.ViewID]ui.View

	tabsModel tabpane.Model

	ctx     context.Context
	cancel  context.CancelFunc
	msgChan chan tview.Msg
}

// New returns a new root model.
func New(binder ui.AppBinder) Model {
	header := fmt.Sprintf(" bluetuith %s (%s)", binder.Configuration().Version, binder.Configuration().Revision)
	ctx, cancel := context.WithCancel(context.Background())

	m := Model{
		AppBinder:   binder,
		header:      stext.New(tview.AlignmentRight).SetContent(header),
		initedViews: make(map[ui.ViewID]ui.View),

		ctx:     ctx,
		cancel:  cancel,
		msgChan: make(chan tview.Msg, 1),
	}

	m.tabsModel = tabpane.New(&m)

	return m
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return m.listenForMsg()
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	vc := m.updateModel(msg)
	return m, vc
}

// View draws this model onto the screen.
func (m Model) View() tview.Element {
	vflex := row.New(m.tabsModel.View(true))
	hflex := column.New(m.header.SetStyle(theme.Current().TitleBar), vflex, m.header.SetStyle(theme.Current().TitleBar))

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
		case kb().Quit.Matches(ms):
			return m.quitCmd()

		default:
		}

	default:
	}

	return nil
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

func kb() *keybindings.Keybindings {
	return keybindings.Current
}

var (
	_ ui.RootView        = (&Model{})
	_ tview.Model[Model] = (Model{})
)
