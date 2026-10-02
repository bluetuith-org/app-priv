package operations

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/bluetuith-org/bluetuith/theme"
)

type cardItem struct {
	box.Widget
	it *item
}

func newCardItem(running *RunningInfo) cardItem {
	it := newItem(running)
	bx := box.New(it).
		BorderSet(tview.BorderSetRound()).
		Borders(tview.BordersAll).
		Height(tview.Fixed(it.Rows(0) + 2)).
		Background(theme.Current().Global.GetBackground())

	c := cardItem{
		it:     it,
		Widget: bx,
	}

	return c
}

func (c cardItem) Rows(width int) int {
	return c.it.Rows(width) + 2
}

type item struct {
	running *RunningInfo
}

func newItem(running *RunningInfo) *item {
	return &item{running}
}

// Draw draws the element onto the screen within the given area.
func (i *item) Draw(screen tview.Screen, area tview.Rectangle) {
	x, y, width, height := area.X, area.Y, area.Width, area.Height

	for _, line := range i.running.description {
		if y >= y+height {
			return
		}

		drawn := 0
		for _, seg := range line {
			drawn += tview.Print(screen, seg.Text, x, y, width, tview.AlignmentLeft, seg.Style)
			if drawn >= x+width {
				break
			}
		}

		y++
	}
}

// Handle translates a message, such as input received within the given area, into the message passed to Update. It returns nil to drop the message.
func (i *item) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}

// Rows returns the number of rows needed for this item, provided a set width.
func (i *item) Rows(int) int {
	return len(i.running.description)
}
