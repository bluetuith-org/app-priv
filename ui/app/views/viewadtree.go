package views

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/Digital-Shane/treeview/v2"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
)

type adTree struct {
	tree    *treeview.Tree[adTreeNode]
	tuiTree *treeview.TuiTreeModel[adTreeNode]

	provider *adTreeProvider

	v   rootView
	ctx context.Context

	width, height int

	focused   bool
	focusedID string

	search     bool
	searchTerm string

	mu sync.Mutex
}

// ViewID returns the view's ID.
func (a *adTree) ViewID() viewID {
	return viewIDAdTree
}

// InitializeView initializes the view.
func (a *adTree) InitializeView(_ *appfeatures.FeatureSet) (inited bool, err error) {
	a.focused = true
	a.width, a.height = 120, 30

	a.ctx = context.Background()
	a.provider = newAdTreeProvider()

	rootNode := newRootAdNode(a)

	kmap := treeview.KeyMap{}

	a.tree = treeview.NewTree(
		[]*treeview.Node[adTreeNode]{rootNode.node},
		treeview.WithExpandFunc(a.expandFunc),
		treeview.WithProvider(a.provider),
	)

	a.tuiTree = treeview.NewTuiTreeModel(
		a.tree,
		treeview.WithTuiDisableNavBar[adTreeNode](true),
		treeview.WithTuiKeyMap[adTreeNode](kmap),
	)

	return true, nil
}

// SetRootView sets the root view.
func (a *adTree) SetRootView(v rootView) {
	a.v = v
}

// AttachToTabView attaches this view to the tabbed view.
func (a *adTree) AttachToTabView() (tabSection, bool) {
	return nil, false
}

// Resize resizes the view.
func (a *adTree) Resize(width, height int) tea.WindowSizeMsg {
	a.width = (width / 2) - 1
	a.height = height - 2

	return tea.WindowSizeMsg{
		Width:  a.width,
		Height: a.height,
	}
}

// SetFocus sets whether the view is currently focused.
func (a *adTree) SetFocus(focused bool) {
	a.updateTreeFn(true, func(rootNode *treeview.Node[adTreeNode]) {
		a.focused = focused

		if !focused {
			a.focusedID = a.tuiTree.GetFocusedID()
			a.tuiTree.ClearAllFocus()

			return
		}

		_, err := a.tuiTree.SetFocusedID(a.ctx, a.focusedID)
		if errors.Is(err, treeview.ErrNodeNotFound) {
			a.tuiTree.SetFocusedID(a.ctx, rootNode.ID())
		}
	})
}

// GetFocus gets whether the view is currently focused.
func (a *adTree) GetFocus() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.focused
}

// UpdateStyles updates the styles for the view.
func (a *adTree) UpdateStyles() {
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (a *adTree) Init() tea.Cmd {
	return a.populate()
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (a *adTree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		msg = a.Resize(m.Width, m.Height)
		a.tuiTree.Update(msg)

	case tea.KeyPressMsg:
		if !a.focused {
			return a, nil
		}

		focusedNode := a.tree.GetFocusedNode()
		if focusedNode == nil {
			return a, nil
		}

		adNode := focusedNode.Data()

		switch {
		case keybindings.MatchesKey(keybindings.KeyNavigateUp, m):
			a.tuiTree.NavigateUp()
			return a, nil

		case keybindings.MatchesKey(keybindings.KeyNavigateDown, m):
			a.tuiTree.NavigateDown()
			return a, nil

		case keybindings.MatchesKey(keybindings.KeySwitch, m):
			a.v.FocusTabView()
			return a, nil

		case keybindings.MatchesKey(keybindings.KeyClose, m):
			a.endSearch()
			return a, nil

		case keybindings.MatchesKey(keybindings.KeySelect, m):
			if a.endSearch() {
				return a, nil
			}

			if adNode.nodeType != nodeTypeAction {
				focusedNode.Toggle()
				return a, nil
			}

			return a, func() tea.Msg {
				return actionStateToOperation(focusedNode.ID(), adNode.actionState)
			}

		case keybindings.MatchesKey(keybindings.KeyFilter, m):
			if a.beginSearch() {
				return a, nil
			}

		default:
		}

		if a.handleSearch(m) {
			return a, nil
		}

		if nmsg, ok := adNode.noder.handleKeys(m); ok {
			return a, func() tea.Msg { return nmsg }
		}

	case actionUpdateMsg:
		return a, a.updateAction(m)
	}

	return a, nil
}

// View renders the program's UI, which can be a string or a [Layer]. The
// view is rendered after every Update.
func (a *adTree) View() tea.View {
	view := tea.View{
		Content:   "Loading adapters...",
		AltScreen: true,
	}

	if a.tuiTree == nil {
		return view
	}

	view.SetContent(lipgloss.JoinVertical(lipgloss.Left, " ", a.tuiTree.View().Content))

	return view
}

// handleRouterMsg handles the routed message.
func (a *adTree) handleRouterMsg(m routerMsg) tea.Cmd {
	return handleRouterMsg(a, m)
}

func (a *adTree) features() *appfeatures.FeatureSet {
	return a.v.Features()
}

func (a *adTree) session() bluetooth.Session {
	return a.v.Session()
}

func (a *adTree) populate() tea.Cmd {
	return func() tea.Msg {
		adapters, err := a.v.Session().Adapters()
		if err != nil {
			// TODO: Log error messages
			return nil
		}

		for _, adapter := range adapters {
			a.addAdapter(adapter, true)
		}

		return msgAdTreeView(treeUpdateMsg{})
	}
}

func (a *adTree) addAdapter(adapter bluetooth.AdapterData, lock bool) {
	devices, err := a.v.Session().Adapter(adapter.AdapterAddress).Devices()
	if err != nil {
		// TODO:Log error messages
		_ = err
		return
	}

	a.updateTreeFn(lock, func(rootNode *treeview.Node[adTreeNode]) {
		newAdapterAdNode(a, rootNode, adapter, devices)
	})
}

func (a *adTree) addDevice(device bluetooth.DeviceData, parentNode *adTreeNode, lock bool) {
	if parentNode == nil {
		n, _ := a.tuiTree.FindByID(
			context.Background(),
			newAdapterNodeID(bluetooth.NewAdapterAddress(device.AssociatedAdapter)).appendSubNodeNib(nibDevicesList).String(),
		)
		if n == nil {
			return
		}

		parentNode = n.Data()
	}

	a.updateTreeFn(lock, func(*treeview.Node[adTreeNode]) {
		newDeviceAdNode(a, parentNode.node, device)
	})
}

func (a *adTree) expandFunc(node *treeview.Node[adTreeNode]) bool {
	data := *node.Data()

	return data.nodeType == nodeTypeRoot ||
		data.nodeType == nodeTypeAdapter ||
		data.nodeType == nodeTypeDevice ||
		data.nodeType == nodeTypeDevicesList
}

func (a *adTree) updateAction(msg actionUpdateMsg) tea.Cmd {
	return func() tea.Msg {
		a.updateTreeFn(true, func(_ *treeview.Node[adTreeNode]) {
			node, err := a.tree.FindByID(context.Background(), msg.id)
			if err != nil {
				return
			}

			if node != nil {
				node.Data().noder.updateAction(node, msg)
			}
		})

		return nil
	}
}

func (a *adTree) updateTreeFn(lock bool, fn func(rootNode *treeview.Node[adTreeNode])) {
	if lock {
		a.mu.Lock()
		defer a.mu.Unlock()
	}

	fn(a.tuiTree.Nodes()[0])
}

func (a *adTree) beginSearch() bool {
	if a.search {
		return false
	}

	a.search = true
	a.tuiTree.BeginSearch()
	a.searchTerm = ""

	return true
}

func (a *adTree) endSearch() bool {
	if !a.search {
		return false
	}

	a.tuiTree.EndSearch()
	a.search = false

	focusedNodes := a.tuiTree.GetAllFocusedIDs()
	if len(focusedNodes) > 0 {
		a.focusedID = focusedNodes[0]
	}

	a.SetFocus(true)

	return true
}

func (a *adTree) handleSearch(m tea.KeyPressMsg) bool {
	if !a.search {
		return false
	}
	text := m.Key().Text
	key := m.String()

	if key == "backspace" && len(a.searchTerm) > 0 {
		a.searchTerm = a.searchTerm[:len(a.searchTerm)-1]
		a.tuiTree.Search(a.searchTerm)

		return true
	}

	if len(text) == 1 && text >= " " && text <= "~" {
		a.searchTerm += text
		a.tuiTree.Search(a.searchTerm)

		return true
	}

	if len(key) == 1 && key >= " " && key <= "~" {
		a.searchTerm += key
		a.tuiTree.Search(a.searchTerm)
	}

	return true
}

type adTreeProvider struct {
	defaultStyle, focusedStyle lipgloss.Style
}

func newAdTreeProvider() *adTreeProvider {
	p := &adTreeProvider{}

	p.defaultStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	p.focusedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("39")).
		Bold(true)

	return p
}

// Icon returns the leading glyph (e.g. folder / file symbol) for the node.
func (a *adTreeProvider) Icon(node *treeview.Node[adTreeNode]) string {
	// TODO: Ascii
	const collapseIndicator = string('\u25b6')
	const expandIndicator = string('\u25bc')

	return useBuffer(len(expandIndicator)+2, func(b *strings.Builder) {
		if node.HasChildren() {
			indicator := expandIndicator
			if !node.IsExpanded() {
				indicator = collapseIndicator
			}

			b.WriteString(indicator)
			b.WriteString(" ")
		}
	})
}

// Format converts the node's data into a human-readable label that follows the icon.
func (a *adTreeProvider) Format(node *treeview.Node[adTreeNode]) string {
	return node.Name()
}

// Style supplies the lipgloss style for the node based on its focus state.
func (a *adTreeProvider) Style(_ *treeview.Node[adTreeNode], isFocused bool) lipgloss.Style {
	if isFocused {
		return a.focusedStyle
	}

	return a.defaultStyle
}

// getAdapterDisplayName returns the display name of the adapter.
func getAdapterDisplayName(adapterData bluetooth.AdapterData) string {
	if name, ok := adapterData.Name.Get(); ok {
		return name
	}

	if adapterData.UniqueName != "" {
		return adapterData.UniqueName
	}

	return adapterData.Address.String()
}

// getDeviceDisplayName returns the display name for the device.
func getDeviceDisplayName(deviceData bluetooth.DeviceEventData) string {
	if name, ok := deviceData.Name.Get(); ok {
		return name
	}

	if alias, ok := deviceData.Alias.Get(); ok {
		return alias
	}

	return deviceData.Address.String()
}

type treeUpdateMsg struct{}

type adTreeMsgC interface {
	treeUpdateMsg | actionUpdateMsg
}

func msgAdTreeView[M adTreeMsgC](msg M) routerMsg {
	return viewIDAdTree.routerMessage(msg)
}
