package views

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Digital-Shane/treeview/v2"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
)

type adActionState struct {
	key keybindings.KeyID

	currentState actionStateSpec
	isToggleable bool

	invokeAction actionInvoker
}

type actionInvoker func() (opCreationInfo, opInvoker)

func newAdActionState(key keybindings.KeyID, state actionStateSpec, isToggleable bool, invoker actionInvoker) *adActionState {
	return &adActionState{key, state, isToggleable, invoker}
}

type actionStateSpec uint8

const (
	actionStateNone actionStateSpec = iota
	actionStateEnabled
	actionStateDisabled
)

type actionUpdateMsg struct {
	id    string
	state actionStateSpec
}

func newActionUpdateMsg(id string, state actionStateSpec) actionUpdateMsg {
	return actionUpdateMsg{id, state}
}

func emptyActionUpdateMsg() actionUpdateMsg {
	return actionUpdateMsg{}
}

func (a *actionUpdateMsg) updateState(stateID string, ad *adActionState) bool {
	if !a.isValid(stateID) {
		return false
	}

	ad.currentState = a.state
	return true
}

func (a *actionUpdateMsg) isValid(stateID string) bool {
	return a.id != "" && a.id == stateID
}

func actionHandleKeyMsg(p tea.KeyPressMsg, node *treeview.Node[adTreeNode]) (tea.Msg, bool) {
	_, actionState, ok := keybindings.IterMatch(p, actionKeyIterator(node))
	if !ok {
		return nil, false
	}

	createInfo, opFunc := actionState.invokeAction()

	return newOpCreateMsg(createInfo, opFunc), true
}

func actionKeyIterator(node *treeview.Node[adTreeNode]) keybindings.IterKeyMatch[*adActionState] {
	return func(yield func(keybindings.KeyID, *adActionState) bool) {
		n := node
		if n == nil {
			return
		}

		ch := node.Children()
		if len(ch) <= int(actionsListNodePos) {
			return
		}

		actionListNode := ch[actionsListNodePos]
		actionNodes := actionListNode.Children()

		for _, ac := range actionNodes {
			data := ac.Data()
			if !yield(data.actionState.key, data.actionState) {
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
