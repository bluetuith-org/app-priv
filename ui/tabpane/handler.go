package tabpane

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/info"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
)

type viewHandler struct {
	*tabs.State

	infoModel info.Model

	indexView map[int]ui.ViewID
	index     int
}

type vm interface {
	ui.View
	tabs.TabSection
}

func newHandler() *viewHandler {
	h := &viewHandler{
		State:     tabs.NewState(),
		indexView: make(map[int]ui.ViewID),
	}

	h.infoModel = info.New()
	h.addModel(h.infoModel)

	return h
}

func (h *viewHandler) addModel[T vm](v T) {
	h.AddSection(v.Icon(), v.Title())

	h.indexView[h.index] = v.ViewID()
	h.index++
}

func (h *viewHandler) init() tview.Cmd {
	return tview.Batch(h.infoModel.Init())
}

func (h *viewHandler) update(id ui.ViewID, focused bool, msg tview.Msg) tview.Cmd {
	var cmd tview.Cmd

	switch id {
	case ui.ViewIDInfo:
		h.infoModel, cmd = h.infoModel.Update(focused, msg)

	default:
	}

	return cmd
}

func (h *viewHandler) activeContent(focused bool) tview.Element {
	id := h.indexView[h.ActiveIndex()]

	switch id {
	case ui.ViewIDInfo:
		return h.infoModel.View(focused)

	default:
	}

	return nil
}

func (h *viewHandler) updateRouted(id ui.ViewID, focused bool, msg tview.Msg) tview.Cmd {
	focused = focused && id == h.indexView[h.ActiveIndex()]

	return h.update(id, focused, msg)
}

func (h *viewHandler) updateAll(focused bool, msg tview.Msg) tview.Cmd {
	var cmds []tview.Cmd

	activeID := h.indexView[h.ActiveIndex()]

	for _, id := range h.indexView {
		cmd := h.update(id, focused && activeID == id, msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if cmds == nil {
		return nil
	}

	return tview.Batch(cmds...)
}
