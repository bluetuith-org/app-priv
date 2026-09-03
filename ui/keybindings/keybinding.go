package keybindings

import (
	"iter"

	tea "charm.land/bubbletea/v2"
)

// KeyID represents a keybinding's ID.
type KeyID int

// Keybinding represents a single Keybinding.
type Keybinding struct {
	ID                       KeyID
	Key, ShortHelp, LongHelp string
}

// NewKeybinding returns a new keybinding.
func NewKeybinding(id KeyID, keyCombo, shorthelp, longhelp string) Keybinding {
	return Keybinding{ID: id, Key: keyCombo, ShortHelp: shorthelp, LongHelp: longhelp}
}

// Matches returns if the keybinding matches the [tea.KeyPressMsg].
func (k *Keybinding) Matches(p tea.KeyPressMsg) bool {
	return p.String() == k.Key
}

// IterKeyMatch describes an iterator that iterates over a keybinding-T pair.
type IterKeyMatch[T any] iter.Seq2[Keybinding, T]

// IterMatch iterates over a set of keybindings and returns whether it matches with the provided [tea.KeyPressMsg].
func IterMatch[T any](p tea.KeyPressMsg, iterator IterKeyMatch[T]) (Keybinding, T, bool) {
	var defVal T

	for key, val := range iterator {
		if key.Key == p.String() {
			return key, val, true
		}
	}

	return Keybinding{}, defVal, false
}
