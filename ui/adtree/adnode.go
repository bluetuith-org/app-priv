package adtree

import (
	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
)

// adNode describes a node in an adapter-device tree.
// noder must not be nil.
type adNode struct {
	*ctree.Node

	id nodeID

	markers []text.Span
	mwidth  int

	content []text.Span
	cwidth  int

	rv ui.RootView
}

func newAdNode(id nodeID, noder adNoder, rv ui.RootView, expanded bool, indent int) *adNode {
	data := &adNode{}

	node := ctree.NewNode(noder).SetIndent(indent).SetExpanded(expanded)
	data.Node = node

	data.id = id
	data.rv = rv

	return data
}

// ID returns the identifier for this node.
func (a *adNode) ID() nodeID {
	return a.id
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
