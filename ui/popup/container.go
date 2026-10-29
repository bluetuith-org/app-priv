package popup

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/button"
	"github.com/ayn2op/tview/center"
	"github.com/ayn2op/tview/column"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/row"
	"github.com/ayn2op/tview/spacer"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/gdamore/tcell/v3"
)

// Container represents a decorative box within which the underlying widget is drawn.
type Container struct {
	ui.Model

	id               string
	closeIcon, title string

	removeMsg Router

	width, height layout.Length
}

// NewContainer returns a new popup container for the underlying widget.
func NewContainer(id string, content ui.Model) *Container {
	return &Container{
		id:    id,
		Model: content,
	}
}

// PopupID returns the popup's unique ID, for messages to be routed to.
func (m *Container) PopupID() string {
	return m.id
}

// SetDimensions sets the layout for the popup.
func (m *Container) SetDimensions(width, height layout.Length) *Container {
	m.width, m.height = width, height
	return m
}

// SetTitle sets up the title for the popup.
func (m *Container) SetTitle(icon, title string) *Container {
	if title != "" {
		m.title = " " + icon + " " + title + " "
	}

	return m
}

// SetCloseButton sets up the close button to close the popup.
func (m *Container) SetCloseButton() *Container {
	m.closeIcon = " " + theme.Icons().Close.String() + " "
	m.removeMsg = RemoveMsg(m.PopupID())

	return m
}

// View draws this model onto the screen.
func (m *Container) View(focused bool) tview.Widget {
	bgStyle := theme.Current().Global

	titleWidget := m.boxTitleWidget(bgStyle)

	sp := spacer.New().
		Width(layout.Fill).
		Height(layout.Fixed(1))

	col := column.New(sp).
		Push(m.Model.View(focused))

	bx := box.New(col).
		Background(bgStyle.GetBackground()).
		BorderSet(tview.BorderSetRound()).
		BorderStyle(theme.Current().App.Border).
		Borders(tview.BordersAll).
		Title(titleWidget).
		TitleAlignment(tview.HorizontalAlignmentRight)

	centered := center.New(
		column.New(bx).
			Width(m.width).
			Height(m.height),
	)

	return centered
}

func (m *Container) boxTitleWidget(bgStyle tcell.Style) row.Widget {
	sp := newSpacer(bgStyle)
	rw := row.New().
		Width(layout.Fill).
		Height(layout.Shrink).
		Push(sp)

	pushed := 0

	if tn, ok := m.renderTitle(bgStyle); ok {
		rw = rw.Push(tn)
		pushed++
	}
	if bn, ok := m.renderCloseButton(bgStyle); ok {
		if pushed > 0 {
			rw = rw.Push(sp)
		}

		rw = rw.Push(bn)
		pushed++
	}

	return rw
}

func (m *Container) renderTitle(bgStyle tcell.Style) (tn text.Widget, ok bool) {
	if m.title == "" {
		return
	}

	titleText := ui.RenderMergedPill(
		[]text.Span{
			{
				Text:  m.title,
				Style: theme.Current().Popups.Title,
			},
		},
		bgStyle.GetBackground(),
	)

	ok = true
	tn = text.Rich(titleText...).
		Style(bgStyle)

	return
}

func (m *Container) renderCloseButton(bgStyle tcell.Style) (bn button.Widget, ok bool) {
	if m.closeIcon == "" {
		return
	}

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

	ok = true
	bn = button.New(closeIcon).
		Style(bgStyle).
		Focused(true).
		OnClick(m.removeMsg)

	return
}

var _ ui.Model = (*Container)(nil)
