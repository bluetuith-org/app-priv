package card

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
)

type item struct {
	desc stext.TextPairs

	header, msg, status richtext.Segment
}

// Layout implements [tview.Element].
func (i item) Layout(limits layout.Limits) layout.Size {
	return layout.Atomic(limits, layout.Fill, layout.Fill)
}

// Size implements [tview.Element].
func (i item) Size() (width layout.Length, height layout.Length) {
	return layout.Fill, layout.Fill
}

// Draw draws the element onto the screen within the given area.
func (i item) Draw(screen tview.Screen, area tview.Rectangle) {
	x, y, width, height := area.X, area.Y, area.Width, area.Height
	if width <= 0 || height <= 0 {
		return
	}

	drawnLine := y

	if i.header.Text != "" {
		tview.Print(screen, i.header.Text, x, drawnLine, width, tview.AlignmentLeft, i.header.Style)

		drawnLine++
	}

	for _, line := range i.desc.Lines() {
		if drawnLine >= y+height {
			return
		}

		drawn := 0

		w := width / 2
		drawingContent := false

		contentX := x + w
		labelX := x + max(0, w-i.desc.MaxLabelWidth())

		for seg, isContent := range line.Iter() {
			dx := labelX
			if isContent {
				dx = contentX
				if !drawingContent {
					w = width / 2
					drawn = 0
					drawingContent = true
				}
			}

			drawn += tview.Print(screen, seg.Text, dx+drawn, drawnLine, w-4, tview.AlignmentLeft, seg.Style)
			if drawn >= w-4 {
				break
			}
		}

		drawnLine++
	}

	if drawnLine >= y+height {
		return
	}

	wd := tview.Print(screen, i.status.Text, x, drawnLine, width, tview.AlignmentLeft, i.status.Style)
	wd++

	tview.Print(screen, i.msg.Text, x+wd, drawnLine, width-wd, tview.AlignmentLeft, i.msg.Style)
}

// Handle translates a message, such as input received within the given area, into the message passed to Update. It returns nil to drop the message.
func (i item) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}

// Rows returns the number of rows needed for this item, provided a set width.
func (i item) Rows(int) int {
	count := 0
	if i.header.Text != "" {
		count++
	}
	if i.msg.Text != "" || i.status.Text != "" {
		count++
	}

	return len(i.desc.Lines()) + count
}
