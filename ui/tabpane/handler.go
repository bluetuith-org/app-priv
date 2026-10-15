package tabpane

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/info"
	"github.com/bluetuith-org/bluetuith/ui/logs"
	"github.com/bluetuith-org/bluetuith/ui/operations"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
)

type viewHandler struct {
	tabs.State

	views []ui.ViewID

	infoModel       info.Model
	operationsModel operations.Model
	logsModel       logs.Model

	rv ui.RootView
}

type vm interface {
	ui.View
	tabs.TabSection
}

func newHandler(rv ui.RootView) *viewHandler {
	h := &viewHandler{
		State: tabs.NewState(),
		rv:    rv,
	}

	h.operationsModel = operations.New(rv)
	h.addModel(h.operationsModel)

	h.logsModel = logs.New(rv)
	h.addModel(h.logsModel)

	h.infoModel = info.New(rv)
	h.addModel(h.infoModel)

	return h
}

func (h *viewHandler) addModel[T vm](v T) {
	h.AddSection(v.Icon(), v.Title())
	h.views = append(h.views, v.ViewID())
}

func (h *viewHandler) init() tview.Cmd {
	return tview.Batch(
		h.infoModel.Init(),
		h.operationsModel.Init(),
		h.logsModel.Init(),
	)
}

func (h *viewHandler) update(id ui.ViewID, msg tview.Msg) tview.Cmd {
	var cmd tview.Cmd

	switch id {
	case ui.ViewIDInfo:
		h.infoModel, cmd = h.infoModel.Update(msg)

	case ui.ViewIDOperations:
		h.operationsModel, cmd = h.operationsModel.Update(msg)

	case ui.ViewIDLog:
		h.logsModel, cmd = h.logsModel.Update(msg)

	default:
	}

	return cmd
}

func (h *viewHandler) header() []richtext.Line {
	var lines []richtext.Line

	if op := h.operationsModel.Header(); op != nil {
		lines = append(lines, op)
	}

	return lines
}

func (h *viewHandler) status() richtext.Line {
	return h.logsModel.StatusBar()
}

func (h *viewHandler) activeID() ui.ViewID {
	return h.views[h.ActiveIndex()]
}

func (h *viewHandler) setActive(msg tview.KeyMsg) (tview.Cmd, bool) {
	var id ui.ViewID

	switch {
	case ui.Kb().Information.Show.Matches(msg):
		id = ui.ViewIDInfo

	case ui.Kb().Operations.Show.Matches(msg):
		id = ui.ViewIDOperations

	case ui.Kb().Transfers.Show.Matches(msg):
		id = ui.ViewIDTransfers

	case ui.Kb().Log.Show.Matches(msg):
		id = ui.ViewIDLog

	default:
		return nil, false
	}

	idx := -1
	for i, vid := range h.views {
		if vid == id {
			idx = i
			break
		}
	}

	if idx < 0 {
		return nil, false
	}

	h.SetActiveIndex(idx)
	h.rv.SetFocused(id)

	return nil, true
}

func (h *viewHandler) activeContent(focused bool) tview.Widget {
	id := h.views[h.ActiveIndex()]

	switch id {
	case ui.ViewIDInfo:
		return h.infoModel.View(focused)

	case ui.ViewIDOperations:
		return h.operationsModel.View(focused)

	case ui.ViewIDLog:
		return h.logsModel.View(focused)

	default:
	}

	return nil
}

func (h *viewHandler) updateRouted(id ui.ViewID, msg tview.Msg) tview.Cmd {
	return h.update(id, msg)
}

func (h *viewHandler) updateAll(msg tview.Msg) tview.Cmd {
	var cmds []tview.Cmd

	for _, id := range h.views {
		cmd := h.update(id, msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if cmds == nil {
		return nil
	}

	if len(cmds) == 1 {
		return cmds[0]
	}

	return tview.Batch(cmds...)
}
