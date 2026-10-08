package ui

import (
	"slices"

	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/gdamore/tcell/v3"
)

// Kb returns the full set of keybindings.
func Kb() *keybindings.Keybindings {
	return keybindings.Current
}

// RenderMergedPill renders a merged pill with multiple segments.
func RenderMergedPill(line richtext.Line) richtext.Line {
	if len(line) == 0 {
		return richtext.Line{}
	}

	firstSeg, lastSeg := line[0], line[len(line)-1]

	firstStyle := tcell.StyleDefault.Foreground(firstSeg.Style.GetBackground())
	lastStyle := tcell.StyleDefault.Foreground(lastSeg.Style.GetBackground())

	line = slices.Insert(
		line, 0,
		richtext.NewSegment(theme.Icons().PillLeft.String(), firstStyle),
	)

	line = append(line, richtext.NewSegment(theme.Icons().PillRight.String(), lastStyle))

	return line
}
