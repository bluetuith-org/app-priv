package views

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/tree"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
)

type adActionState struct {
	key keybindings.Keybinding

	currentState actionStateSpec
	isToggleable bool

	invokeAction actionInvoker
}

type actionInvoker func() (opCreationInfo, opInvoker)

func newAdActionState(key keybindings.Keybinding, state actionStateSpec, isToggleable bool, invoker actionInvoker) *adActionState {
	return &adActionState{key, state, isToggleable, invoker}
}

type actionStateSpec uint8

const (
	actionStateNone actionStateSpec = iota
	actionStateEnabled
	actionStateDisabled
)

type actionUpdateMsg struct {
	id        string
	stateSpec actionStateSpec
}

func emptyActionUpdateMsg() actionUpdateMsg {
	return actionUpdateMsg{}
}

func (a *actionUpdateMsg) updateState(stateID string, ad *adActionState) bool {
	if !a.isValid(stateID) {
		return false
	}

	ad.currentState = a.stateSpec
	return true
}

func (a *actionUpdateMsg) isValid(stateID string) bool {
	return a.id != "" && a.id == stateID
}

func actionHandleKeyMsg(p tview.KeyMsg, node *adTreeNode) (routerMsg, bool) {
	_, result, ok := keybindings.IterMatch(p, actionKeyIterator(node))
	if !ok {
		return emptyRouterMsg(), false
	}

	return acStateToOpMsg(result.id, result.state), true
}

func acStateToOpMsg(id string, state *adActionState) routerMsg {
	if state == nil || id == "" {
		return emptyRouterMsg()
	}

	createInfo, opFunc := state.invokeAction()
	createInfo.updateID(id)

	return msgOpCreate(createInfo, opFunc)
}

type actionKeyIterResult struct {
	state *adActionState
	id    string
}

func actionKeyIterator(node *adTreeNode) keybindings.IterKeyMatch[*actionKeyIterResult] {
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
			data := getNodeReference(ac)

			res.id = data.id.String()
			res.state = data.actionState

			if !yield(data.actionState.key, res) {
				return
			}
		}
	}
}

func boolToActionState(val bool) actionStateSpec {
	if !val {
		return actionStateDisabled
	}

	return actionStateEnabled
}

//revive:disable
func getNodeReference(n *tree.Node) *adTreeNode {
	return n.Reference().(*adTreeNode)
}
