package card

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/layout"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

// Widget holds a card and contents.
type Widget struct {
	cc *Content

	bgColor     color.Color
	borderStyle tcell.Style
	borderSet   tview.BorderSet
}

var _ tview.Widget = Widget{}

// New returns a new card.
func New(cc *Content) Widget {
	w := Widget{
		cc: cc,
	}

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

// Layout implements [tview.Widget].
func (w Widget) Layout(limits layout.Limits) layout.Size {
	return layout.Sized(limits, layout.Fill, layout.Shrink, func(l layout.Limits) layout.Size {
		h := w.cc.ensureSize(l.Max.Width)
		return layout.Size{Height: h + 2}
	})
}

// Size implements [tview.Widget].
func (w Widget) Size() (width layout.Length, height layout.Length) {
	return layout.Fill, layout.Shrink
}

// Draw draws the element onto the screen within the given area.
func (w Widget) Draw(screen tview.Screen, area tview.Rectangle) {
	box.New(w.cc).
		Borders(tview.BordersAll).
		BorderSet(w.borderSet).
		BorderStyle(w.borderStyle).
		Background(w.bgColor).
		Draw(screen, area)
}

// Handle implements [tview.Widget].
func (w Widget) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}
