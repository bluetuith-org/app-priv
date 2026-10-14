package operations

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/list"
	"github.com/ayn2op/tview/richtext"
	"github.com/ayn2op/tview/scrollbar"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
	"github.com/bluetuith-org/bluetuith/ui/widgets/tabs"
	"github.com/gdamore/tcell/v3"
)

// Model represents the operations view.
type Model struct {
	scrollState viewport.ScrollState
	scrollBar   scrollbar.Widget

	count int

	cursor, target int

	run map[string]*RunningInfo
	ord []*RunningInfo

	rv ui.RootView
}

// New returns a new operations view.
func New(rv ui.RootView) Model {
	return Model{
		cursor:    -1,
		target:    -1,
		scrollBar: scrollbar.New(),
		run:       make(map[string]*RunningInfo),
		rv:        rv,
	}
}

// Title returns the title of the tab section.
func (m Model) Title() string {
	return fmt.Sprintf("%s (%s)", m.ViewID().String(), ui.Kb().Operations.Show.KeyName)
}

// Icon returns the icon associated with the tab section.
func (m Model) Icon() string {
	return theme.Icons().Operations.String()
}

// ViewID returns the view's ID.
func (m Model) ViewID() ui.ViewID {
	return ui.ViewIDOperations
}

// Header returns the content to print in the title bar.
func (m *Model) Header() richtext.Line {
	if len(m.run) == 0 {
		return nil
	}

	style := theme.Current().Operations.TitleBarIndicator

	return richtext.NewLine(
		richtext.NewSegment(" ", style),
		richtext.NewSegment(strconv.Itoa(len(m.run)), style),
		richtext.NewSegment(" Running ", style),
	)
}

// Keybinds implements keybindings for the list.
func (m Model) Keybinds(msg tview.KeyMsg) list.Action {
	switch {
	case ui.Kb().NavigateUp.Matches(msg):
		return list.ActionSelectUp

	case ui.Kb().NavigateDown.Matches(msg):
		return list.ActionSelectDown
	}

	return list.ActionNone
}

// ScrollKeybinds implements keybindings for the viewport.
func (m Model) ScrollKeybinds(msg tview.KeyMsg) viewport.Action {
	switch {
	case ui.Kb().NavigateTop.Matches(msg):
		return viewport.ActionTop

	case ui.Kb().NavigateBottom.Matches(msg):
		return viewport.ActionBottom
	}

	return viewport.ActionNone
}

// OnChange implements [ui.KcVirtualHandler].
func (m Model) OnChange(cursor int) tview.Msg {
	return listChg(cursor)
}

// Init returns a command to run when the model starts, or nil.
func (m Model) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m Model) Update(msg tview.Msg) (Model, tview.Cmd) {
	cmd := m.updateModel(msg)
	return m, cmd
}

// View draws this model onto the screen.
func (m Model) View(focused bool) tview.Widget {
	ls := list.New(m.cursor, len(m.ord), func(index int) tview.Widget {
		info := m.ord[index]
		style := theme.Current().Global

		borders := tview.BorderSetRound()
		borderStyle := style
		if m.cursor == index && focused {
			borderStyle = theme.Current().Operations.Selection
		}

		return card.New(info.cc).
			Borders(tview.BordersAll).
			BorderStyle(borderStyle).
			BorderSet(borders).
			Background(style.GetBackground())
	}).
		Focused(focused).
		OnChange(m.OnChange).
		Keybind(m.Keybinds).
		SelectedStyle(tcell.StyleDefault)

	vp := viewport.New(ls, m.scrollState).
		Target(ls.Target(m.target)).
		ScrollBar(m.scrollBar, viewport.ScrollBarVisibilityAutomatic).
		OnChildMsg(func(c *viewport.ScrollState, msg tview.Msg) tview.Msg { return scrollMsg{change: *c, msg: msg} }).
		Axes(layout.Axes{Width: true, Height: true}).
		Keybind(m.ScrollKeybinds).
		Focused(focused).
		OnChange(func(c *viewport.ScrollState) tview.Msg { return scrollMsg{change: *c} })

	return vp
}

func (m *Model) setCursor(index int) {
	if index >= 0 && index != m.cursor {
		m.scrollState.ScrollToTarget()
		m.target = index
	}

	m.cursor = max(index, -1)
}

func (m *Model) updateModel(msg tview.Msg) tview.Cmd {
	switch ms := msg.(type) {
	case listChg:
		m.setCursor(int(ms))
		m.rv.SetFocused(m.ViewID())

	case scrollMsg:
		m.scrollState = viewport.ScrollState(ms.change)
		if ms.msg != nil {
			return m.updateModel(ms.msg)
		}

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
		m.deleteOperation(ms.id)

	case tview.KeyMsg:

		switch {
		case ui.Kb().Operations.ClearAll.Matches(ms):
			m.clearOperations()

		case ui.Kb().Operations.Cancel.Matches(ms):
		}
	}

	return nil
}

func (m *Model) createNewOperation(msg CreateMsg) (*RunningInfo, Invoker, error) {
	if msg.getID() == "" {
		return nil, nil, fmt.Errorf("%w: msg ID empty on create", errOpInternal)
	}

	runInfo, ok := m.run[msg.getID()]
	if ok && runInfo.stage == operationInProgress {
		return nil, nil, errOpAlreadyInProgress
	}

	m.count++
	header := "#" + strconv.Itoa(m.count)

	info := newRunningInfo(m.rv, msg.updateHeader(header))
	m.run[info.id] = info
	m.ord = slices.Insert(m.ord, 0, info)

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

func (m *Model) deleteOperation(id string) {
	delete(m.run, id)
}

func (m *Model) clearOperations() {
	clear(m.run)
	m.ord = nil
}

func (m *Model) deleteOperationCmd(id string) tview.Cmd {
	return func() tview.Msg {
		return msgOpDelete(id)
	}
}

type (
	listChg int

	scrollMsg struct {
		change viewport.ScrollState
		msg    tview.Msg
	}
)

var (
	_ ui.Model[Model] = Model{}
	_ tabs.TabSection = Model{}
)

var (
	errOpAlreadyInProgress = errors.New("this operation is already in progress")
	errOpInternal          = errors.New("(op) an internal error has occurred")
)
