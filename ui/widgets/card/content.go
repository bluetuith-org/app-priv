package card

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/richtext"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// Content describes a card content layout.
type Content struct {
	headerRight, headerLeft richtext.Segment

	content richtext.Text
	cl      int

	message              richtext.Text
	cRendered, mRendered bool

	cb ContentBuilder

	width int
}

var _ tview.Widget = (*Content)(nil)

// NewCardContent returns new [Content].
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

	hd := c.drawSegment(screen, c.headerLeft, x, y, area)
	c.drawSegment(screen, c.headerRight, x+(width-c.headerRight.Width()), y, tview.Rectangle{
		X:      x + hd,
		Y:      y,
		Width:  width - hd,
		Height: area.Height,
	})

	y += 2
	cArea := tview.Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: len(c.content) + 1,
	}

	vstart := max(c.cl-width, 0) / 2
	canvas := max(c.cl, area.Width)
	for _, line := range c.content {
		start := (canvas - line.Width()) / 2

		c.drawLine(screen, line, cArea.X+start-vstart, cArea.Y, cArea)
		cArea.Y++
	}

	y += cArea.Height
	mArea := tview.Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: len(c.message),
	}

	for _, line := range c.message {
		c.drawLine(screen, line, mArea.X, mArea.Y, mArea)
		mArea.Y++
	}
}

// Handle implements [tview.Widget].
func (c *Content) Handle(msg tview.Msg, _ tview.Rectangle) tview.Msg {
	return msg
}

// Layout implements [tview.Widget].
func (c *Content) Layout(limits layout.Limits) layout.Size {
	return layout.Sized(limits, layout.Fill, layout.Shrink, func(l layout.Limits) layout.Size {
		h := c.ensureSize(l.Max.Width)
		return layout.Size{Height: h}
	})
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

func (c *Content) drawLine(screen tview.Screen, line richtext.Line, x, y int, area tview.Rectangle) int {
	var textWidth int

	for _, seg := range line {
		textWidth += c.drawSegment(screen, seg, x+textWidth, y, area)
	}

	return textWidth
}

func (c *Content) ensureSize(width int) int {
	if c.width == width && c.cRendered && c.mRendered {
		return c.height()
	}

	if c.width != width {
		c.cRendered = false
		c.mRendered = false
	}

	if !c.cRendered {
		c.refreshHeaders()
		c.refreshBody(width)
	}

	if !c.mRendered {
		c.refreshMsg(width)
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

	c.cl = 0
	for _, line := range c.cb.Content() {
		c.content = append(c.content, richtext.WrapWords(line, width)...)
		c.cl = max(c.cl, line.Width())
	}

	c.cRendered = true
}

func (c *Content) refreshMsg(width int) {
	msgMarker, msgContent := c.cb.Message()
	spaceBg := msgContent.Style.GetBackground()
	spaceSeg := richtext.NewSegment(" ", tcell.Style{}.Background(spaceBg))

	c.message = richtext.WrapWords(richtext.Line{msgMarker, spaceSeg, msgContent}, width)
	c.mRendered = true
}

func (c *Content) height() int {
	return 3 + len(c.content) + len(c.message)
}
