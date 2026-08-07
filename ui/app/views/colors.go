package views

import (
	"image/color"

	tint "github.com/lrstanley/bubbletint/v2"

	"charm.land/lipgloss/v2"
)

func adaptBright(c color.Color, amount float64) color.Color {
	if tint.Current().Dark {
		return lipgloss.Lighten(c, amount)
	}

	return lipgloss.Darken(c, amount)
}
