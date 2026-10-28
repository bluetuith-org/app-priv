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
	return NewContainer(&TextView{id: id, content: content})
}

// ViewID returns the view's ID.
func (m *TextView) ViewID() ui.ViewID {
	return ui.ViewIDNone
}

// Init returns a command to run when the model starts, or nil.
func (m *TextView) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m *TextView) Update(msg tview.Msg, rsp *ui.Response) {
	switch ms := msg.(type) {
	case R[infoScrollMsg]:
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

// PopupID returns the unique ID for the popup,
// for messages to be routed to.
func (m *TextView) PopupID() string {
	return m.id
}

// Width returns the width to be occupied by the popup.
func (m *TextView) Width() layout.Length {
	return layout.Fixed(60)
}

// Height returns the height to be occupied by the popup.
func (m *TextView) Height() layout.Length {
	return layout.Shrink.Max(40)
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
	return NewRouterMsg(m.id, infoScrollMsg{change: *c})
}

// OnViewportChildMsg processes a message that was passed to its underlying view in the viewport and returns a message.
func (m *TextView) OnViewportChildMsg(c *viewport.ScrollState, msg tview.Msg) tview.Msg {
	return NewRouterMsg(m.id, infoScrollMsg{change: *c, msg: msg})
}

type infoScrollMsg struct {
	change viewport.ScrollState
	msg    tview.Msg
}
