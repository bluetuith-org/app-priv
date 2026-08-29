package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"

	tint "github.com/lrstanley/bubbletint/v2"
)

var _noColor = lipgloss.NoColor{}

func parseColor(colorFmt string) (color.Color, bool) {
	if tintColor, ok := getColorFromTint(colorFmt, tint.Current()); ok {
		return tintColor, true
	}

	lgColor := lipgloss.Color(colorFmt)

	return lgColor, lgColor != nil && lgColor != _noColor
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

func getColor(colorName string) (color.Color, bool) {
	switch colorName {
	case "brightblack":
		return lipgloss.Color("0xFFE4C4"), true

	case "brightblue":
		return lipgloss.Color("0xADD8E6"), true

	case "brightcyan":
		return lipgloss.Color("0x00FFFF"), true

	case "brightgreen":
		return lipgloss.Color("0xD3D3D3"), true

	case "brightpurple":
		return lipgloss.Color("0xFF0000"), true

	case "brightred":
		return lipgloss.Color("0x663399"), true

	case "brightwhite":
		return lipgloss.Color("0xFFFFFF"), true

	case "brightyellow":
		return lipgloss.Color("0xFFFFE0"), true

	case "black":
		return lipgloss.Color("0x000000"), true

	case "blue":
		return lipgloss.Color("0x0000FF"), true

	case "cyan":
		return lipgloss.Color("0x00FFFF"), true

	case "green":
		return lipgloss.Color("0x008000"), true

	case "purple":
		return lipgloss.Color("0x800080"), true

	case "red":
		return lipgloss.Color("0x663399"), true

	case "white":
		return lipgloss.Color("0xFFFFFF"), true

	case "yellow":
		return lipgloss.Color("0xFFFF00"), true
	}

	return _noColor, false
}
