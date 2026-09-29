package adtree

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
)

// actionStateSpec specifies the current state of the action.
type actionStateSpec uint8

// The different types of action states.
const (
	actionStateNone actionStateSpec = iota
	actionStateEnabled
	actionStateDisabled
)

// actionUpdateMsg holds an update to the action state matched by the action's ID.
type actionUpdateMsg struct {
	id        string
	stateSpec actionStateSpec
}

// emptyActionUpdateMsg returns an empty action updated message.
func emptyActionUpdateMsg() actionUpdateMsg {
	return actionUpdateMsg{}
}

// updateState updates the action's state from this update message.
func (a *actionUpdateMsg) updateState(stateID string, ad *adActionState) bool {
	if !a.isValid(stateID) {
		return false
	}

	ad.currentState = a.stateSpec
	return true
}

// isValid checks if the message is valid.
func (a *actionUpdateMsg) isValid(stateID string) bool {
	return a.id != "" && a.id == stateID
}

// adActionState holds the state of an executed action.
type adActionState struct {
	key keybindings.Keybinding

	currentState actionStateSpec
	isToggleable bool
}

// newAdActionState creates a new action state, which can be mutated depending on the state of the executed action.
func newAdActionState(key keybindings.Keybinding, state actionStateSpec, isToggleable bool) *adActionState {
	return &adActionState{key, state, isToggleable}
}

// actionHandleKeyMsg maps a keybinding to an action.
func actionHandleKeyMsg(p tview.KeyMsg, node *ctree.Node) (ui.RouterMsg, bool) {
	_, result, ok := keybindings.IterMatch(p, actionKeyIterator(node))
	if !ok {
		return ui.EmptyRouterMsg(), false
	}

	return acStateToOpMsg(result.id, result.state), true
}

// TODO:
func acStateToOpMsg(id string, state *adActionState) ui.RouterMsg {
	return ui.EmptyRouterMsg()
}

// actionKeyIterResult holds the action-key iteration result.
type actionKeyIterResult struct {
	state *adActionState
	id    string
}

// actionKeyIterator returns an iterator that iterates over [nodeTypeAction] nodes.
func actionKeyIterator(node *ctree.Node) keybindings.IterKeyMatch[*actionKeyIterResult] {
	return func(yield func(keybindings.Keybinding, *actionKeyIterResult) bool) {
		n := node
		if n == nil {
			return
		}

		ch := node.Children()
		if len(ch) <= int(relPosActionsListNode) {
			return
		}

		actionListNode := ch[relPosActionsListNode]
		actionNodes := actionListNode.Children()

		res := &actionKeyIterResult{}
		for _, ac := range actionNodes {
			data := ac.Data[*adNode]()

			res.id = data.ID()
			res.state = data.actionState

			if !yield(data.actionState.key, res) {
				return
			}
		}
	}
}

// boolToActionState converts a boolean value to an action state.
func boolToActionState(val bool) actionStateSpec {
	if !val {
		return actionStateDisabled
	}

	return actionStateEnabled
}
