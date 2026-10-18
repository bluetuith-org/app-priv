package keybindings

import (
	"iter"

	"github.com/ayn2op/tview"
	"github.com/gdamore/tcell/v3"
)

var _emptyKeybinding = Keybinding{}

// KeyID represents a keybinding's ID.
type KeyID int

// Keybinding represents a single Keybinding.
type Keybinding struct {
	TcellKey

	ID      KeyID
	KeyName string
}

// newKeybinding returns a new keybinding.
func newKeybinding(id KeyID, keyCombo TcellKey, keyName string) Keybinding {
	return Keybinding{TcellKey: keyCombo, ID: id, KeyName: keyName}
}

// Matches returns if the keybinding matches the [tea.KeyPressMsg].
func (k Keybinding) Matches(p tview.KeyMsg) bool {
	return p.Key() == k.Key && p.Str() == k.Str && p.Modifiers() == k.Mod
}

// IterKeyMatch describes an iterator that iterates over a keybinding-T pair.
type IterKeyMatch[T any] iter.Seq2[Keybinding, T]

// IterMatch iterates over a set of keybindings and returns whether it matches with the provided [tview.KeyMsg].
func IterMatch[T any](p tview.KeyMsg, iterator IterKeyMatch[T]) (Keybinding, T, bool) {
	var defVal T

	for key, val := range iterator {
		if key.Matches(p) {
			return key, val, true
		}
	}

	return _emptyKeybinding, defVal, false
}

// TcellKey represents a [tcell.EventKey].
type TcellKey struct {
	Key tcell.Key
	Str string
	Mod tcell.ModMask
}

func newTcellKey(key tcell.Key, str string, mod tcell.ModMask) TcellKey {
	return TcellKey{key, str, mod}
}
