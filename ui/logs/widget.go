package logs

import (
	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/box"
	"github.com/ayn2op/tview/list"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/gdamore/tcell/v3"
)

type widget struct {
	list.Widget

	h *handler

	cursor              int
	selectItem, focused bool
}

func newWidget(h *handler, cursor int, selectItem bool) widget {
	w := widget{h: h, focused: selectItem}

	w.Widget = list.New(cursor, w.h.count(), w.Items).
		SelectedStyle(tcell.StyleDefault).
		Focused(w.focused).
		OnChange(w.OnChange).
		Keybind(w.ListKeybinds)

	return w
}

// Draw draws the widget onto the screen within the given area.
func (w widget) Draw(screen tview.Screen, area tview.Rectangle) {
	w.h.Lock()
	defer w.h.Unlock()

	w.Widget.Draw(screen, area)
}

func (w widget) Items(index int) tview.Widget {
	style := theme.Current().Global

	bx := box.New(w.h.rb.Get(index)).
		Background(style.GetBackground())

	if w.selectItem && w.cursor == index {
		bx = bx.Borders(tview.BordersAll).
			BorderSet(tview.BorderSetRound()).
			BorderStyle(theme.Current().Log.Selection)
	}

	return bx
}

// ListKeybinds implements keybindings for the list.
func (w widget) ListKeybinds(msg tview.KeyMsg) list.Action {
	switch {
	case ui.Kb().NavigateUp.Matches(msg):
		return list.ActionSelectUp

	case ui.Kb().NavigateDown.Matches(msg):
		return list.ActionSelectDown
	}

	return list.ActionNone
}

// OnChange implements a handler that returns a message when the list's state is changed.
func (w widget) OnChange(index int) tview.Msg {
	return ui.ViewIDLog.RouterMessage(listChg(index))
}
