package operations

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/list"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
	"github.com/gdamore/tcell/v3"
)

// Model represents the operations view.
type Model struct {
	state list.SelectionState

	focused bool

	run map[string]*RunningInfo
	ord []*RunningInfo

	count int

	rv ui.RootView
}

// New returns a new operations view.
func New(rv ui.RootView) Model {
	var state list.SelectionState

	state.SetCursor(0)

	return Model{
		state: state,
		run:   make(map[string]*RunningInfo),
		rv:    rv,
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
func (m Model) Keybinds(msg tview.KeyMsg) list.Action {
	switch {
	case kb().NavigateUp.Matches(msg):
		return list.ActionSelectUp

	case kb().NavigateDown.Matches(msg):
		return list.ActionSelectDown
	}

	return list.ActionNone
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
func (m Model) View(focused bool) tview.Widget {
	ls := list.New(m.state, len(m.ord), func(index int) tview.Widget {
		info := m.ord[index]
		style := theme.Current().Global

		borders := tview.BorderSetRound()
		borderStyle := style
		if m.state.Cursor() == index && focused {
			borders = tview.BorderSetThick()
			borderStyle = theme.Current().Operations.Selection
		}

		return card.New(info.cc).
			Background(style.GetBackground()).
			BorderStyle(borderStyle).
			BorderSet(borders)
	}).
		Focused(focused).
		OnChange(m.OnChange).
		Keybind(m.Keybinds).
		SelectedStyle(tcell.StyleDefault)

	return ls
}

func (m *Model) updateModel(focused bool, msg tview.Msg) tview.Cmd {
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

	case tview.KeyMsg:
		if !focused {
			return nil
		}

		switch {
		case kb().Operations.ClearAll.Matches(ms):
			m.clearOperations()

		case kb().Operations.Cancel.Matches(ms):
		}
	}

	return nil
}

func (m *Model) createNewOperation(msg CreateMsg) (*RunningInfo, Invoker, error) {
	if msg.id == "" {
		return nil, nil, fmt.Errorf("%w: msg ID empty on create", errOpInternal)
	}

	runInfo, ok := m.run[msg.id]
	if ok && runInfo.stage == operationInProgress {
		return nil, nil, errOpAlreadyInProgress
	}

	m.count++
	header := "#" + strconv.Itoa(m.count)

	info := newRunningInfo(m.rv, msg.CreationInfo.updateHeader(header))
	for range 10 {
		m.run[info.id] = info
		m.ord = slices.Insert(m.ord, 0, info)
	}

	return info, msg.opAction, nil
}

func (m *Model) processUpdateMsg(msg opUpdateMsg) bool {
	if msg.id == "" {
		return false
	}

	info, ok := m.run[msg.id]
	if !ok {
		return false
	}

	info.fromUpdateMsg(msg)

	return true
}

func (m *Model) clearOperations() {
	clear(m.run)
	m.ord = nil
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
