package adtree

import (
	"strconv"

	"github.com/ayn2op/tview"
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/keybindings"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

type subNode struct {
	*adNode
	noder adNoder
}

func newSubNode(adn *adNode, parentNoder adNoder) *subNode {
	return &subNode{adn, parentNoder}
}

// AdNode returns the base node data.
func (s *subNode) AdNode() *adNode {
	return s.adNode
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (s *subNode) HandleKeys(p tview.KeyMsg) (ui.Router, bool) {
	return s.noder.HandleKeys(p)
}

// Information returns a message to print adapter/device information.
func (s *subNode) Information() []text.Span {
	return s.noder.Information()
}

func (s *subNode) refreshNode(style tcell.Style) {
	c := s.Content()

	for i := range c {
		c[i].Style = style
	}
}

type alNode struct {
	*subNode
	actions []*acNode
}

func newActionsListNode(a *adNode, parentNoder adNoder) *alNode {
	al := &alNode{}

	alAdNode := newAdNode(
		a.id.appendSubNodeNib(nibActionsList),
		al, a.rv,
		false, 2,
	)

	al.subNode = newSubNode(alAdNode, parentNoder)
	al.AdNode().SetContent(buildLabel(al.Icon(), "Actions", al.NodeStyle()))

	alAdNode.SetChildren(make([]*ctree.Node, 0, 5))
	a.AddChild(alAdNode.Node)

	return al
}

// NodeStyle returns the style to be applied for this node's label.
func (a *alNode) NodeStyle() tcell.Style {
	return theme.Current().ADTree.ActionsList.Style
}

// Icon returns the Icon associated with this node.
func (a *alNode) Icon() string {
	return theme.Icons().Actions.String()
}

// AddAction adds an action item to the list.
func (a *alNode) AddAction(ac *acNode) {
	a.AdNode().AddChild(ac.AdNode().Node)
	a.actions = append(a.actions, ac)
}

// ClearActions clears all the actions from the list.
func (a *alNode) ClearActions() {
	a.AdNode().ClearChildren()
	a.actions = nil
}

// Actions returns the currently stored actions.
func (a *alNode) Actions() []*acNode {
	return a.actions
}

// Refresh refreshes the content of the node.
func (a *alNode) Refresh() {
	a.subNode.refreshNode(a.NodeStyle())
}

type dlNode struct {
	*subNode
	anode *adapterNode
}

func newDeviceListNode(a *adapterNode, devices []bluetooth.DeviceData) *dlNode {
	dl := &dlNode{}

	dlAdNode := newAdNode(
		a.id.appendSubNodeNib(nibDevicesList),
		dl, a.rv,
		true, 2,
	)

	dl.subNode = newSubNode(dlAdNode, a)
	dl.anode = a
	dl.AdNode().SetContent(buildLabel(dl.Icon(), "Devices", dl.NodeStyle()))

	for _, device := range devices {
		newDeviceNode(a.rv, dl, device)
	}

	a.AddChild(dlAdNode.Node)

	return dl
}

// NodeStyle returns the style to be applied for this node's label.
func (a *dlNode) NodeStyle() tcell.Style {
	return theme.Current().ADTree.DevicesList
}

// Icon returns the Icon associated with this node.
func (a *dlNode) Icon() string {
	return theme.Icons().Devices.String()
}

// AdapterNode returns the adapter node.
func (a *dlNode) AdapterNode() *adapterNode {
	return a.anode
}

// Refresh refreshes the content of the node.
func (a *dlNode) Refresh() {
	a.subNode.refreshNode(a.NodeStyle())
}

type acNode struct {
	*subNode

	actionState *adActionState
}

func newActionNode(a *alNode, key keybindings.Keybinding, state actionStateSpec, isToggleable bool, invoker actionInvoker) *acNode {
	ac := &acNode{}

	actionNode := newAdNode(
		a.AdNode().id.appendSubNodeTextNib(nibAction, strconv.Itoa(int(key.ID))),
		ac, a.AdNode().rv,
		false, 2,
	)

	ac.subNode = newSubNode(actionNode, a.noder)
	ac.actionState = newAdActionState(key, state, isToggleable, invoker)

	muted := theme.Current().App.PlaceHolder
	actionNode.SetMarkers(stext.NewLine(
		stext.NewSpan("(", muted),
		stext.NewSpan(key.KeyName, muted),
		stext.NewSpan(")", muted),
	))

	a.AddAction(ac)

	return ac
}

// Refresh refreshes the content of the node.
func (a *acNode) Refresh() {
	a.subNode.refreshNode(a.NodeStyle())
}

// NodeStyle returns the style to be applied for this node's label.
func (a *acNode) NodeStyle() tcell.Style {
	return theme.Current().ADTree.ActionsList.Nodes
}

// Icon returns the Icon associated with this node.
func (a *acNode) Icon() string {
	return ""
}
