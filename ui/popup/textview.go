package popup

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/scrollbar"
	"github.com/ayn2op/tview/text"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
)

// TextView describes a popup with scrollable text.
type TextView struct {
	scrollState viewport.ScrollState

	id      string
	content []text.Span
}

// NewTextView returns a new scrollable text view popup.
func NewTextView(id string, content []text.Span) *Container {
	return NewContainer(id, &TextView{id: id, content: content}).
		SetDimensions(layout.Fixed(60), layout.Shrink.Max(40)).
		SetCloseButton()
}

// Init returns a command to run when the model starts, or nil.
func (m *TextView) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m *TextView) Update(msg tview.Msg, rsp *ui.Response) {
	switch ms := msg.(type) {
	case R[scrollMsg]:
		m.scrollState = ms.Value.change
		if ms.Value.msg != nil {
			m.Update(ms.Value.msg, rsp)
		}

	default:
	}
}

// View draws this model onto the screen.
func (m *TextView) View(focused bool) tview.Widget {
	bgStyle := theme.Current().Global
	scroll := scrollbar.New().Style(bgStyle)

	tv := text.Rich(m.content...).
		Style(bgStyle)

	vp := viewport.New(tv, m.scrollState).
		ScrollBar(scroll, viewport.ScrollBarVisibilityAutomatic).
		Axes(layout.Axes{Width: true, Height: true}).
		Focused(focused).
		Keybind(m.ViewportKeybinds).
		OnChange(m.OnViewportChange).
		OnChildMsg(m.OnViewportChildMsg)

	return vp
}

// ViewportKeybinds implements keybindings for the viewport.
func (m *TextView) ViewportKeybinds(msg tview.KeyMsg) viewport.Action {
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
func (m *TextView) OnViewportChange(c *viewport.ScrollState) tview.Msg {
	return NewRouterMsg(m.id, scrollMsg{change: *c})
}

// OnViewportChildMsg processes a message that was passed to its underlying view in the viewport and returns a message.
func (m *TextView) OnViewportChildMsg(c *viewport.ScrollState, msg tview.Msg) tview.Msg {
	return NewRouterMsg(m.id, scrollMsg{change: *c, msg: msg})
}

type scrollMsg struct {
	change viewport.ScrollState
	msg    tview.Msg
}
