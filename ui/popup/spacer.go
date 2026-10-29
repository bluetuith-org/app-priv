package popup

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/gdamore/tcell/v3"
)

// Spacer is an empty area that draws nothing and handles no messages.
type Spacer struct {
	style tcell.Style
}

var _ tview.Widget = Spacer{}

// newSpacer returns a spacer that takes no space until it is given a width or height.
func newSpacer(style tcell.Style) Spacer {
	return Spacer{style: style}
}

// Size returns the width and height of the spacer.
func (w Spacer) Size() (width, height layout.Length) {
	return layout.Fill, layout.Fill
}

// Layout returns the size of the spacer within limits.
func (w Spacer) Layout(limits layout.Limits) layout.Size {
	return layout.Atomic(limits, layout.Fill, layout.Fill)
}

// Draw draws nothing.
func (w Spacer) Draw(screen tview.Screen, area tview.Rectangle) {
	screen.FillArea(area.X, area.Y, area.Width, 1, '─', w.style)
}

// Handle returns msg unchanged.
func (Spacer) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg { return msg }
