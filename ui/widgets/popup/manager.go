package popup

import (
	"iter"
	"slices"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/inert"
	"github.com/bluetuith-org/bluetuith/ui"
)

// Manager describes a popup manager.
type Manager struct {
	popups []Popup
}

// NewManager returns a new popup manager.
func NewManager() *Manager {
	return &Manager{}
}

// HandleKey sends all key events to the focused popup, if present.
func (p *Manager) HandleKey(msg tview.KeyMsg, rsp *ui.Response) bool {
	if len(p.popups) == 0 {
		return false
	}

	p.popups[len(p.popups)-1].Update(msg, rsp)

	return true
}

// HandleMsg handles popup related messages.
func (p *Manager) HandleMsg(msg Router, rsp *ui.Response) bool {
	switch ms := msg.(type) {
	case R[popupDeleteMsg]:
		p.pop()

	case R[popupRemoveMsg]:
		p.removeByID(ms.id)

	case R[popupModifyMsg]:
		p.add(ms.Value.popup)

	default:
		if !msg.IsValid() {
			return false
		}

		for _, p := range p.popups {
			if p.PopupID() == msg.PopupID() {
				p.Update(msg, rsp)
				break
			}
		}
	}

	return true
}

// HasPopups returns if popups exist or not.
func (p *Manager) HasPopups() bool {
	return len(p.popups) != 0
}

// Popups returns an iterator over the registered popups.
func (p *Manager) Popups() iter.Seq2[tview.Widget, bool] {
	return func(yield func(tview.Widget, bool) bool) {
		if len(p.popups) == 0 {
			return
		}

		focusedIdx := len(p.popups) - 1

		for idx := 0; idx < len(p.popups)-1; idx++ {
			if !yield(inert.New(p.popups[idx].View(false)), false) {
				return
			}
		}

		if !yield(p.popups[focusedIdx].View(true), true) {
			return
		}
	}
}

func (p *Manager) add(popup Popup) {
	idx := -1
	pid := popup.PopupID()

	for i, p := range p.popups {
		if p.PopupID() == pid {
			idx = i
			break
		}
	}

	if idx >= 0 {
		p.popups[idx] = popup
		return
	}

	p.popups = append(p.popups, popup)
}

func (p *Manager) pop() {
	p.popups = slices.Delete(p.popups, len(p.popups)-1, len(p.popups))
}

func (p *Manager) removeByID(id string) {
	idx := -1

	for i, popup := range p.popups {
		if popup.PopupID() == id {
			idx = i
			break
		}
	}

	if idx < 0 {
		return
	}

	p.popups = slices.Delete(p.popups, idx, idx+1)
}
