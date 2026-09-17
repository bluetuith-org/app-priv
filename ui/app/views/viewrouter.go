package views

import (
	"github.com/ayn2op/tview"
)

type viewID int

const (
	viewIDNone viewID = iota
	viewIDAdTree
	viewIDTabs
	viewIDInfo
	viewIDOperations
	viewIDStatusBar
	viewIDLog
	_viewIDMax
)

const maxViews = _viewIDMax - 1

func (v viewID) isValid() bool {
	return v != viewIDNone
}

func (v viewID) routerMessage(t tview.Msg) routerMsg {
	return newRouterMsg(v, t)
}

type routerMsg struct {
	id  viewID
	msg tview.Msg
}

func newRouterMsg(id viewID, msg tview.Msg) routerMsg {
	return routerMsg{id, msg}
}

func emptyRouterMsg() routerMsg {
	return routerMsg{}
}

func (r *routerMsg) isValid() bool {
	return r.id != viewIDNone
}

func (r routerMsg) sendRoutedMsg(rv rootView) tview.Cmd {
	return rv.SendRoutedUpdateMsg(r)
}

func handleRouterMsg(v view, msg routerMsg) tview.Cmd {
	return v.Update(msg.msg)
}
