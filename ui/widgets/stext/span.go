package stext

import (
	"strings"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/rockorager/go-uucode"
)

// NewSpan returns a new styled text container.
func NewSpan(txt string, style tview.Style) text.Span {
	return text.Span{Text: txt, Style: style}
}

// SpanWidth returns the width of the [text.Span].
func SpanWidth(span text.Span) int {
	return uucode.StringWidth(span.Text)
}

// DrawSpan draws a segment of styled text onto the screen.
func DrawSpan(screen tview.Screen, span text.Span, x, y int, area tview.Rectangle) int {
	textWidth := 0

	style := span.Style
	// The cells beneath a span mostly share a style, so it is merged once for each run of them.
	var beneath, merged tview.Style
	known := false
	clusters := uucode.NewGraphemeWidthIterator(span.Text)
	for c, ok := clusters.Next(); ok; c, ok = clusters.Next() {
		cluster, width := span.Text[c.Start:c.End], max(c.Width, 1)
		// What is past the right edge is not drawn.
		if x >= area.X+area.Width {
			return textWidth
		}
		if x >= area.X && x+width <= area.X+area.Width {
			if _, under, _ := screen.Get(x, y); !known || under != beneath {
				beneath, merged, known = under, tview.MergeStyle(under, style), true
			}
			screen.Put(x, y, cluster, merged)
		}
		x += width
		textWidth += width
	}

	return textWidth
}

// SpanBuilder builds the spans of a text.
type SpanBuilder struct {
	spans []text.Span
}

// Write appends value in style.
func (b *SpanBuilder) Write(value string, style tview.Style) {
	if value != "" {
		b.spans = append(b.spans, text.Span{Text: value, Style: style})
	}
}

// WriteSpans appends spans.
func (b *SpanBuilder) WriteSpans(spans ...text.Span) {
	b.spans = append(b.spans, spans...)
}

// LineEmpty reports whether nothing was written since the last line break.
func (b *SpanBuilder) LineEmpty() bool {
	return len(b.spans) == 0 || strings.HasSuffix(b.spans[len(b.spans)-1].Text, "\n")
}

// NewLine starts a new line.
func (b *SpanBuilder) NewLine() {
	b.Write("\n", tview.Style{})
}

// Finish returns the spans, without a line break at the end, which would add an empty line.
func (b *SpanBuilder) Finish() []text.Span {
	if n := len(b.spans); n > 0 {
		last := &b.spans[n-1]
		if last.Text = strings.TrimSuffix(last.Text, "\n"); last.Text == "" {
			b.spans = b.spans[:n-1]
		}
	}
	return b.spans
}

// NewLine returns a slice of styled spans.
func NewLine(span ...text.Span) []text.Span {
	return span
}

// LineWidth returns the width of the styled line.
func LineWidth(line []text.Span) int {
	return text.Width(line)
}

// DrawLine draws line starting at x, keeping only the cells within area.
func DrawLine(screen tview.Screen, line []text.Span, x, y int, area tview.Rectangle) int {
	textWidth := 0

	for _, span := range line {
		textWidth += DrawSpan(screen, span, x+textWidth, y, area)
	}

	return textWidth
}
