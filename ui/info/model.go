package info

import (
	"fmt"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/scrollbar"
	"github.com/ayn2op/tview/text"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
)

// Model represents an information section.
type Model struct {
	c *cached

	scrollState viewport.ScrollState
	scrollBar   scrollbar.Widget

	placeholder []text.Span

	rv ui.RootView
}

// ViewID returns the view's ID.
func (m *Model) ViewID() ui.ViewID {
	return ui.ViewIDInfo
}

// New returns a new information section.
func New(rv ui.RootView) *Model {
	return &Model{
		rv:        rv,
		c:         newCached(nil),
		scrollBar: scrollbar.New(),
		placeholder: []text.Span{
			stext.NewSpan(
				"Settle on an adapter or device to view information.",
				theme.Current().App.PlaceHolder,
			),
		},
	}
}

// Init returns a command to run when the model starts, or nil.
func (m *Model) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m *Model) Update(msg tview.Msg, rsp *ui.Response) {
	switch ms := msg.(type) {
	case ui.R[scrollMsg]:
		m.scrollState = ms.Value.change
		if ms.Value.msg != nil {
			m.Update(ms.Value.msg, rsp)
		}

	case ui.R[infoMsg]:
		m.c.reset(ms.Value.content)

	case ui.R[resetMsg]:
		m.c.reset(nil)

	default:
	}
}

// View draws this model onto the screen.
func (m *Model) View(focused bool) tview.Widget {
	txt := m.placeholder
	if m.c.hasContent() {
		txt = m.c.getContent()
	}

	t := text.Rich(txt...).
		Style(theme.Current().Global)

	vp := viewport.New(t, m.scrollState).
		ScrollBar(m.scrollBar, viewport.ScrollBarVisibilityAutomatic).
		Axes(layout.Axes{Width: true, Height: true}).
		Focused(focused).
		Keybind(m.ViewportKeybinds).
		OnChange(m.OnViewportChange).
		OnChildMsg(m.OnViewportChildMsg)

	return vp
}

// Title returns the title of the tab section.
func (m *Model) Title() string {
	return fmt.Sprintf("%s (%s)", m.ViewID().String(), ui.Kb().Information.Show.KeyName)
}

// Icon returns the icon associated with the tab section.
func (m *Model) Icon() string {
	return theme.Icons().Info.String()
}

// Header returns the content to print in the title bar area.
func (m *Model) Header() []text.Span {
	return nil
}

// Status returns the content to print in the status bar area.
func (m *Model) Status() []text.Span {
	return nil
}

// ViewKeybind returns the keybinding used to focus the view.
func (m *Model) ViewKeybind() keybindings.Keybinding {
	return ui.Kb().Information.Show
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

	case ui.Kb().NavigateUp.Matches(msg):
		return viewport.ActionUp

	case ui.Kb().NavigateDown.Matches(msg):
		return viewport.ActionDown
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

type scrollMsg struct {
	change viewport.ScrollState
	msg    tview.Msg
}

var _ ui.TabSection = (*Model)(nil)
