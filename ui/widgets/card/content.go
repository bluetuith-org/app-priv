package card

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

// Content describes a card content layout.
type Content struct {
	headerRight, headerLeft text.Span

	content [][]text.Span
	cl      int

	message              [][]text.Span
	cRendered, mRendered bool

	cb ContentBuilder

	width int
}

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

	hd := stext.DrawSpan(screen, c.headerLeft, x, y, area)
	stext.DrawSpan(screen, c.headerRight, x+(width-stext.SpanWidth(c.headerRight)), y, tview.Rectangle{
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
		start := (canvas - stext.LineWidth(line)) / 2

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

func (c *Content) drawLine(screen tview.Screen, line []text.Span, x, y int, area tview.Rectangle) int {
	var textWidth int

	for _, seg := range line {
		textWidth += stext.DrawSpan(screen, seg, x+textWidth, y, area)
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
	c.content = text.Lines(c.cb.Content(), width, text.WrappingWord)
	for _, line := range c.content {
		c.cl = max(c.cl, stext.LineWidth(line))
	}

	c.cRendered = true
}

func (c *Content) refreshMsg(width int) {
	msgMarker, msgContent := c.cb.Message()
	spaceBg := msgContent.Style.GetBackground()
	spaceSeg := stext.NewSpan(" ", tcell.Style{}.Background(spaceBg))

	c.message = text.Lines([]text.Span{msgMarker, spaceSeg, msgContent}, width, text.WrappingGlyph)
	c.mRendered = true
}

func (c *Content) height() int {
	return 3 + len(c.content) + len(c.message)
}

var _ tview.Widget = (*Content)(nil)
