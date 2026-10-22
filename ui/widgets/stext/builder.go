package stext

import (
	"github.com/ayn2op/tview/text"
	"github.com/gdamore/tcell/v3"
)

// Builder holds a new text builder.
type Builder struct {
	buf     SpanBuilder
	bgStyle tcell.Style

	kstyle, vstyle tcell.Style
}

// NewBuilder creates a new text builder.
func NewBuilder(bgStyle tcell.Style) *Builder {
	return &Builder{bgStyle: bgStyle}
}

// ResetKVStyle resets the styles for key-value pairs.
func (b *Builder) ResetKVStyle(kstyle, vstyle tcell.Style) {
	b.kstyle, b.vstyle = kstyle, vstyle
}

func (b *Builder) AddKV(key, val string) {
	b.buf.Write(key, b.kstyle)
	b.Space()
	b.buf.Write(val, b.vstyle)
	b.Newline()
}

func (b *Builder) AddKVFunc(key string, valFn func(bb *Builder)) {
	b.buf.Write(key, b.kstyle)
	b.Space()
	valFn(b)
	b.buf.NewLine()
}

// AddKVStyled appends a key-value pair.
func (b *Builder) AddKVStyled(key string, keyStyle tcell.Style, val string, valStyle tcell.Style) {
	b.buf.Write(key, keyStyle)
	b.Space()
	b.buf.Write(val, valStyle)
	b.buf.NewLine()
}

// AddKVStyledFunc appends a key-value pair to the builder using a custom function.
func (b *Builder) AddKVStyledFunc(key string, keyStyle tcell.Style, valFn func(bb *Builder)) {
	b.buf.Write(key, keyStyle)
	b.Space()
	valFn(b)
	b.buf.NewLine()
}

// Newline appends a newline.
func (b *Builder) Newline() {
	b.buf.NewLine()
}

// Space appends a space.
func (b *Builder) Space() {
	b.buf.Write(" ", b.bgStyle)
}

// Appendln appends a value-style pair and adds a newline.
func (b *Builder) Appendln(val string, valStyle tcell.Style) {
	b.Append(val, valStyle)
	b.buf.NewLine()
}

// Append appends a single value-style pair to the builder.
func (b *Builder) Append(val string, valStyle tcell.Style) {
	b.buf.Write(val, valStyle)
}

func (b *Builder) AppendSpans(spans ...text.Span) {
	b.buf.WriteSpans(spans...)
}

// Text returns the built text.
func (b *Builder) Text() []text.Span {
	return b.buf.Finish()
}
