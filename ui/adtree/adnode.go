package adtree

import (
	"slices"
	"strconv"
	"strings"

	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
	"github.com/gdamore/tcell/v3"
	"github.com/google/uuid"
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

// adNode describes a node in an adapter-device tree.
// noder must not be nil.
type adNode struct {
	*ctree.Node

	id       nodeID
	nodeType adTreeNodeType

	noder adNoder

	markers []text.Span
	mwidth  int

	content []text.Span
	cwidth  int

	rv ui.RootView
}

func newAdNode(ntype adTreeNodeType, expanded bool, indent int, id nodeID, noder adNoder, rv ui.RootView) *adNode {
	data := &adNode{}

	node := ctree.NewNode(noder).SetIndent(indent).SetExpanded(expanded)
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

func (a *adNode) ParentNode() (adNoder, bool) {
	parent, ok := a.Node.Parent()
	if !ok {
		return nil, false
	}

	return nget(parent), true
}

func nget(n *ctree.Node) adNoder {
	return n.Data[adNoder]()
}

func nvget[v any](n *ctree.Node) v {
	return n.Data[v]()
}

func printPlaceHolder(v string) string {
	if v == "" {
		return "<None>"
	}

	return v
}

func printAddress(rv ui.RootView, address bluetooth.MacAddress, placeholder bool) string {
	if rv.Configuration().Values.RedactAddress {
		if placeholder {
			return "<Redacted>"
		}

		return ""
	}

	return address.String()
}

// getAdapterDisplayName returns the display name of the adapter.
func getAdapterDisplayName(rv ui.RootView, adapterData bluetooth.AdapterData, withAddress bool) string {
	var display string

	redactAddress := rv.Configuration().Values.RedactAddress

	if adapterData.UniqueName != "" {
		display = adapterData.UniqueName
	}

	if name, ok := adapterData.Name.Get(); ok {
		display = name
	}

	var sb strings.Builder

	if display == "" && !redactAddress {
		display = adapterData.Address.String()
		withAddress = false
	}

	sb.WriteString(printPlaceHolder(display))
	if withAddress && !redactAddress {
		sb.WriteString(" (")
		sb.WriteString(adapterData.Address.String())
		sb.WriteString(")")
	}

	return sb.String()
}

// getDeviceDisplayName returns the display name for the device.
func getDeviceDisplayName(rv ui.RootView, deviceData bluetooth.DeviceEventData, withAddress bool) string {
	var display string

	redactAddress := rv.Configuration().Values.RedactAddress

	if alias, ok := deviceData.Alias.Get(); ok {
		display = alias
	}

	if name, ok := deviceData.Name.Get(); ok {
		display = name
	}

	var sb strings.Builder

	if display == "" && !redactAddress {
		display = deviceData.Address.String()
		withAddress = false
	}

	sb.WriteString(printPlaceHolder(display))
	if withAddress && !redactAddress {
		sb.WriteString(" (")
		sb.WriteString(deviceData.Address.String())
		sb.WriteString(")")
	}

	return sb.String()
}

func printServiceUUIDs(b *stext.Builder, uuids uuid.UUIDs) {
	if len(uuids) == 0 {
		b.AddKVFunc("Services:", func(bb *stext.Builder) {
			bb.Newline()
			bb.Append("<None>", theme.Current().App.PlaceHolder)
		})
		return
	}

	b.AddKV("Services ("+strconv.Itoa(len(uuids))+"):", "")

	nameStyle := theme.Current().Info.ServiceName
	uuidStyle := theme.Current().Info.ServiceIdentifier

	slices.SortFunc(uuids, func(a, b uuid.UUID) int {
		aid, bid := a.ID(), b.ID()
		switch {
		case aid < bid:
			return -1

		case aid > bid:
			return 1
		}

		return 0
	})

	for _, u := range uuids {
		name := "- " + bluetooth.ServiceType(u)

		rep := ""
		switch {
		case bluetooth.IsVendorSpecificUUID(u):
			rep = u.String()

		default:
			rep = "0x" + strconv.FormatUint(uint64(u.ID()), 16)
		}

		b.AddKVStyledFunc(name, nameStyle, func(bb *stext.Builder) {
			bb.Newline()
			bb.Append("(", uuidStyle)
			bb.Append(rep, uuidStyle)
			bb.Append(")", uuidStyle)
			bb.Newline()
		})
	}
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
