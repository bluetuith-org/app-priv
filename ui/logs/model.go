package logs

import (
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/richtext"
	"github.com/ayn2op/tview/scrollbar"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
)

//revive:disable

type Model struct {
	h *handler

	cursor, target int
	lastIndex      int

	scrollState viewport.ScrollState
	scrollBar   scrollbar.Widget

	follow bool
}

func New(rv ui.RootView) Model {
	h := newHandler(rv, 1000)
	ui.GetLogger().AddHandler(h)

	return Model{
		cursor:      -1,
		target:      -1,
		h:           h,
		scrollState: viewport.ScrollState{},
		scrollBar:   scrollbar.New(),
		follow:      true,
	}
}

// ViewID returns the view's ID.
func (m Model) ViewID() ui.ViewID {
	return ui.ViewIDLog
}

// Title returns the title of the tab section.
func (m Model) Title() string {
	return fmt.Sprintf("Log (%s)", ui.Kb().Log.Show.KeyName)
}

// Icon returns the icon associated with the tab section.
func (m Model) Icon() string {
	return theme.Icons().Log.String()
}

// StatusBar returns the content for the status bar.
func (m *Model) StatusBar() richtext.Line {
	return m.h.getLatest()
}

// ScrollKeybinds implements keybindings for the viewport.
func (m Model) ScrollKeybinds(msg tview.KeyMsg) viewport.Action {
	switch {
	case ui.Kb().NavigateTop.Matches(msg):
		return viewport.ActionTop

	case ui.Kb().NavigateBottom.Matches(msg):
		return viewport.ActionBottom
	}

	return viewport.ActionNone
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m Model) Update(focused bool, msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(focused, msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(focused bool) tview.Widget {
	ls := newWidget(m.h, m.cursor, focused && !m.follow)

	vp := viewport.New(ls, m.scrollState).
		Target(ls.Target(m.target)).
		ScrollBar(m.scrollBar, viewport.ScrollBarVisibilityAutomatic).
		OnChildMsg(func(c *viewport.ScrollState, msg tview.Msg) tview.Msg { return scrollMsg{change: *c, msg: msg} }).
		Axes(layout.Axes{Width: true, Height: true}).
		Keybind(m.ScrollKeybinds).
		Focused(focused).
		OnChange(func(c *viewport.ScrollState) tview.Msg { return scrollMsg{change: *c} })

	bx := box.New(vp).
		TitleStyle(theme.Current().Log.FollowIndicator).
		TitleAlignment(tview.AlignmentRight).
		Background(theme.Current().Global.GetBackground())

	if m.follow {
		bx = bx.Title("- Following -")
	}

	return bx
}

func (m *Model) setCursor(index int) {
	if index >= 0 && index != m.cursor {
		m.scrollState.ScrollToTarget()
		m.target = index
	}

	m.cursor = max(index, -1)
}

func (m *Model) setFollow() {
	if m.follow {
		m.setCursor(m.lastIndex)
		m.scrollState.ScrollToEnd()
	}
}

func (m *Model) updateModel(focused bool, msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case newLogMsg:
		m.lastIndex = ms.index
		if !focused {
			m.setFollow()
		}

		return nil

	case listChg:
		m.setCursor(int(ms))
		return nil

	case scrollMsg:
		m.scrollState = ms.change
		if ms.msg != nil {
			return m.updateModel(focused, ms.msg)
		}

		return nil

	case tview.KeyMsg:
		if ui.Kb().Log.Follow.Matches(ms) {
			m.follow = !m.follow
			m.setFollow()
		}

	case ui.FocusMsg:
		if ms.ID != m.ViewID() {
			m.setFollow()
		}
	}

	return nil
}

type (
	listChg int

	scrollMsg struct {
		change viewport.ScrollState
		msg    tview.Msg
	}
)

var (
	_ ui.Model[Model] = Model{}
	_ tabs.TabSection = Model{}
	_ ui.View         = (*Model)(nil)
)
