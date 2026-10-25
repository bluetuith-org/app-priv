package popup

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/button"
	"github.com/ayn2op/tview/center"
	"github.com/ayn2op/tview/column"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/scrollbar"
	"github.com/ayn2op/tview/text"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
)

// ScrollableText represents a popup with scrollable text.
type ScrollableText struct {
	scrollState viewport.ScrollState

	id      string
	content []text.Span

	closeIcon string
	removeMsg R[popupRemoveMsg]
}

// NewScrollableText returns a new scrollable popup with text.
func NewScrollableText(id string, content []text.Span) *ScrollableText {
	return &ScrollableText{
		id:      id,
		content: content,

		closeIcon: " " + theme.Icons().Close.String() + " ",
		removeMsg: RemoveMsg(id),
	}
}

// ViewID returns the view's ID.
func (m *ScrollableText) ViewID() ui.ViewID {
	return ui.ViewIDNone
}

// Init returns a command to run when the model starts, or nil.
func (m *ScrollableText) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m *ScrollableText) Update(msg tview.Msg, rsp *ui.Response) {
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
func (m *ScrollableText) View(focused bool) tview.Widget {
	bgStyle := theme.Current().Global
	scroll := scrollbar.New().Style(bgStyle)

	t := text.Rich(m.content...).
		Style(bgStyle)

	closeIconText := ui.RenderMergedPill(
		[]text.Span{
			{
				Text:  m.closeIcon,
				Style: theme.Current().Popups.CloseIcon,
			},
		},
		bgStyle.GetBackground(),
	)

	closeIcon := text.Rich(closeIconText...).
		Style(bgStyle)

	closeButton := button.New(closeIcon).
		Style(bgStyle).
		Focused(true).
		OnClick(m.removeMsg)

	vp := viewport.New(t, m.scrollState).
		ScrollBar(scroll, viewport.ScrollBarVisibilityAutomatic).
		Axes(layout.Axes{Width: true, Height: true}).
		Focused(focused).
		Keybind(m.ViewportKeybinds).
		OnChange(m.OnViewportChange).
		OnChildMsg(m.OnViewportChildMsg)

	bx := box.New(vp).
		Background(bgStyle.GetBackground()).
		BorderSet(tview.BorderSetRound()).
		BorderStyle(theme.Current().App.Border).
		Borders(tview.BordersAll).
		Title(closeButton).
		TitleAlignment(tview.HorizontalAlignmentRight)

	centered := center.New(
		column.New(bx).
			Width(m.Width()).
			Height(m.Height()),
	)

	return centered
}

// PopupID returns the unique ID for the popup,
// for messages to be routed to.
func (m *ScrollableText) PopupID() string {
	return m.id
}

// Width returns the width to be occupied by the popup.
func (m *ScrollableText) Width() layout.Length {
	return layout.Fixed(60)
}

// Height returns the height to be occupied by the popup.
func (m *ScrollableText) Height() layout.Length {
	return layout.Shrink.Max(40)
}

// ButtonKeybinds implements keybindings for the close button on the popup.
func (m *ScrollableText) ButtonKeybinds(tview.KeyMsg) button.Action {
	return button.ActionPress
}

// ViewportKeybinds implements keybindings for the viewport.
func (m *ScrollableText) ViewportKeybinds(msg tview.KeyMsg) viewport.Action {
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
func (m *ScrollableText) OnViewportChange(c *viewport.ScrollState) tview.Msg {
	return NewRouterMsg(m.id, infoScrollMsg{change: *c})
}

// OnViewportChildMsg processes a message that was passed to its underlying view in the viewport and returns a message.
func (m *ScrollableText) OnViewportChildMsg(c *viewport.ScrollState, msg tview.Msg) tview.Msg {
	return NewRouterMsg(m.id, infoScrollMsg{change: *c, msg: msg})
}

type infoScrollMsg struct {
	change viewport.ScrollState
	msg    tview.Msg
}

var _ Popup = (*ScrollableText)(nil)
