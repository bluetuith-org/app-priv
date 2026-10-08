package adtree

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/operations"
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

	invokeAction actionInvoker
}

type actionInvoker func(*adActionState) (operations.CreationInfo, operations.Invoker)

// newAdActionState creates a new action state, which can be mutated depending on the state of the executed action.
func newAdActionState(key keybindings.Keybinding, state actionStateSpec, isToggleable bool, invoker actionInvoker) *adActionState {
	return &adActionState{key, state, isToggleable, invoker}
}

// actionHandleKeyMsg maps a keybinding to an action.
func actionHandleKeyMsg(p tview.KeyMsg, node *ctree.Node) (ui.RouterMsg, bool) {
	_, result, ok := keybindings.IterMatch(p, actionKeyIterator(node))
	if !ok {
		return ui.EmptyRouterMsg(), false
	}

	return acStateToOpMsg(result.id, result.state), true
}

func acStateToOpMsg(id string, state *adActionState) ui.RouterMsg {
	if state == nil || id == "" {
		return ui.EmptyRouterMsg()
	}

	createInfo, opFunc := state.invokeAction(state)
	createInfo.UpdateID(id)

	return operations.MsgCreate(createInfo, opFunc)
}

// actionKeyIterResult holds the action-key iteration result.
type actionKeyIterResult struct {
	state *adActionState
	id    string
}

// actionKeyIterator returns an iterator that iterates over [nodeTypeAction] nodes.
func actionKeyIterator(node *ctree.Node) keybindings.IterKeyMatch[actionKeyIterResult] {
	return func(yield func(keybindings.Keybinding, actionKeyIterResult) bool) {
		n := node
		if n == nil {
			return
		}

		ndat := n.Data[*adNode]()
		pos, ok := ndat.noder.SubnodePosition(relPosActionsListNode)
		if !ok {
			return
		}

		ch := node.Children()
		if len(ch) == 0 {
			return
		}

		actionListNode := ch[pos]
		actionNodes := actionListNode.Children()

		for _, ac := range actionNodes {
			data := ac.Data[*adNode]()

			res := actionKeyIterResult{
				id:    data.ID(),
				state: data.actionState,
			}

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
