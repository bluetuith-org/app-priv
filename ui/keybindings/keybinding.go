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
	tcellKey

	ID                  KeyID
	ShortHelp, LongHelp string
}

// newKeybinding returns a new keybinding.
func newKeybinding(id KeyID, keyCombo tcellKey, shorthelp, longhelp string) Keybinding {
	return Keybinding{tcellKey: keyCombo, ID: id, ShortHelp: shorthelp, LongHelp: longhelp}
}

// Matches returns if the keybinding matches the [tea.KeyPressMsg].
func (k *Keybinding) Matches(p tview.KeyMsg) bool {
	return p.Key() == k.key && p.Str() == k.str && p.Modifiers() == k.mod
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

type tcellKey struct {
	key tcell.Key
	str string
	mod tcell.ModMask
}

func newTcellKey(key tcell.Key, str string, mod tcell.ModMask) tcellKey {
	return tcellKey{key, str, mod}
}
