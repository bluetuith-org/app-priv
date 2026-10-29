package adtree

import (
	"log/slog"
	"slices"
	"strconv"

	"github.com/ayn2op/tview/text"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/bluetuith-org/bluetuith/ui/widgets/stext"
)

// rootNode describes a root tree node.
type rootNode struct {
	*adNode
	emptyNoder

	title []text.Span

	adapterCount, devicesCount int
	acPos, dcPos               int
}

// newRootNode creates a new root node.
func newRootNode(rv ui.RootView) *rootNode {
	rn := &rootNode{}

	rootAdNode := newAdNode(
		newRootNodeID(),
		rn, rv,
		true, 2,
	)

	rootAdNode.SetSelectable(false)
	rn.adNode = rootAdNode

	rootAdNode.SetChildren(make([]*ctree.Node, 0, 10))
	rn.Refresh()

	return rn
}

// AdNode returns the base node data.
func (r *rootNode) AdNode() *adNode {
	return r.adNode
}

// RootID returns the root ID for the subnodes of an adapter or device node.
func (r *rootNode) RootID() string {
	return r.id.String()
}

// Refresh refreshes the content of the node.
func (r *rootNode) Refresh() {
	r.renderTitle()
}

func (r *rootNode) Title() []text.Span {
	return r.title
}

func (r *rootNode) renderTitle() {
	var line []text.Span

	tstyle := theme.Current().ADTree.TitleBarIndicator

	line = append(line, stext.NewSpan(" Adapters (", tstyle))
	line = append(line, stext.NewSpan(strconv.Itoa(r.adapterCount), tstyle))
	line = append(line, stext.NewSpan(") / Devices (", tstyle))
	line = append(line, stext.NewSpan(strconv.Itoa(r.devicesCount), tstyle))
	line = append(line, stext.NewSpan(") ", tstyle))

	r.title = line
	r.acPos, r.dcPos = 1, 3
}

func (r *rootNode) updateStats(fn func(rt *rootNode)) {
	fn(r)

	titleLength := len(r.title)
	if titleLength <= r.acPos || titleLength <= r.dcPos {
		return
	}

	r.title[r.acPos].Text = strconv.Itoa(r.adapterCount)
	r.title[r.dcPos].Text = strconv.Itoa(r.devicesCount)
}

func (r *rootNode) addAdapter(adapter bluetooth.AdapterData) {
	devices, err := r.rv.Session().Adapter(adapter.AdapterAddress).Devices()
	if err != nil {
		slog.Error("An error occurred while getting device information from the adapter",
			"internal error", err,
			"adapter", getAdapterDisplayName(r.rv, adapter, true),
		)

		return
	}

	dvCount := len(devices)
	r.updateStats(func(rt *rootNode) {
		rt.adapterCount++
		rt.devicesCount += dvCount
	})

	newAdapterNode(r.rv, r, adapter, devices)
}

func (r *rootNode) addDevice(device bluetooth.DeviceData) {
	dlNodeID := newAdapterNodeID(device.AdapterAddress()).appendSubNodeNib(nibDevicesList)

	deviceListNode, _, ok := findTreeNodeByID[*dlNode](r, dlNodeID.String())
	if !ok {
		return
	}

	r.updateStats(func(rt *rootNode) {
		rt.devicesCount++
	})

	newDeviceNode(r.rv, deviceListNode, device)
}

func (r *rootNode) updateAdapter(adapterEvent bluetooth.AdapterEventData, remove bool) {
	adNodeID := newAdapterNodeID(adapterEvent.AdapterAddress)

	adapterNode, pos, ok := findTreeNodeByID[*adapterNode](r, adNodeID.String())
	if !ok {
		return
	}

	if remove {
		parentNode := r.Node
		if parentNode != nil {
			parentNode.SetChildren(slices.Delete(parentNode.Children(), pos, pos+1))
		}

		dvCount := len(adapterNode.Children())
		r.updateStats(func(rt *rootNode) {
			rt.adapterCount = max(rt.adapterCount-1, 0)
			rt.devicesCount = max(rt.devicesCount-dvCount, 0)
		})

		return
	}

	adapterNode.SetAdapterEventData(adapterEvent)
}

func (r *rootNode) updateDevice(deviceEvent bluetooth.DeviceEventData, remove bool) {
	dvNodeID := newDeviceNodeID(deviceEvent.DeviceAddress)

	deviceNode, pos, ok := findTreeNodeByID[*deviceNode](r, dvNodeID.String())
	if !ok {
		return
	}

	if remove {
		parentNode, ok := deviceNode.ParentNode()
		if ok {
			node := parentNode.AdNode()
			node.SetChildren(slices.Delete(node.Children(), pos, pos+1))
		}

		r.updateStats(func(rt *rootNode) {
			rt.devicesCount = max(rt.devicesCount-1, 0)
		})

		return
	}

	deviceNode.SetDeviceEventData(deviceEvent)
}

// findTreeNodeByID walks the nodes starting from the root node, splitting "id" into individual
// chunks until the matching node is found by its ID.
func findTreeNodeByID[V any](rootNode *rootNode, id string) (V, int, bool) {
	var val V

	currNode := rootNode.adNode
	currNodePos := -1

	for nodeID := range iterNodeID(id) {
		found := false

		for pos, node := range currNode.Children() {
			data := nget(node).AdNode()
			if data.ID() == nodeID {
				currNode = data
				currNodePos = pos

				found = true
				break
			}
		}

		if !found {
			return val, -1, false
		}
	}

	if m, ok := currNode.Node.Reference().(V); ok {
		return m, currNodePos, currNode != nil
	}

	return val, -1, false
}

var _ adNoder = (*rootNode)(nil)
