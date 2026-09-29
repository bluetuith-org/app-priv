package theme

import (
	"cmp"
	"image/color"

	"github.com/gdamore/tcell/v3"
	tc "github.com/gdamore/tcell/v3/color"
	tint "github.com/lrstanley/bubbletint/v2"
)

var _noColor = tcell.ColorDefault

func parseColor(colorFmt string) (color.Color, bool) {
	if tintColor, ok := getColorFromTint(colorFmt, tint.Current()); ok {
		return tintColor, true
	}

	lgColor := tc.GetColor(colorFmt)

	return lgColor, lgColor != _noColor
}

// darkenColor takes a color and makes it darker by a specific percentage (0-1, clamped).
// Taken from: https://github.com/charmbracelet/lipgloss/blob/main/color.go
func darkenColor(c color.Color, percent float64) color.Color {
	if c == nil {
		return nil
	}

	mult := 1.0 - clamp(percent, 0, 1)

	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(float64(r>>8) * mult),
		G: uint8(float64(g>>8) * mult),
		B: uint8(float64(b>>8) * mult),
		A: uint8(min(255, float64(a>>8))),
	}
}

// lightenColor makes a color lighter by a specific percentage (0-1, clamped).
// Taken from: https://github.com/charmbracelet/lipgloss/blob/main/color.go
func lightenColor(c color.Color, percent float64) color.Color {
	if c == nil {
		return nil
	}

	add := 255 * clamp(percent, 0, 1)

	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(min(255, float64(r>>8)+add)),
		G: uint8(min(255, float64(g>>8)+add)),
		B: uint8(min(255, float64(b>>8)+add)),
		A: uint8(min(255, float64(a>>8))),
	}
}

func getColorFromTint(colorName string, selectedTint *tint.Tint) (color.Color, bool) {
	switch colorName {
	case "brightblack":
		return selectedTint.BrightBlack, true

	case "brightblue":
		return selectedTint.BrightBlue, true

	case "brightcyan":
		return selectedTint.BrightCyan, true

	case "brightgreen":
		return selectedTint.BrightGreen, true

	case "brightpurple":
		return selectedTint.BrightPurple, true

	case "brightred":
		return selectedTint.BrightRed, true

	case "brightwhite":
		return selectedTint.BrightWhite, true

	case "brightyellow":
		return selectedTint.BrightYellow, true

	case "black":
		return selectedTint.Black, true

	case "blue":
		return selectedTint.Blue, true

	case "cyan":
		return selectedTint.Cyan, true

	case "green":
		return selectedTint.Green, true

	case "purple":
		return selectedTint.Purple, true

	case "red":
		return selectedTint.Red, true

	case "white":
		return selectedTint.White, true

	case "yellow":
		return selectedTint.Yellow, true
	}

	return _noColor, false
}

//go:inline
func clamp[T cmp.Ordered](v, low, high T) T {
	if high < low {
		high, low = low, high
	}
	return min(high, max(low, v))
}
