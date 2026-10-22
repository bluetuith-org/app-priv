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
	sub   *adNode
	noder adNoder
}

func newSubNode(adn *adNode, parentNoder adNoder) *subNode {
	return &subNode{adn, parentNoder}
}

// AdNode returns the base node data.
func (s *subNode) AdNode() *adNode {
	return s.sub
}

// HandleKeys handles a [tview.KeyMsg] and returns a message.
func (s *subNode) HandleKeys(p tview.KeyMsg) (ui.Router, bool) {
	return s.noder.HandleKeys(p)
}

// RootID returns the root ID for the subnodes of an adapter or device node.
func (s *subNode) RootID() string {
	return s.noder.RootID()
}

// Information returns a message to print adapter/device information.
func (s *subNode) Information() []text.Span {
	return s.noder.Information()
}

// PopulateActions populates all actions within a node of type [nodeTypeAction].
func (s *subNode) PopulateActions() {
	s.noder.PopulateActions()
}

// Refresh refreshes the content of the node.
func (s *subNode) Refresh() {
	style := s.sub.noder.NodeStyle()
	c := s.sub.Content()

	for i := range c {
		c[i].Style = style
	}
}

// UpdateActionNode updates the node of type [nodeTypeAction] with the message.
func (s *subNode) UpdateActionNode(actionNode *acNode, updateMsg actionUpdateMsg) {
	s.noder.UpdateActionNode(actionNode, updateMsg)
}

// SetAdapterEventData sets the adapter event data for the node.
func (s *subNode) SetAdapterEventData(ev bluetooth.AdapterEventData) {
	s.noder.SetAdapterEventData(ev)
}

// SetDeviceEventData sets the device event data for the node.
func (s *subNode) SetDeviceEventData(ev bluetooth.DeviceEventData) {
	s.noder.SetDeviceEventData(ev)
}

type alNode struct {
	*subNode
}

func newActionsListNode(a *adNode) *alNode {
	al := &alNode{}

	alAdNode := newAdNode(
		nodeTypeActionsList,
		false, 2,
		a.id.appendSubNodeNib(nibActionsList),
		al, a.rv,
	)

	al.subNode = newSubNode(alAdNode, a.noder)
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

type dlNode struct {
	*subNode
	anode *adapterNode
}

func newDeviceListNode(a *adapterNode, devices []bluetooth.DeviceData) *dlNode {
	dl := &dlNode{}

	dlAdNode := newAdNode(
		nodeTypeDevicesList,
		true, 2,
		a.id.appendSubNodeNib(nibDevicesList),
		dl, a.rv,
	)

	dl.subNode = newSubNode(dlAdNode, a.noder)
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

func (a *dlNode) AdapterNode() *adapterNode {
	return a.anode
}

type acNode struct {
	*subNode

	actionState *adActionState
}

func newActionNode(a *alNode, key keybindings.Keybinding, state actionStateSpec, isToggleable bool, invoker actionInvoker) *acNode {
	ac := &acNode{}

	actionNode := newAdNode(
		nodeTypeAction,
		false, 2,
		a.AdNode().id.appendSubNodeTextNib(nibAction, strconv.Itoa(int(key.ID))),
		ac, a.AdNode().rv,
	)

	ac.subNode = newSubNode(actionNode, a.noder)
	ac.actionState = newAdActionState(key, state, isToggleable, invoker)

	muted := theme.Current().App.PlaceHolder
	actionNode.SetMarkers(stext.NewLine(
		stext.NewSpan("(", muted),
		stext.NewSpan(key.KeyName, muted),
		stext.NewSpan(")", muted),
	))

	a.AdNode().AddChild(actionNode.Node)
	a.UpdateActionNode(ac, emptyActionUpdateMsg())

	return ac
}

// NodeStyle returns the style to be applied for this node's label.
func (a *acNode) NodeStyle() tcell.Style {
	return theme.Current().ADTree.ActionsList.Nodes
}

// Icon returns the Icon associated with this node.
func (a *acNode) Icon() string {
	return ""
}
