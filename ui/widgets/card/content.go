package card

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/richtext"
	"github.com/ayn2op/tview/textview"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

//revive:disable

type Content struct {
	headerRight, headerLeft richtext.Segment

	content              richtext.Text
	message              richtext.Text
	cRendered, mRendered bool

	cb ContentBuilder

	width  int
	inited bool
}

var _ tview.Widget = (*Content)(nil)

func NewCardContent(cb ContentBuilder) *Content {
	c := &Content{cb: cb}

	return c
}

// Refresh refreshes the entire content of the card.
func (c *Content) Refresh() *Content {
	c.cRendered = false
	c.mRendered = false

	return c
}

// RefreshMessage refreshes the message.
func (c *Content) RefreshMessage() *Content {
	c.mRendered = false

	return c
}

// Draw implements [tview.Widget].
func (c *Content) Draw(screen tview.Screen, area tview.Rectangle) {
	x, y, width, _ := area.X, area.Y, area.Width, area.Height
	_, style, _ := screen.Get(x+width/2, y+area.Height/2)

	hd := c.drawSegment(screen, c.headerLeft, x, y, area)
	c.drawSegment(screen, c.headerRight, x+(width-c.headerRight.Width()), y, tview.Rectangle{
		X:      x + hd,
		Y:      y,
		Width:  width - hd,
		Height: area.Height,
	})

	y += 2
	c.tv(tview.AlignmentCenter, style, c.content).Draw(screen, tview.Rectangle{X: x, Y: y, Width: width, Height: len(c.content) + 1})

	y += len(c.content) + 1
	c.tv(tview.AlignmentLeft, style, c.message).Draw(screen, tview.Rectangle{X: x, Y: y, Width: width, Height: len(c.message)})
}

// Handle implements [tview.Widget].
func (c *Content) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}

// Layout implements [tview.Widget].
func (c *Content) Layout(limits layout.Limits) layout.Size {
	return layout.Atomic(limits, layout.Fill, layout.Fixed(c.height()))
}

// Size implements [tview.Widget].
func (c *Content) Size() (width layout.Length, height layout.Length) {
	return layout.Fill, layout.Shrink
}

// drawSegment draws a segment starting at x, keeping only the cells within area.
func (c *Content) drawSegment(screen tview.Screen, segment richtext.Segment, x, y int, area tview.Rectangle) int {
	var textWidth int

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

	return textWidth
}

func (c *Content) ensureSize(width int) int {
	if !c.inited {
		c.refreshHeaders()
		c.refreshMessage(width)
		c.refreshBody(width)

		c.inited = true
		c.width = width

		return c.height()
	}

	if c.width == width && c.cRendered && c.mRendered {
		return c.height()
	}

	if c.width != width {
		c.cRendered = false
		c.mRendered = false
	}

	if !c.cRendered {
		c.refreshBody(width)
	}

	if !c.mRendered {
		c.refreshMessage(width)
	}

	c.width = width

	return c.height()
}

func (c *Content) refreshHeaders() *Content {
	c.headerLeft, c.headerRight = c.cb.Headers()

	return c
}

func (c *Content) refreshBody(width int) {
	c.content = nil

	for _, line := range c.cb.Content() {
		c.content = append(c.content, richtext.WrapWords(line, width)...)
	}

	c.cRendered = true
}

func (c *Content) refreshMessage(width int) {
	msgMarker, msgContent := c.cb.Message()
	spaceBg := msgContent.Style.GetBackground()
	spaceSeg := richtext.NewSegment(" ", tcell.Style{}.Background(spaceBg))

	c.message = richtext.WrapWords(richtext.Line{msgMarker, spaceSeg, msgContent}, width)
	c.mRendered = true
}

func (c *Content) height() int {
	return 3 + len(c.content) + len(c.message)
}

func (c *Content) tv(align tview.Alignment, style tcell.Style, content richtext.Text) textview.Widget {
	return textview.New(content).Alignment(align).Wrap(false).Style(style)
}
