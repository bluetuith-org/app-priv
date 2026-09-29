package ctree

import (
	"iter"

	"github.com/ayn2op/tview"
	"github.com/gdamore/tcell/v3"
)

// Provider describes a formatter to print tree-node contents.
type Provider interface {
	// Keybind converts keybindings to actions.
	Keybind(msg tview.KeyMsg) (Action, bool)

	// OnChange processes a change in the tree's state and returns a [tview.Msg]
	OnChange(chg Change) tview.Msg

	// MarkerStyle returns the styles for the markers.
	MarkerStyle(node *Node) tcell.Style

	// StyledLabels returns a sequence of text and their associated styles.
	StyledLabels(node *Node, selected bool) iter.Seq2[string, tcell.Style]
}
