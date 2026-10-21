package ui

import (
	"slices"

	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

// Kb returns the full set of keybindings.
func Kb() *keybindings.Keybindings {
	return keybindings.Current
}

// RenderMergedPill renders a merged pill with multiple segments.
func RenderMergedPill(line []text.Span, bgColor color.Color) []text.Span {
	if len(line) == 0 {
		return []text.Span{}
	}

	firstSeg, lastSeg := line[0], line[len(line)-1]

	firstStyle := tcell.StyleDefault.Foreground(firstSeg.Style.GetBackground()).
		Background(bgColor)

	lastStyle := tcell.StyleDefault.Foreground(lastSeg.Style.GetBackground()).
		Background(bgColor)

	line = slices.Insert(
		line, 0,
		stext.NewSpan(theme.Icons().PillLeft.String(), firstStyle),
	)

	line = append(line, stext.NewSpan(theme.Icons().PillRight.String(), lastStyle))

	return line
}
