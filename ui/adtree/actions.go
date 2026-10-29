package adtree

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
)

// actionStateSpec specifies the current state of the action.
type actionStateSpec uint8

// The different types of action states.
const (
	actionStateNone actionStateSpec = iota
	actionStateEnabled
	actionStateDisabled
)

func (a actionStateSpec) String() string {
	switch a {
	case actionStateEnabled:
		return "On"

	case actionStateDisabled:
		return "Off"
	}

	return ""
}

func (a actionStateSpec) Format(on, off string) string {
	switch a {
	case actionStateEnabled:
		return on

	case actionStateDisabled:
		return off
	}

	return ""
}

// actionUpdateMsg holds an update to the action state matched by the action's ID.
type actionUpdateMsg struct {
	id        nodeID
	stateSpec actionStateSpec
}

// emptyActionUpdateMsg returns an empty action updated message.
func emptyActionUpdateMsg() actionUpdateMsg {
	return actionUpdateMsg{}
}

// updateState updates the action's state from this update message.
func (a *actionUpdateMsg) updateState(stateID nodeID, ad *adActionState) bool {
	if !a.isValid(stateID) {
		return false
	}

	ad.currentState = a.stateSpec
	return true
}

// isValid checks if the message is valid.
func (a *actionUpdateMsg) isValid(stateID nodeID) bool {
	return stateID.IsValid() && a.id == stateID
}

// adActionState holds the state of an executed action.
type adActionState struct {
	key keybindings.Keybinding

	currentState actionStateSpec
	isToggleable bool

	invokeAction actionInvoker
}

type actionInvoker func(*adActionState) (operations.CreationInfo, operations.Invoker)

// newAdActionState creates a new action state, which can be mutated depending on the state of the executed action.
func newAdActionState(key keybindings.Keybinding, state actionStateSpec, isToggleable bool, invoker actionInvoker) *adActionState {
	return &adActionState{key, state, isToggleable, invoker}
}

// actionHandleKeyMsg maps a keybinding to an action.
func actionHandleKeyMsg(p tview.KeyMsg, node *alNode) (ui.Router, bool) {
	_, result, ok := keybindings.IterMatch(p, actionKeyIterator(node))
	if !ok {
		return ui.EmptyRouterMsg(), false
	}

	msg, ok := createOperationMsg(result)
	if !ok {
		return ui.EmptyRouterMsg(), false
	}

	return msg, true
}

func createOperationMsg(acNode *acNode) (ui.Router, bool) {
	id := acNode.AdNode().ID()
	state := acNode.actionState

	if state == nil || !id.IsValid() {
		return ui.EmptyRouterMsg(), false
	}

	createInfo, opFunc := state.invokeAction(state)
	createInfo.UpdateID(id.ToUint64())

	return operations.MsgCreate(createInfo, opFunc), true
}

// actionKeyIterator returns an iterator that iterates over [nodeTypeAction] nodes.
func actionKeyIterator(node *alNode) keybindings.IterKeyMatch[*acNode] {
	return func(yield func(keybindings.Keybinding, *acNode) bool) {
		n := node
		if n == nil {
			return
		}

		ch := node.AdNode().Children()
		if len(ch) == 0 {
			return
		}

		for _, ac := range ch {
			data := nvget[*acNode](ac)
			if !yield(data.actionState.key, data) {
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
