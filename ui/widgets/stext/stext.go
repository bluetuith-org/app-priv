package stext

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/text"
	"github.com/gdamore/tcell/v3"
)

// Widget draws a single line of text.
type Widget struct {
	cr, cl        []text.Span
	width, height layout.Length
	align         tview.Alignment
	bgStyle       tcell.Style
}

var _ tview.Widget = Widget{}

// New draws content in the style of the cells beneath it, so it takes the colors of the element it is drawn in.
func New() Widget {
	return Widget{
		height: layout.Fixed(1),
		width:  layout.Fill,
	}
}

// Push pushes a line of content with alignment.
func (w Widget) Push(line []text.Span, align tview.Alignment) Widget {
	switch align {
	case tview.AlignmentLeft:
		w.cl = append(w.cl, line...)

	case tview.AlignmentRight:
		w.cr = append(w.cr, line...)
	}

	return w
}

// PushSpace appends a space between segments.
func (w Widget) PushSpace(align tview.Alignment) Widget {
	switch align {
	case tview.AlignmentLeft:
		w.cl = append(w.cl, NewSpan(" ", w.bgStyle))

	case tview.AlignmentRight:
		w.cr = append(w.cr, NewSpan(" ", w.bgStyle))
	}

	return w
}

// Background sets the background style.
func (w Widget) Background(bgStyle tcell.Style) Widget {
	w.bgStyle = bgStyle

	return w
}

// Size returns the width of the text and a height of one line.
func (w Widget) Size() (width, height layout.Length) {
	return w.width, w.height
}

// Layout returns the size of the text within limits.
func (w Widget) Layout(limits layout.Limits) layout.Size {
	return layout.Atomic(limits, w.width, w.height)
}

// Draw draws the text from the top-left corner of area, cut off at its right edge.
func (w Widget) Draw(screen tview.Screen, area tview.Rectangle) {
	if area.Height <= 0 {
		return
	}

	screen.FillArea(area.X, area.Y, area.Width, area.Height, ' ', w.bgStyle)

	textWidth := DrawLine(screen, w.cl, area.X, area.Y, area)

	DrawLine(screen, w.cr, area.X+(area.Width-LineWidth(w.cr)), area.Y, tview.Rectangle{
		X:      area.X + textWidth,
		Y:      area.Y,
		Width:  area.Width - textWidth,
		Height: area.Height,
	})
}

// Handle passes msg through unchanged.
func (Widget) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}
