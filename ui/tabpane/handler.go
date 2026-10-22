package tabpane

import (
	"github.com/ayn2op/tview"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/info"
	"github.com/bluetuith-org/bluetuith/ui/logs"
	"github.com/bluetuith-org/bluetuith/ui/operations"
)

type viewHandler struct {
	activeIndex int

	views  []vm
	labels []string

	rv ui.RootView
}

type vm interface {
	ui.Model
	ui.TabSection
}

func newHandler(rv ui.RootView) *viewHandler {
	h := &viewHandler{
		rv: rv,
	}

	h.addModel(operations.New(rv))
	h.addModel(logs.New(rv))
	h.addModel(info.New(rv))

	return h
}

func (h *viewHandler) addModel[T vm](v T) {
	h.views = append(h.views, v)

	icon, title := v.Icon(), v.Title()
	h.labels = append(h.labels, icon+" "+title)
}

func (h *viewHandler) init() tview.Cmd {
	r := ui.NewResponse()
	for _, views := range h.views {
		r.PushCmd(views.Init())
	}

	return r.Cmds()
}

func (h *viewHandler) update(id ui.ViewID, msg tview.Msg, rsp *ui.Response) {
	for _, v := range h.views {
		if v.ViewID() != id {
			continue
		}

		v.Update(msg, rsp)
		return
	}
}

func (h *viewHandler) active() vm {
	return h.views[h.activeIndex]
}

func (h *viewHandler) setActiveIndex(idx int) {
	h.activeIndex = idx
}

func (h *viewHandler) setActive(msg tview.KeyMsg) bool {
	var id ui.ViewID

	idx := -1
	for i, view := range h.views {
		if !view.ViewKeybind().Matches(msg) {
			continue
		}

		idx = i
		id = view.ViewID()

		break
	}
	if idx < 0 {
		return false
	}

	h.activeIndex = idx
	h.rv.SetFocused(id)

	return true
}

func (h *viewHandler) updateRouted(id ui.ViewID, msg tview.Msg, rsp *ui.Response) {
	h.update(id, msg, rsp)
}

func (h *viewHandler) updateAll(msg tview.Msg, rsp *ui.Response) {
	for _, view := range h.views {
		view.Update(msg, rsp)
	}
}
