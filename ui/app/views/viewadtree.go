package views

import (
	"context"
	"errors"
	"strings"

	"github.com/Digital-Shane/treeview/v2"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/ui/keybindings"
	"github.com/bluetuith-org/bluetuith/ui/theme"
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

	rootNode *rootAdNode
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
	a.provider = newAdTreeProvider(a.v)

	a.rootNode = newRootAdNode(a)

	kmap := treeview.KeyMap{}

	a.tree = treeview.NewTree(
		[]*treeview.Node[adTreeNode]{a.rootNode.Node},
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
	a.height = height

	return tea.WindowSizeMsg{
		Width:  a.width,
		Height: a.height,
	}
}

// SetFocus sets whether the view is currently focused.
func (a *adTree) SetFocus(focused bool) {
	a.focused = focused

	if !focused {
		a.focusedID = a.tuiTree.GetFocusedID()
		a.tuiTree.ClearAllFocus()

		return
	}

	_, err := a.tuiTree.SetFocusedID(a.ctx, a.focusedID)
	if errors.Is(err, treeview.ErrNodeNotFound) {
		a.tuiTree.SetFocusedID(a.ctx, a.rootNode.ID())
	}
}

// GetFocus gets whether the view is currently focused.
func (a *adTree) GetFocus() bool {
	return a.focused
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

		a.rootNode.RLock()
		defer a.rootNode.RUnlock()

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

			return a, acStateToOpMsg(focusedNode.ID(), adNode.actionState).sendRoutedMsg(a.v)

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
			return a, nmsg.sendRoutedMsg(a.v)
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

// HandleRouterMsg handles the routed message.
func (a *adTree) HandleRouterMsg(m routerMsg) tea.Cmd {
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
			if err := a.addAdapter(adapter); err != nil {
				_ = err
				return nil
			}
		}

		return msgAdTreeUpdate()
	}
}

func (a *adTree) addAdapter(adapter bluetooth.AdapterData) error {
	devices, err := a.v.Session().Adapter(adapter.AdapterAddress).Devices()
	if err != nil {
		// TODO:Log error messages
		_ = err
		return err
	}

	a.rootNode.addAdapter(adapter, devices)

	return nil
}

func (a *adTree) addDevice(device bluetooth.DeviceData) {
	a.rootNode.addDevice(device)
}

func (a *adTree) expandFunc(node *treeview.Node[adTreeNode]) bool {
	data := *node.Data()

	return data.nodeType == nodeTypeRoot ||
		data.nodeType == nodeTypeAdapter ||
		data.nodeType == nodeTypeDevice ||
		data.nodeType == nodeTypeDevicesList
}

func (a *adTree) updateAction(msg actionUpdateMsg) tea.Cmd {
	a.rootNode.updateAction(msg)
	return nil
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

	switch {
	case len(text) == 1 && text >= " " && text <= "~":
		a.searchTerm += text

	case len(key) == 1 && key >= " " && key <= "~":
		a.searchTerm += key
	}

	a.tuiTree.Search(a.searchTerm)

	return true
}

type adTreeProvider struct {
	defaultStyle, focusedStyle lipgloss.Style

	v rootView
}

func newAdTreeProvider(v rootView) *adTreeProvider {
	p := &adTreeProvider{}

	p.v = v

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
	collapseIndicator := theme.Icons().TriangleRight.GetIcon()
	expandIndicator := theme.Icons().TriangleDown.GetIcon()

	return useStringBuffer(len(expandIndicator)+2, func(b *strings.Builder) {
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

func msgAdTreeUpdate() routerMsg {
	return viewIDAdTree.routerMessage(treeUpdateMsg{})
}

func msgAdActionUpdate(id string, stateSpec actionStateSpec) routerMsg {
	return viewIDAdTree.routerMessage(actionUpdateMsg{id, stateSpec})
}
