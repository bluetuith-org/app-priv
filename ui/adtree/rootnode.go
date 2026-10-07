package adtree

import (
	"slices"
	"strconv"

	"github.com/ayn2op/tview/richtext"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/theme"
	"github.com/bluetuith-org/bluetuith/ui"
	"github.com/bluetuith-org/bluetuith/ui/widgets/ctree"
	"github.com/gdamore/tcell/v3"
)

// rootNode describes a root tree node.
type rootNode struct {
	*adNode
	emptyNoder

	adapterCount, devicesCount int
}

// newRootNode creates a new root node.
func newRootNode(rv ui.RootView) *rootNode {
	rn := &rootNode{}

	rootAdNode := newAdNode(
		nodeTypeRoot,
		true, newRootNodeID(),
		rn, rv,
	)

	rn.adNode = rootAdNode
	rootAdNode.SetChildren(make([]*ctree.Node, 0, 10))

	rn.Refresh()

	return rn
}

// NodeStyle returns the style to be applied for this node's label.
func (r *rootNode) NodeStyle() tcell.Style {
	return theme.Current().ADTree.Title
}

// Icon returns the Icon associated with this node.
func (r *rootNode) Icon() string {
	return theme.Icons().Adapters.String()
}

// RootID returns the root ID for the subnodes of an adapter or device node.
func (r *rootNode) RootID() string {
	return r.id.String()
}

// Refresh refreshes the content of the node.
func (r *rootNode) Refresh() {
	r.SetContent(buildLabel(r.Icon(), "Adapters", r.NodeStyle()))
}

// Information returns a message to print adapter/device information.
func (r *rootNode) Information() (string, richtext.Text) {
	builder := ui.NewBuilder(theme.Current().Global)

	style := theme.Current().Info.PropertyName
	vstyle := theme.Current().Info.PropertyValue
	tstyle := theme.Current().Info.Title

	builder.Appendln(" Overview ", tstyle)
	builder.Newline()

	builder.AddKV("Adapters:", style, strconv.Itoa(r.adapterCount), vstyle)
	builder.AddKV("Devices:", style, strconv.Itoa(r.devicesCount), vstyle)

	return r.id.String(), builder.Text()
}

func (r *rootNode) addAdapter(adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) {
	r.adapterCount++
	r.devicesCount += len(devices)

	newAdapterNode(r.rv, r, adapter, devices)
}

func (r *rootNode) addDevice(device bluetooth.DeviceData) {
	dlNodeID := newAdapterNodeID(device.AdapterAddress()).appendSubNodeNib(nibDevicesList)

	deviceListNode, _, ok := findTreeNodeByID(r, dlNodeID.String(), nodeTypeDevicesList)
	if !ok {
		return
	}

	r.devicesCount++

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

		r.adapterCount = max(r.adapterCount-1, 0)
		r.devicesCount = max(r.devicesCount-len(adapterNode.Children()), 0)

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

		r.devicesCount = max(r.devicesCount-1, 0)

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
