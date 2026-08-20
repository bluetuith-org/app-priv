package views

import tea "charm.land/bubbletea/v2"

type viewID int

const (
	viewIDNone viewID = iota
	viewIDAdTree
	viewIDTabs
	viewIDInfo
	viewIDOperations
	_viewIDMax
)

const maxViews = _viewIDMax - 1

func (v viewID) isValid() bool {
	return v != viewIDNone
}

func (v viewID) routerMessage(t tea.Msg) routerMsg {
	return newRouterMsg(v, t)
}

type routerMsg struct {
	id  viewID
	msg tea.Msg
}

func newRouterMsg(id viewID, msg tea.Msg) routerMsg {
	return routerMsg{id, msg}
}

func handleRouterMsg(v viewer, msg routerMsg) tea.Cmd {
	_, cmd := v.Update(msg.msg)
	return cmd
}
