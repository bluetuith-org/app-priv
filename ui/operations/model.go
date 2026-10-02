package operations

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/list"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
	"github.com/jba/omap/ordered"
)

// Model represents the operations view.
type Model struct {
	state list.SelectionState

	running *ordered.Map[string, *RunningInfo]
	focused bool

	rv ui.RootView
}

// New returns a new operations view.
func New(rv ui.RootView) Model {
	var state list.SelectionState

	state.SetCursor(0)

	return Model{
		state: state,
		running: ordered.NewMap[string, *RunningInfo](func(s1, s2 string) int {
			return strings.Compare(s1, s2)
		}),
		rv: rv,
	}
}

// Title returns the title of the tab section.
func (m Model) Title() string {
	return "Operations"
}

// Icon returns the icon associated with the tab section.
func (m Model) Icon() string {
	return theme.Icons().Operations.String()
}

// ViewID returns the view's ID.
func (m Model) ViewID() ui.ViewID {
	return ui.ViewIDOperations
}

// Keybinds implements [ui.KcVirtualHandler].
func (m Model) Keybinds(msg tview.KeyMsg) (list.Action, bool) {
	switch {
	case kb().NavigateUp.Matches(msg):
		return list.ActionSelectUp, true

	case kb().NavigateDown.Matches(msg):
		return list.ActionSelectDown, true
	}

	return 0, false
}

// OnChange implements [ui.KcVirtualHandler].
func (m Model) OnChange(change list.Change) tview.Msg {
	return listChg(change)
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m Model) Update(focused bool, msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(focused, msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(focused bool) tview.Element {
	ls := list.New(m.state, m.running.Len(), func(index int) tview.Element {
		_, info := m.running.Nth(index)
		if info == nil {
			return card.Widget{}
		}

		header, msg, status := info.buildSegments()
		style := theme.Current().Global

		borders := tview.BorderSetRound()
		borderStyle := style
		if m.state.Cursor() == index && focused {
			borders = tview.BorderSetThick()
			borderStyle = theme.Current().Operations.Selection
		}

		return card.New(info.description).
			Header(header.Text, header.Style).
			Status(status.Text, status.Style).
			Message(msg.Text, msg.Style).
			Background(style.GetBackground()).
			BorderStyle(borderStyle).
			BorderSet(borders)
	}).
		Focused(focused).
		OnChange(m.OnChange).
		Keybind(m.Keybinds)

	return ls
}

func (m *Model) updateModel(_ bool, msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case listChg:
		m.state.Apply(list.Change(ms))

	case CreateMsg:
		info, action, err := m.createNewOperation(ms)
		if err != nil {
			return nil
		}

		return tview.Sequence(func() tview.Msg {
			return action(info)
		}, m.deleteOperationCmd(info.id))

	case opUpdateMsg:
		m.processUpdateMsg(ms)

	case opDeleteMsg:
		//m.removeOperation(ms)

	case opErrorMsg:
	}

	return nil
}

func (m *Model) createNewOperation(msg CreateMsg) (*RunningInfo, Invoker, error) {
	if msg.id == "" {
		return nil, nil, fmt.Errorf("%w: msg ID empty on create", errOpInternal)
	}

	_, ok := m.running.Get(msg.id)
	if ok {
		return nil, nil, errOpAlreadyInProgress
	}

	header := "Operation #" + strconv.Itoa(m.running.Len()+1)

	info := newRunningInfo(m.rv, msg.CreationInfo.updateHeader(header))
	m.running.Set(msg.id, info)

	return info, msg.opAction, nil
}

func (m *Model) processUpdateMsg(msg opUpdateMsg) bool {
	if msg.id == "" {
		return false
	}

	info, ok := m.running.Get(msg.id)
	if !ok {
		return false
	}

	info.fromUpdateMsg(msg)

	return true
}

func (m *Model) removeOperation(msg opDeleteMsg) {
	m.running.Delete(msg.id)
}

func (m *Model) deleteOperationCmd(id string) tview.Cmd {
	return func() tview.Msg {
		time.Sleep(1 * time.Second)
		return msgOpDelete(id)
	}
}

type listChg list.Change

func kb() *keybindings.Keybindings {
	return keybindings.Current
}

var (
	_ ui.Model[Model] = Model{}
	_ tabs.TabSection = Model{}
)

var (
	errOpAlreadyInProgress = errors.New("this operation is already in progress")
	errOpInternal          = errors.New("(op) an internal error has occurred")
)
