package operations

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/layout"
	"github.com/ayn2op/tview/list"
	"github.com/ayn2op/tview/scrollbar"
	"github.com/ayn2op/tview/text"
	"github.com/ayn2op/tview/viewport"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/card"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
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
func New(rv ui.RootView) *Model {
	return &Model{
		cursor:    0,
		target:    0,
		scrollBar: scrollbar.New(),
		run:       make(map[string]*RunningInfo),
		rv:        rv,
	}
}

// Title returns the title of the tab section.
func (m *Model) Title() string {
	return fmt.Sprintf("%s (%s)", m.ViewID().String(), ui.Kb().Operations.Show.KeyName)
}

// Icon returns the icon associated with the tab section.
func (m *Model) Icon() string {
	return theme.Icons().Operations.String()
}

// ViewID returns the view's ID.
func (m *Model) ViewID() ui.ViewID {
	return ui.ViewIDOperations
}

// Header returns the content to print in the title bar.
func (m *Model) Header() []text.Span {
	if len(m.run) == 0 {
		return nil
	}

	style := theme.Current().Operations.TitleBarIndicator

	return stext.NewLine(
		stext.NewSpan(" ", style),
		stext.NewSpan(strconv.Itoa(len(m.run)), style),
		stext.NewSpan(" Running ", style),
	)
}

// Status returns the content to print in the status bar area.
func (m *Model) Status() []text.Span {
	return nil
}

// ViewKeybind returns the keybinding used to focus the view.
func (m *Model) ViewKeybind() keybindings.Keybinding {
	return ui.Kb().Operations.Show
}

// ListKeybinds implements keybindings for the list.
func (m *Model) ListKeybinds(msg tview.KeyMsg) list.Action {
	switch {
	case ui.Kb().NavigateUp.Matches(msg):
		return list.ActionSelectUp

	case ui.Kb().NavigateDown.Matches(msg):
		return list.ActionSelectDown
	}

	return list.ActionNone
}

// ViewportKeybinds implements keybindings for the viewport.
func (m *Model) ViewportKeybinds(msg tview.KeyMsg) viewport.Action {
	switch {
	case ui.Kb().NavigateTop.Matches(msg):
		return viewport.ActionTop

	case ui.Kb().NavigateBottom.Matches(msg):
		return viewport.ActionBottom
	}

	return viewport.ActionNone
}

// OnListChange processes a change in the list and returns a message.
func (m *Model) OnListChange(cursor int) tview.Msg {
	return m.ViewID().RouterMessage(listChg(cursor))
}

// OnViewportChange processes a change in the viewport and returns a message.
func (m *Model) OnViewportChange(c *viewport.ScrollState) tview.Msg {
	return m.ViewID().RouterMessage(scrollMsg{change: *c})
}

// OnViewportChildMsg processes a message that was passed to its underlying view in the viewport and returns a message.
func (m *Model) OnViewportChildMsg(c *viewport.ScrollState, msg tview.Msg) tview.Msg {
	return m.ViewID().RouterMessage(scrollMsg{change: *c, msg: msg})
}

// Init returns a command to run when the model starts, or nil.
func (m *Model) Init() tview.Cmd {
	return nil
}

// Update receives messages when this model has focus.
func (m *Model) Update(msg tview.Msg, rsp *ui.Response) {
	switch ms := msg.(type) {
	case ui.R[listChg]:
		m.setCursor(int(ms.Value))
		m.rv.SetFocused(m.ViewID())

	case ui.R[scrollMsg]:
		m.scrollState = viewport.ScrollState(ms.Value.change)
		if ms.Value.msg != nil {
			m.Update(ms.Value.msg, rsp)
		}

	case ui.R[CreateMsg]:
		info, action, err := m.createNewOperation(ms.Value)
		if err != nil {
			return
		}

		rsp.PushCmd(func() tview.Msg {
			defer m.deleteOperationCmd(info.id)
			return action(info)
		})

	case ui.R[opUpdateMsg]:
		m.processUpdateMsg(ms.Value)

	case ui.R[opDeleteMsg]:
		m.deleteOperation(ms.Value.id)

	case tview.KeyMsg:
		switch {
		case ui.Kb().Operations.ClearAll.Matches(ms):
			m.clearOperations()

		case ui.Kb().Operations.Cancel.Matches(ms):
		}
	}
}

// View draws this model onto the screen.
func (m *Model) View(focused bool) tview.Widget {
	ls := list.New(m.cursor, len(m.ord), func(index int) tview.Widget {
		info := m.ord[index]
		style := theme.Current().Global

		borders := tview.BorderSetRound()
		borderStyle := style
		if m.cursor == index && focused {
			borderStyle = theme.Current().Operations.Selection
			borders = tview.BorderSetThick()
		}

		return card.New(info.cc).
			Borders(tview.BordersAll).
			BorderStyle(borderStyle).
			BorderSet(borders).
			Background(style.GetBackground())
	}).
		Focused(focused).
		SelectedStyle(tcell.StyleDefault).
		OnChange(m.OnListChange).
		Keybind(m.ListKeybinds)

	vp := viewport.New(ls, m.scrollState).
		Target(ls.Target(m.target)).
		ScrollBar(m.scrollBar, viewport.ScrollBarVisibilityAutomatic).
		Axes(layout.Axes{Width: true, Height: true}).
		Focused(focused).
		Keybind(m.ViewportKeybinds).
		OnChange(m.OnViewportChange).
		OnChildMsg(m.OnViewportChildMsg)

	return vp
}

func (m *Model) setCursor(index int) {
	if index >= 0 {
		m.scrollState.ScrollToTarget()
		m.target = index
	}

	m.cursor = max(index, -1)
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
	_ ui.Model      = (*Model)(nil)
	_ ui.TabSection = (*Model)(nil)
)

var (
	errOpAlreadyInProgress = errors.New("this operation is already in progress")
	errOpInternal          = errors.New("(op) an internal error has occurred")
)
