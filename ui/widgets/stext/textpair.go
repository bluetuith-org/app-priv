package stext

import (
	"iter"

	"github.com/ayn2op/tview/richtext"
	"github.com/gdamore/tcell/v3"
	"github.com/rivo/uniseg"
)

// TextPair holds a styled label and content.
type TextPair struct {
	content              richtext.Line
	labelIdx, labelWidth int
}

// Iter iterates over all the segments of this text pair.
func (t *TextPair) Iter() iter.Seq2[richtext.Segment, bool] {
	return func(yield func(richtext.Segment, bool) bool) {
		for i, c := range t.content {
			if !yield(c, i >= t.labelIdx-1) {
				return
			}
		}
	}
}

// LabelWidth returns the label width.
func (t *TextPair) LabelWidth() int {
	return t.labelWidth
}

// TextPairs holds a set of text pairs.
type TextPairs struct {
	maxLabelWidth int
	pairs         []TextPair
}

// NewTextPairs creates a set of text pairs.
func NewTextPairs(pairs ...TextPair) TextPairs {
	ml := 0
	for _, p := range pairs {
		if p.labelWidth > ml {
			ml = p.labelWidth
		}
	}

	return TextPairs{ml, pairs}
}

// Lines returns the text pairs.
func (t *TextPairs) Lines() []TextPair {
	return t.pairs
}

// MaxLabelWidth returns the maximum label width.
func (t *TextPairs) MaxLabelWidth() int {
	return t.maxLabelWidth
}

// TextPairBuilder builds a [TextPair].
type TextPairBuilder struct {
	pairs                richtext.Line
	curr                 bool
	labelIdx, labelWidth int
}

// WriteLabel writes the label text.
func (t *TextPairBuilder) WriteLabel(value string, style tcell.Style) {
	t.write(false, richtext.NewSegment(value, style))
}

// WriteContent writes the content.
func (t *TextPairBuilder) WriteContent(value string, style tcell.Style) {
	t.write(true, richtext.NewSegment(value, style))
}

// Reset resets the builder.
func (t *TextPairBuilder) Reset() {
	t.pairs = nil
	t.curr = false
	t.labelIdx = 0
}

// Finish returns the built text pair.
func (t *TextPairBuilder) Finish() TextPair {
	defer t.Reset()

	return TextPair{t.pairs, t.labelIdx, t.labelWidth}
}

func (t *TextPairBuilder) write(lb bool, seg richtext.Segment) {
	if !t.curr && lb {
		t.labelIdx = len(t.pairs)
		t.curr = true
	}

	if !t.curr {
		t.labelWidth += uniseg.StringWidth(seg.Text)
	}

	t.pairs = append(t.pairs, seg)
}
