package logs

import (
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/scrollbar"
	"github.com/ayn2op/tview/text"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
)

// Model represents a log view.
type Model struct {
	h *handler

	cursor, target int
	lastIndex      int

	scrollState viewport.ScrollState
	scrollBar   scrollbar.Widget

	follow bool
}

// New returns a new log view.
func New(rv ui.RootView) Model {
	h := newHandler(rv, 1000)
	ui.AddLogHandler(h)

	return Model{
		cursor:    0,
		target:    0,
		h:         h,
		scrollBar: scrollbar.New(),
		follow:    true,
	}
}

// ViewID returns the view's ID.
func (m Model) ViewID() ui.ViewID {
	return ui.ViewIDLog
}

// Title returns the title of the tab section.
func (m Model) Title() string {
	return fmt.Sprintf("%s (%s)", m.ViewID().String(), ui.Kb().Log.Show.KeyName)
}

// Icon returns the icon associated with the tab section.
func (m Model) Icon() string {
	return theme.Icons().Log.String()
}

// StatusBar returns the content for the status bar.
func (m *Model) StatusBar() []text.Span {
	return m.h.getLatest()
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

// OnListChange processes a change in the list and returns a message.
func (m Model) OnListChange(cursor int) tview.Msg {
	return m.ViewID().RouterMessage(listChg(cursor))
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
	return nil
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(focused bool) tview.Widget {
	ls := newWidget(m.h, m.cursor, focused && !m.follow)

	vp := viewport.New(ls, m.scrollState).
		Target(ls.Target(m.target)).
		ScrollBar(m.scrollBar, viewport.ScrollBarVisibilityAutomatic).
		Axes(layout.Axes{Width: true, Height: true}).
		Focused(focused).
		Keybind(m.ViewportKeybinds).
		OnChange(m.OnViewportChange).
		OnChildMsg(m.OnViewportChildMsg)

	bx := box.New(vp).
		TitleStyle(theme.Current().Log.FollowIndicator).
		TitleAlignment(tview.AlignmentRight).
		Background(theme.Current().Global.GetBackground())

	if m.follow {
		bx = bx.Title("- Following -")
	}

	return bx
}

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case ui.R[newLogMsg]:
		m.lastIndex = ms.Value.index
		m.setFollow()

	case ui.R[listChg]:
		m.h.rv.SetFocused(m.ViewID())
		if !m.follow {
			m.setCursor(int(ms.Value))
		}

	case ui.R[scrollMsg]:
		if !m.follow {
			m.scrollState = ms.Value.change
		}
		if ms.Value.msg != nil {
			return m.updateModel(ms.Value.msg)
		}

	case tview.KeyMsg:
		if ui.Kb().Log.Follow.Matches(ms) {
			m.follow = !m.follow
			m.setFollow()
		}
	}

	return nil
}

func (m *Model) setCursor(index int) {
	if index < 0 {
		return
	}

	m.scrollState.ScrollToTarget()
	m.target = index
	m.cursor = max(index, -1)
}

func (m *Model) setFollow() {
	if m.follow {
		m.setCursor(m.lastIndex)
	}
}

type (
	listChg int

	scrollMsg struct {
		change viewport.ScrollState
		msg    tview.Msg
	}

	newLogMsg struct {
		index int
	}
)

var (
	_ ui.Model[Model] = Model{}
	_ ui.TabSection   = Model{}
	_ ui.View         = (*Model)(nil)
)
