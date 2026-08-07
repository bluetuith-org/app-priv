package views

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Digital-Shane/treeview/v2"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/appfeatures"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
)

type (
	treeUpdate viewUpdate
	nodeID     = string
)

type adTree struct {
	tree    *treeview.Tree[*adTreeNode]
	tuiTree *treeview.TuiTreeModel[*adTreeNode]

	provider *adTreeProvider

	v   rootView
	ctx context.Context

	width, height int

	focused   bool
	focusedID string

	mu sync.Mutex
}

// InitializeView initializes the view.
func (a *adTree) InitializeView(_ *appfeatures.FeatureSet) (inited bool, err error) {
	a.focused = true
	a.width, a.height = 120, 30

	a.ctx = context.Background()
	a.provider = newAdTreeProvider()

	rootNode := newAdTreeNode().withRootNode()
	rootNode.node.SetChildren(make([]*treeview.Node[*adTreeNode], 0, 10))

	kmap := treeview.DefaultKeyMap()
	kmap.Quit = []string{}

	a.tree = treeview.NewTree(
		[]*treeview.Node[*adTreeNode]{rootNode.node},
		treeview.WithExpandFunc(a.expandFunc),
		treeview.WithProvider(a.provider),
	)

	a.tuiTree = treeview.NewTuiTreeModel(
		a.tree,
		treeview.WithTuiDisableNavBar[*adTreeNode](true),
		treeview.WithTuiKeyMap[*adTreeNode](kmap),
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
	a.updateTreeFn(true, func(rootNode *treeview.Node[*adTreeNode]) {
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
	return a.populate
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (a *adTree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		msg = a.Resize(m.Width, m.Height)

	case tea.KeyPressMsg:
		if !a.focused {
			return a, nil
		}
	}

	if a.tuiTree != nil {
		_, cmd = a.tuiTree.Update(msg)
	}

	return a, cmd
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

func (a *adTree) populate() tea.Msg {
	adapters, err := a.v.Session().Adapters()
	if err != nil {
		// TODO: Log error messages
		return nil
	}

	for _, adapter := range adapters {
		a.addAdapter(adapter, true)
	}

	return treeUpdate{}
}

func (a *adTree) addAdapter(adapter bluetooth.AdapterData, lock bool) {
	devices, err := a.v.Session().Adapter(adapter.AdapterAddress).Devices()
	if err != nil {
		// TODO:Log error messages
		_ = err
	}

	ln := 1
	if len(devices) > 0 {
		ln = len(devices)
	}

	a.updateTreeFn(lock, func(rootNode *treeview.Node[*adTreeNode]) {
		adNode := newAdTreeNode().withAdapter(adapter)

		actionsNode := newAdTreeNode().withActionsList(adNode.getID())
		actionsNode.node.SetChildren(make([]*treeview.Node[*adTreeNode], 0, 5))

		devicesListNode := newAdTreeNode().withDevicesList(adNode.getID())
		devicesListNode.node.SetChildren(make([]*treeview.Node[*adTreeNode], 0, ln))

		for _, device := range devices {
			a.addDevice(device, devicesListNode, false)
		}

		adNode.node.SetChildren([]*treeview.Node[*adTreeNode]{actionsNode.node, devicesListNode.node})

		rootCh := rootNode.Children()
		rootCh = append(rootCh, adNode.node)
		rootNode.SetChildren(rootCh)
	})
}

func (a *adTree) addDevice(device bluetooth.DeviceData, parentNode *adTreeNode, lock bool) {
	if parentNode == nil {
		n, _ := a.tuiTree.FindByID(context.Background(), filepath.Join(devicesListIDNib, adapterIDNib, device.AssociatedAdapter.String()))
		if n == nil {
			return
		}

		parentNode = *n.Data()
	}

	a.updateTreeFn(lock, func(*treeview.Node[*adTreeNode]) {
		dvNode := newAdTreeNode().withDevice(device)

		actionsNode := newAdTreeNode().withActionsList(dvNode.getID())
		actionsNode.node.SetChildren(make([]*treeview.Node[*adTreeNode], 0, 5))

		dvNode.node.SetChildren([]*treeview.Node[*adTreeNode]{actionsNode.node})

		adCh := parentNode.node.Children()
		adCh = append(adCh, dvNode.node)

		parentNode.node.SetChildren(adCh)
	})
}

func (a *adTree) expandFunc(node *treeview.Node[*adTreeNode]) bool {
	data := *node.Data()

	return data.nodeType == nodeTypeRoot ||
		data.nodeType == nodeTypeAdapter ||
		data.nodeType == nodeTypeDevice ||
		data.nodeType == nodeTypeDevicesList
}

func (a *adTree) updateTreeFn(lock bool, fn func(rootNode *treeview.Node[*adTreeNode])) {
	if lock {
		a.mu.Lock()
		defer a.mu.Unlock()
	}

	fn(a.tuiTree.Nodes()[0])
}

type adTreeNodeType uint8

const (
	nodeTypeRoot adTreeNodeType = iota
	nodeTypeAdapter
	nodeTypeDevice
	nodeTypeDevicesList
	nodeTypeAction
	nodeTypeActionsList
)

type idNib = string

const (
	devicesListIDNib idNib = "dl"
	actionsListIDNib idNib = "al"
	actionIDNib      idNib = "ac"
	adapterIDNib     idNib = "ad"
	deviceIDNib      idNib = "dv"
)

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
func (a *adTreeProvider) Icon(node *treeview.Node[*adTreeNode]) string {
	// TODO: Ascii
	const collapseIndicator = string('\u25b6')
	const expandIndicator = string('\u25bc')

	return useBuffer(func(b *strings.Builder) {
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
func (a *adTreeProvider) Format(node *treeview.Node[*adTreeNode]) string {
	return node.Name()
}

// Style supplies the lipgloss style for the node based on its focus state.
func (a *adTreeProvider) Style(_ *treeview.Node[*adTreeNode], isFocused bool) lipgloss.Style {
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
