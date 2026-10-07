package ui

import (
	"github.com/ayn2op/tview/richtext"
	"github.com/gdamore/tcell/v3"
)

//revive:disable

// Builder holds a new text builder.
type Builder struct {
	buf     richtext.Builder
	bgStyle tcell.Style
}

// NewBuilder creates a new text builder.
func NewBuilder(bgStyle tcell.Style) *Builder {
	return &Builder{bgStyle: bgStyle}
}

// AddKV appends a key-value pair.
func (b *Builder) AddKV(key string, keyStyle tcell.Style, val string, valStyle tcell.Style) {
	b.buf.Write(key, keyStyle)
	b.Space()
	b.buf.Write(val, valStyle)
	b.buf.NewLine()
}

// AddKVFunc appends a key-value pair to the builder using a custom function.
func (b *Builder) AddKVFunc(key string, keyStyle tcell.Style, valFn func(bb *Builder)) {
	b.buf.Write(key, keyStyle)
	b.Space()
	valFn(b)
	b.buf.NewLine()
}

// Newline appends a newline.
func (b *Builder) Newline() {
	b.buf.Write("", b.bgStyle)
	b.buf.NewLine()
}

// Space appends a space.
func (b *Builder) Space() {
	b.buf.Write(" ", b.bgStyle)
}

// AppendLine appends a [richtext.Line].
func (b *Builder) AppendLine(line richtext.Line) {
	b.buf.WriteText(richtext.Text{line})
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

// Text returns the built text.
func (b *Builder) Text() richtext.Text {
	return b.buf.Finish()
}
