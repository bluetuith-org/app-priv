package info

import "github.com/ayn2op/tview/text"

// Content represents the information that is printed on the screen.
type Content interface {
	// Information returns the styled information text.
	Information() []text.Span
}

type cached struct {
	content Content
	spans   []text.Span
}

func newCached(content Content) *cached {
	return &cached{content: content}
}

func (c *cached) reset(content Content) {
	c.content = content
	c.spans = nil
}

func (c *cached) hasContent() bool {
	return c.content != nil
}

func (c *cached) getContent() []text.Span {
	if c.content == nil {
		return nil
	}

	if c.spans == nil {
		c.spans = c.content.Information()
	}

	return c.spans
}
