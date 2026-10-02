package card

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/layout"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

// Widget holds a card and contents.
type Widget struct {
	item

	bgColor     color.Color
	borderStyle tcell.Style
	borderSet   tview.BorderSet
}

var _ tview.Element = Widget{}

// New returns a new card.
func New(desc stext.TextPairs) Widget {
	w := Widget{
		desc: desc,
	}

	return w
}

// Header sets the header of the card.
func (w Widget) Header(text string, style tcell.Style) Widget {
	w.header.Text = text
	w.header.Style = style

	return w
}

// Status sets the status to be displayed.
func (w Widget) Status(text string, style tcell.Style) Widget {
	w.status.Text = text
	w.status.Style = style

	return w
}

// Message sets the message to be displayed.
func (w Widget) Message(text string, style tcell.Style) Widget {
	w.msg.Text = text
	w.msg.Style = style

	return w
}

// Background sets the background color for the card.
func (w Widget) Background(c color.Color) Widget {
	w.bgColor = c

	return w
}

// BorderStyle sets the border style for the card.
func (w Widget) BorderStyle(style tcell.Style) Widget {
	w.borderStyle = style

	return w
}

// BorderSet sets the border-set for the card.
func (w Widget) BorderSet(set tview.BorderSet) Widget {
	w.borderSet = set
	return w
}

// Rows returns the number of rows this card will occupy.
func (w Widget) Rows(width int) int {
	return w.item.Rows(width) + 2
}

// Layout implements [tview.Element].
func (w Widget) Layout(limits layout.Limits) layout.Size {
	return layout.Atomic(limits, layout.Fill, layout.Fixed(w.Rows(0)))
}

// Size implements [tview.Element].
func (w Widget) Size() (width layout.Length, height layout.Length) {
	return layout.Fill, layout.Fixed(w.Rows(0))
}

// Draw draws the element onto the screen within the given area.
func (w Widget) Draw(screen tview.Screen, area tview.Rectangle) {
	box.New(w.item).
		Height(layout.Fixed(w.item.Rows(0))).
		Borders(tview.BordersAll).
		BorderSet(w.borderSet).
		BorderStyle(w.borderStyle).
		Background(w.bgColor).
		Draw(screen, area)
}
