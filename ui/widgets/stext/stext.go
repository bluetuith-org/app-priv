package stext

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/richtext"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// Widget draws a single line of text.
type Widget struct {
	cr, cl        richtext.Line
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
func (w Widget) Push(line richtext.Line, align tview.Alignment) Widget {
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
		w.cl = append(w.cl, richtext.NewSegment(" ", w.bgStyle))

	case tview.AlignmentRight:
		w.cr = append(w.cr, richtext.NewSegment(" ", w.bgStyle))
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

	textWidth := w.drawLine(screen, w.cl, area.X, area.Y, area)

	w.drawLine(screen, w.cr, area.X+(area.Width-w.cr.Width()), area.Y, tview.Rectangle{
		X:      area.X + textWidth,
		Y:      area.Y,
		Width:  area.Width - textWidth,
		Height: area.Height,
	})
}

// drawLine draws line starting at x, keeping only the cells within area.
//
// Taken from: https://github.com/ayn2op/tview/blob/master/textview/widget.go
func (w Widget) drawLine(screen tview.Screen, line richtext.Line, x, y int, area tview.Rectangle) int {
	var textWidth int

	for _, segment := range line {
		style := segment.Style
		str, state := segment.Text, -1
		for len(str) > 0 {
			var cluster string
			var width int
			cluster, str, width, state = uniseg.FirstGraphemeClusterInString(str, state)
			width = max(width, 1)
			if x >= area.X && x+width <= area.X+area.Width {
				screen.Put(x, y, cluster, style)
			}
			x += width
			textWidth += width
		}
	}

	return textWidth
}

// Handle passes msg through unchanged.
func (Widget) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}
