package adtree

import (
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
		nodeTypeRoot,
		true, 2,
		newRootNodeID(),
		rn, rv,
	)

	rootAdNode.SetSelectable(false)
	rn.adNode = rootAdNode

	rootAdNode.SetChildren(make([]*ctree.Node, 0, 10))
	rn.Refresh()

	return rn
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

func (r *rootNode) addAdapter(adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) {
	dvCount := len(devices)
	r.updateStats(func(rt *rootNode) {
		rt.adapterCount++
		rt.devicesCount += dvCount
	})

	newAdapterNode(r.rv, r, adapter, devices)
}

func (r *rootNode) addDevice(device bluetooth.DeviceData) {
	dlNodeID := newAdapterNodeID(device.AdapterAddress()).appendSubNodeNib(nibDevicesList)

	deviceListNode, _, ok := findTreeNodeByID(r, dlNodeID.String(), nodeTypeDevicesList)
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

	adapterNode, pos, ok := findTreeNodeByID(r, adNodeID.String(), nodeTypeAdapter)
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

	adapterNode.noder.SetAdapterEventData(adapterEvent)
}

func (r *rootNode) updateDevice(deviceEvent bluetooth.DeviceEventData, remove bool) {
	dvNodeID := newDeviceNodeID(deviceEvent.DeviceAddress)

	deviceNode, pos, ok := findTreeNodeByID(r, dvNodeID.String(), nodeTypeDevice)
	if !ok {
		return
	}

	if remove {
		parentNode, ok := deviceNode.ParentNode()
		if ok {
			parentNode.SetChildren(slices.Delete(parentNode.Children(), pos, pos+1))
		}

		r.updateStats(func(rt *rootNode) {
			rt.devicesCount = max(rt.devicesCount-1, 0)
		})

		return
	}

	deviceNode.noder.SetDeviceEventData(deviceEvent)
}

func (r *rootNode) updateAction(updateMsg actionUpdateMsg) {
	node, _, ok := findTreeNodeByID(r, updateMsg.id, nodeTypeAction)
	if !ok {
		return
	}

	node.noder.UpdateActionNode(node, updateMsg)
}

// findTreeNodeByID walks the nodes starting from the root node, splitting "id" into individual
// chunks until the matching node is found by its ID.
func findTreeNodeByID(rootNode *rootNode, id string, nodeType adTreeNodeType) (*adNode, int, bool) {
	currNode := rootNode.adNode
	currNodePos := -1

	for nodeID := range iterNodeID(id) {
		found := false

		for pos, node := range currNode.Children() {
			data := node.Data[*adNode]()
			if data.ID() == nodeID {
				currNode = data
				currNodePos = pos

				found = true
				break
			}
		}

		if !found {
			return nil, -1, false
		}
	}

	return currNode, currNodePos, currNode != nil && currNode.nodeType == nodeType
}

var _ adNoder = (*rootNode)(nil)
