package stext

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/richtext"
	"github.com/gdamore/tcell/v3"
)

type alignedSegment struct {
	richtext.Segment
	align tview.Alignment
}

// Widget draws a single line of text.
type Widget struct {
	content       []alignedSegment
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
	for _, seg := range line {
		w.content = append(w.content, alignedSegment{Segment: seg, align: align})
	}

	return w
}

// PushSpace appends a space between segments.
func (w Widget) PushSpace() Widget {
	w.content = append(w.content, alignedSegment{Segment: richtext.NewSegment(" ", w.bgStyle)})

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

	var textWidth int
	for _, seg := range w.content {
		if textWidth >= area.Width {
			return
		}

		textWidth += tview.Print(screen, seg.Text, area.X+textWidth, area.Y, area.Width-textWidth, seg.align, seg.Style)
	}
}

// Handle passes msg through unchanged.
func (Widget) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}
