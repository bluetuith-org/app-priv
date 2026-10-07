package card

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/richtext"
	"github.com/ayn2op/tview/textview"
	"github.com/gdamore/tcell/v3"
)

//revive:disable

type Content struct {
	headerRight, headerLeft richtext.Segment

	content richtext.Text
	message richtext.Text

	cRendered, mRendered bool
	cWrapped, mWrapped   richtext.Text

	width int
}

var _ tview.Widget = (*Content)(nil)

func NewCardContent(content richtext.Text) *Content {
	c := &Content{content: content}

	return c
}

func (c *Content) Headers(left, right richtext.Segment) *Content {
	c.headerLeft = left
	c.headerRight = right

	return c
}

func (c *Content) Message(marker richtext.Segment, msg richtext.Segment) *Content {
	spaceBg := msg.Style.GetBackground()
	spaceSeg := richtext.NewSegment(" ", tcell.Style{}.Background(spaceBg))

	c.message = richtext.Text{richtext.Line{marker, spaceSeg, msg}}
	c.mRendered = false

	return c
}

// Draw implements [tview.Widget].
func (c *Content) Draw(screen tview.Screen, area tview.Rectangle) {
	x, y, width, _ := area.X, area.Y, area.Width, area.Height
	_, style, _ := screen.Get(x+width/2, y+area.Height/2)

	hd := tview.Print(screen, c.headerLeft.Text, x, y, width, tview.AlignmentLeft, c.headerLeft.Style)
	tview.Print(screen, c.headerRight.Text, x+hd, y, width-hd, tview.AlignmentRight, c.headerRight.Style)

	y += 2
	c.tv(tview.AlignmentCenter, style, c.cWrapped).Draw(screen, tview.Rectangle{X: x, Y: y, Width: width, Height: len(c.cWrapped) + 1})

	y += len(c.cWrapped) + 1
	c.tv(tview.AlignmentLeft, style, c.mWrapped).Draw(screen, tview.Rectangle{X: x, Y: y, Width: width, Height: len(c.mWrapped)})
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

func (c *Content) ensureSize(width int) int {
	if c.width == width && c.cRendered && c.mRendered {
		return c.height()
	}

	if c.width != width {
		c.cRendered = false
		c.mRendered = false
	}

	if !c.cRendered {
		c.cWrapped = nil
		for _, line := range c.content {
			c.cWrapped = append(c.cWrapped, richtext.WrapWords(line, width)...)
		}

		c.cRendered = true
	}

	if !c.mRendered {
		c.mWrapped = nil
		for _, line := range c.message {
			c.mWrapped = append(c.mWrapped, richtext.WrapWords(line, width)...)
		}

		c.mRendered = true
	}

	c.width = width

	return c.height()
}

func (c *Content) height() int {
	return 3 + len(c.cWrapped) + len(c.mWrapped)
}

func (c *Content) tv(align tview.Alignment, style tcell.Style, content richtext.Text) textview.Widget {
	return textview.New(content).Alignment(align).Wrap(false).Style(style)
}
