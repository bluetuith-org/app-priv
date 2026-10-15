package ui

import (
	"slices"

	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

// Kb returns the full set of keybindings.
func Kb() *keybindings.Keybindings {
	return keybindings.Current
}

// RenderMergedPill renders a merged pill with multiple segments.
func RenderMergedPill(line richtext.Line, bgColor color.Color) richtext.Line {
	if len(line) == 0 {
		return richtext.Line{}
	}

	firstSeg, lastSeg := line[0], line[len(line)-1]

	firstStyle := tcell.StyleDefault.Foreground(firstSeg.Style.GetBackground()).
		Background(bgColor)

	lastStyle := tcell.StyleDefault.Foreground(lastSeg.Style.GetBackground()).
		Background(bgColor)

	line = slices.Insert(
		line, 0,
		richtext.NewSegment(theme.Icons().PillLeft.String(), firstStyle),
	)

	line = append(line, richtext.NewSegment(theme.Icons().PillRight.String(), lastStyle))

	return line
}
