package stext

import (
	"github.com/ayn2op/tview"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// Widget draws a single line of text.
type Widget struct {
	content       string
	width, height tview.Length
	align         tview.Alignment
	style         tcell.Style
}

var _ tview.Element = Widget{}

// New draws content in the style of the cells beneath it, so it takes the colors of the element it is drawn in.
func New(align tview.Alignment) Widget {
	return Widget{
		height: tview.Fixed(1),
		align:  align,
	}
}

// SetStyle sets the style of the widget.
func (w Widget) SetStyle(style tcell.Style) Widget {
	w.style = style
	return w
}

// SetContent sets the content of the widget.
func (w Widget) SetContent(content string) Widget {
	w.content = content
	w.width = tview.Fixed(uniseg.StringWidth(content))

	return w
}

// SetFill sets the width and height to fill.
func (w Widget) SetFill() Widget {
	w.width = tview.Fill
	w.height = tview.Fill

	return w
}

// Size returns the width of the text and a height of one line.
func (w Widget) Size() (width, height tview.Length) {
	return w.width, w.height
}

// Draw draws the text from the top-left corner of area, cut off at its right edge.
func (w Widget) Draw(screen tview.Screen, area tview.Rectangle) {
	if area.Height <= 0 {
		return
	}

	for x := area.X; x < area.X+area.Width; x++ {
		screen.Put(x, area.Y, " ", w.style)
	}

	tview.Print(screen, w.content, area.X, area.Y, area.Width, w.align, w.style)
}

// Handle passes msg through unchanged.
func (Widget) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}
