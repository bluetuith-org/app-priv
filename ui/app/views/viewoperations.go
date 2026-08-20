package views

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
)

var (
	errOpNoneInProgress    = errors.New("no operations in progress")
	errOpAlreadyInProgress = errors.New("this operation is already in progress")
	errOpInternal          = errors.New("(op) an internal error has occurred")
)

type operationsView struct {
	sync.Mutex

	width, height int
	focused       bool

	mapIDToIndex map[string]int
	orderedList  []*opRunningInfo

	v rootView
}

// InitializeView initializes the view.
func (o *operationsView) InitializeView(_ *appfeatures.FeatureSet) (inited bool, err error) {
	o.mapIDToIndex = make(map[string]int)
	o.orderedList = make([]*opRunningInfo, 0, 10)

	return true, nil
}

// SetRootView sets the root view.
func (o *operationsView) SetRootView(v rootView) {
	o.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (o *operationsView) AttachToTabView() (tabSection, bool) {
	return o, true
}

// Resize resizes the view.
func (o *operationsView) Resize(width int, height int) tea.WindowSizeMsg {
	o.width, o.height = width, height

	return tea.WindowSizeMsg{Width: o.width, Height: o.height}
}

// SetFocus sets whether the view is currently focused.
func (o *operationsView) SetFocus(focused bool) {
	o.focused = focused
}

// GetFocus gets whether the view is currently focused.
func (o *operationsView) GetFocus() bool {
	return o.focused
}

// UpdateStyles updates the styles for the view.
func (o *operationsView) UpdateStyles() {
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (o *operationsView) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (o *operationsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		msg = o.Resize(m.Width, m.Height)

	case tea.KeyPressMsg:
		if !o.focused {
			return o, nil
		}

		switch {
		case keybindings.MatchesKey(keybindings.KeySelect, m):
		default:
		}

	case opCreateMsg:
		opinfo, opAction, err := o.createNewOperation(m)
		if err != nil {
			return o, nil
		}

		cmd = tea.Sequence(func() tea.Msg {
			return opAction(opinfo)
		}, o.deleteOperationCmd(opinfo.id))

	case opUpdateMsg:
		o.processUpdateMsg(m)

	case opDeleteMsg:
		o.removeOperation(m)

	case opErrorMsg:
		// TODO: Handle errors.
		_ = m
	}

	return o, cmd
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (o *operationsView) View() tea.View {
	mv := lipgloss.NewStyle().
		Width(o.width).
		Height(o.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	text := ""
	if len(o.orderedList) > 0 {
		text = useBuffer(len(o.orderedList)+10, func(b *strings.Builder) {
			for _, op := range o.orderedList {
				b.WriteString(op.id)
				b.WriteString(": \n")
				b.WriteString(op.message)
				b.WriteString("\n\n")
			}
		})
	}

	return tea.NewView(mv.Render(text))
}

func (o *operationsView) Title() string {
	return "Operations"
}

func (o *operationsView) Icon() *iconVariant {
	return o.v.Icons().CogWheel
}

func (o *operationsView) createNewOperation(msg opCreateMsg) (*opRunningInfo, opInvoker, error) {
	if msg.id == "" {
		return nil, nil, fmt.Errorf("%w: msg ID empty on create", errOpInternal)
	}

	if len(o.mapIDToIndex) == 0 {
		return nil, nil, errOpNoneInProgress
	}

	_, ok := o.mapIDToIndex[msg.id]
	if ok {
		return nil, nil, errOpAlreadyInProgress
	}

	info := newOpRunningInfo(o.v, msg.id, msg.description, msg.message)

	o.orderedList = append(o.orderedList, info)
	o.mapIDToIndex[info.id] = len(o.orderedList)

	return info, msg.opAction, nil
}

func (o *operationsView) processUpdateMsg(msg opUpdateMsg) bool {
	if msg.id == "" {
		return false
	}

	info, ok := o.findInfoByID(msg.id)
	if !ok {
		return false
	}

	if msg.description != "" {
		info.description = msg.description
	}

	if msg.message != "" {
		info.message = msg.message
	}

	return true
}

func (o *operationsView) removeOperation(msg opDeleteMsg) {
	idx, ok := o.mapIDToIndex[msg.id]
	if !ok {
		return
	}

	o.orderedList = slices.Delete(o.orderedList, idx, idx+1)
}

func (o *operationsView) deleteOperationCmd(id string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return newOpDeleteMsg(id)
	}
}

func (o *operationsView) findInfoByID(id string) (*opRunningInfo, bool) {
	idx, ok := o.mapIDToIndex[id]
	if !ok {
		return nil, false
	}

	if idx < 0 || idx > len(o.orderedList) {
		return nil, false
	}

	return o.orderedList[idx], true
}

type opInvoker func(ov *opRunningInfo) tea.Msg

type opCreationInfo struct {
	id string

	description string
	message     string
}

func newOpCreationInfo(id, desc, msg string) opCreationInfo {
	return opCreationInfo{id, desc, msg}
}

type opRunningInfo struct {
	opCreationInfo

	v rootView
}

func newOpRunningInfo(v rootView, id, desc, msg string) *opRunningInfo {
	return &opRunningInfo{opCreationInfo: newOpCreationInfo(id, desc, msg), v: v}
}

func (o *opRunningInfo) info(msg string) {
	o.v.SendMsg(newOpUpdateMsg(o.id, "", msg))
}

func (o *opRunningInfo) updateDescription(desc string) {
	o.v.SendMsg(newOpUpdateMsg(o.id, desc, ""))
}

func (o *opRunningInfo) opSuccess(state actionStateSpec) actionUpdateMsg {
	return newActionUpdateMsg(o.id, state)
}

func (o *opRunningInfo) opError(err error) opErrorMsg {
	return newOpErrorMsg(err)
}

type opCreateMsg struct {
	opCreationInfo

	opAction opInvoker
}

func newOpCreateMsg(creationInfo opCreationInfo, action opInvoker) opCreateMsg {
	return opCreateMsg{opCreationInfo: creationInfo, opAction: action}
}

type opUpdateMsg struct {
	opCreationInfo
}

func newOpUpdateMsg(id, desc, msg string) opUpdateMsg {
	return opUpdateMsg{opCreationInfo: newOpCreationInfo(id, desc, msg)}
}

type opDeleteMsg struct {
	id string
}

func newOpDeleteMsg(id string) opDeleteMsg {
	return opDeleteMsg{id}
}

type opErrorMsg struct {
	err error
}

func newOpErrorMsg(err error) opErrorMsg {
	return opErrorMsg{err}
}
