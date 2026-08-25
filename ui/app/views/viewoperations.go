package views

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
)

var (
	errOpAlreadyInProgress = errors.New("this operation is already in progress")
	errOpInternal          = errors.New("(op) an internal error has occurred")
)

type operationsView struct {
	width, height int
	focused       bool

	mapIDToIndex map[string]int
	orderedList  []*opRunningInfo

	vp viewport.Model

	v rootView
}

// ViewID returns the view's ID.
func (o *operationsView) ViewID() viewID {
	return viewIDOperations
}

// InitializeView initializes the view.
func (o *operationsView) InitializeView(_ *appfeatures.FeatureSet) (inited bool, err error) {
	o.mapIDToIndex = make(map[string]int)
	o.orderedList = make([]*opRunningInfo, 0, 10)

	o.vp = viewport.New()

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
		logError(m.err)
	}

	return o, cmd
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (o *operationsView) View() tea.View {
	text := ""
	if len(o.orderedList) > 0 {
		text = useStringBuffer(len(o.orderedList)+10, func(b *strings.Builder) {
			for _, op := range o.orderedList {
				b.WriteString(op.id)
				b.WriteString(": \n")
				b.WriteString(op.message)
				b.WriteString("\n\n")
			}
		})
	}

	o.vp.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	o.vp.SetWidth(o.width)
	o.vp.SetHeight(o.height)
	o.vp.SetContent(text)

	return tea.NewView(o.vp.View())
}

func (o *operationsView) Title() string {
	return "Operations"
}

func (o *operationsView) Icon() *iconVariant {
	return o.v.Icons().CogWheel
}

// HandleRouterMsg handles the routed message.
func (o *operationsView) HandleRouterMsg(m routerMsg) tea.Cmd {
	return handleRouterMsg(o, m)
}

func (o *operationsView) createNewOperation(msg opCreateMsg) (*opRunningInfo, opInvoker, error) {
	if msg.id == "" {
		return nil, nil, fmt.Errorf("%w: msg ID empty on create", errOpInternal)
	}

	_, ok := o.mapIDToIndex[msg.id]
	if ok {
		return nil, nil, errOpAlreadyInProgress
	}

	info := newOpRunningInfo(o.v, msg.opCreationInfo)

	o.orderedList = append(o.orderedList, info)
	o.mapIDToIndex[info.id] = len(o.orderedList) - 1

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
	delete(o.mapIDToIndex, msg.id)
}

func (o *operationsView) deleteOperationCmd(id string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return msgOpDelete(id)
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

func newOpCreationInfo(desc, msg string) opCreationInfo {
	return opCreationInfo{"", desc, msg}
}

func (o *opCreationInfo) updateID(id string) {
	o.id = id
}

type opRunningInfo struct {
	opCreationInfo

	v rootView
}

func newOpRunningInfo(v rootView, creationInfo opCreationInfo) *opRunningInfo {
	return &opRunningInfo{opCreationInfo: creationInfo, v: v}
}

func (o *opRunningInfo) info(msg string) {
	o.v.SendMsg(msgOpUpdate(o.opCreationInfo, "", msg))
}

func (o *opRunningInfo) updateDescription(desc string) {
	o.v.SendMsg(msgOpUpdate(o.opCreationInfo, desc, ""))
}

func (o *opRunningInfo) opSuccess(state actionStateSpec) routerMsg {
	return msgAdActionUpdate(o.id, state)
}

func (o *opRunningInfo) opError(err error) routerMsg {
	return msgOpError(err)
}

type opCreateMsg struct {
	opCreationInfo

	opAction opInvoker
}

func msgOpCreate(creationInfo opCreationInfo, action opInvoker) routerMsg {
	return viewIDOperations.routerMessage(opCreateMsg{opCreationInfo: creationInfo, opAction: action})
}

type opUpdateMsg struct {
	opCreationInfo
}

func msgOpUpdate(creationInfo opCreationInfo, desc, msg string) routerMsg {
	creationInfo.description = desc
	creationInfo.message = msg

	return viewIDOperations.routerMessage(opUpdateMsg{opCreationInfo: creationInfo})
}

type opDeleteMsg struct {
	id string
}

func msgOpDelete(id string) routerMsg {
	return viewIDOperations.routerMessage(opDeleteMsg{id})
}

type opErrorMsg struct {
	err error
}

func msgOpError(err error) routerMsg {
	return viewIDOperations.routerMessage(opErrorMsg{err})
}
