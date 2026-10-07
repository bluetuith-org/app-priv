package ui

import (
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/gdamore/tcell/v3"
)

// Kb returns the full set of keybindings.
func Kb() *keybindings.Keybindings {
	return keybindings.Current
}

// RenderPillIcon renders a pill icon with text.
func RenderPillIcon(text string, style tcell.Style) richtext.Line {
	bgStyle := tcell.StyleDefault.Foreground(style.GetBackground())

	return richtext.Line{
		richtext.NewSegment(theme.Icons().PillLeft.String(), bgStyle),
		richtext.NewSegment(text, style),
		richtext.NewSegment(theme.Icons().PillRight.String(), bgStyle),
	}
}
