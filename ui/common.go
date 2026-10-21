package ui

import (
	"github.com/bluetuith-org/bluetuith/keybindings"
)

// Kb returns the full set of keybindings.
func Kb() *keybindings.Keybindings {
	return keybindings.Current
}
