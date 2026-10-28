package popup

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/button"
	"github.com/ayn2op/tview/center"
	"github.com/ayn2op/tview/column"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
)

// Container represents a decorative box within which the underlying widget is drawn.
type Container struct {
	Overlay

	closeIcon string
	removeMsg R[popupRemoveMsg]
}

// NewContainer returns a new popup container for the underlying widget.
func NewContainer(content Overlay) *Container {
	id := content.PopupID()
	return &Container{
		Overlay: content,

		closeIcon: " " + theme.Icons().Close.String() + " ",
		removeMsg: RemoveMsg(id),
	}
}

// View draws this model onto the screen.
func (m *Container) View(focused bool) tview.Widget {
	bgStyle := theme.Current().Global

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

	bx := box.New(m.Overlay.View(focused)).
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

var _ Overlay = (*Container)(nil)
