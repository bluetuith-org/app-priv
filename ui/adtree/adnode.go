package adtree

import (
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
)

// adTreeNodeType specifies the type of the tree node.
type adTreeNodeType uint8

// The different types of tree nodes.
const (
	nodeTypeRoot adTreeNodeType = iota
	nodeTypeAdapter
	nodeTypeDevice
	nodeTypeDevicesList
	nodeTypeAction
	nodeTypeActionsList
)

// adSubNodePos specifies the position of sub-nodes under
// adapter, device or root nodes.
type adSubNodePos uint8

// The different types of sub nodes.
const (
	relPosActionsListNode adSubNodePos = iota
	relPosDevicesListNode
)

// adNode describes a node in an adapter-device tree.
// noder must not be nil.
type adNode struct {
	*ctree.Node

	id       nodeID
	nodeType adTreeNodeType

	noder       adNoder
	actionState *adActionState

	markers []text.Span
	mwidth  int

	content []text.Span
	cwidth  int

	rv ui.RootView
}

func newAdNode(ntype adTreeNodeType, expanded bool, indent int, id nodeID, noder adNoder, rv ui.RootView) *adNode {
	data := &adNode{}

	node := ctree.NewNode(data).SetIndent(indent).SetExpanded(expanded)
	data.Node = node

	data.id = id
	data.nodeType = ntype
	data.noder = noder
	data.rv = rv

	return data
}

// ID returns the identifier for this node.
func (a *adNode) ID() string {
	return a.id.String()
}

// Noder returns the interface that implements [adNoder].
func (a *adNode) Noder() adNoder {
	return a.noder
}

// Content returns the name of the node.
func (a *adNode) Content() []text.Span {
	return a.content
}

// SetContent sets the name for the node.
func (a *adNode) SetContent(content []text.Span) {
	a.content = content
	a.cwidth = stext.LineWidth(a.content)
}

func (a *adNode) labelWidth() int {
	return a.mwidth + a.cwidth
}

func (a *adNode) Markers() []text.Span {
	return a.markers
}

func (a *adNode) SetMarkers(markers []text.Span) {
	a.markers = append(markers, stext.NewSpan(" ", theme.Current().Global))
	a.mwidth = stext.LineWidth(a.markers)
}

func (a *adNode) ParentNode() (*adNode, bool) {
	parent, ok := a.Node.Parent()
	if !ok {
		return nil, false
	}

	return parent.Data[*adNode](), true
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

func buildLabel(icon, name string, style tcell.Style) []text.Span {
	line := make([]text.Span, 0, 3)

	if icon != "" {
		line = append(
			line,
			stext.NewSpan(icon, style),
			stext.NewSpan(" ", style),
		)
	}

	if name != "" {
		line = append(line, stext.NewSpan(name, style))
	}

	return line
}
