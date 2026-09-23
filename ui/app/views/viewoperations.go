package views

import "github.com/ayn2op/tview"

type opInvoker func(ov *opRunningInfo) tview.Msg

type opCreationInfo struct {
	id string

	description string
	message     string
}

func newOpCreationInfo(desc, msg string) opCreationInfo {
	return opCreationInfo{"", desc, msg}
}

func (o *opCreationInfo) updateID(id string) {
	o.id = id
}

type opRunningInfo struct {
	opCreationInfo

	v rootView
}

func newOpRunningInfo(v rootView, creationInfo opCreationInfo) *opRunningInfo {
	return &opRunningInfo{opCreationInfo: creationInfo, v: v}
}

func (o *opRunningInfo) info(msg string) {
	o.v.SendMsg(msgOpUpdate(o.opCreationInfo, "", msg))
}

func (o *opRunningInfo) updateDescription(desc string) {
	o.v.SendMsg(msgOpUpdate(o.opCreationInfo, desc, ""))
}

func (o *opRunningInfo) opSuccess(state actionStateSpec) routerMsg {
	return msgAdActionUpdate(o.id, state)
}

func (o *opRunningInfo) opError(err error) routerMsg {
	return msgOpError(err)
}

type opCreateMsg struct {
	opCreationInfo

	opAction opInvoker
}

func msgOpCreate(creationInfo opCreationInfo, action opInvoker) routerMsg {
	return viewIDOperations.routerMessage(opCreateMsg{opCreationInfo: creationInfo, opAction: action})
}

type opUpdateMsg struct {
	opCreationInfo
}

func msgOpUpdate(creationInfo opCreationInfo, desc, msg string) routerMsg {
	creationInfo.description = desc
	creationInfo.message = msg

	return viewIDOperations.routerMessage(opUpdateMsg{opCreationInfo: creationInfo})
}

type opDeleteMsg struct {
	id string
}

func msgOpDelete(id string) routerMsg {
	return viewIDOperations.routerMessage(opDeleteMsg{id})
}

type opErrorMsg struct {
	err error
}

func msgOpError(err error) routerMsg {
	return viewIDOperations.routerMessage(opErrorMsg{err})
}
