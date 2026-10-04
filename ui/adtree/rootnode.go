package adtree

import (
	"slices"
	"sync"

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

	sync.RWMutex
}

// newRootNode creates a new root node.
func newRootNode(rv ui.RootView) *rootNode {
	rn := &rootNode{}

	rootAdNode := newAdNode(
		nodeTypeRoot, "Adapters",
		true, newRootNodeID(),
		rn, rv,
	)

	rootAdNode.Node.SetChildren(make([]*ctree.Node, 0, 10))
	rn.adNode = rootAdNode

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

func (r *rootNode) addAdapter(adapter bluetooth.AdapterData, devices []bluetooth.DeviceData) {
	r.Lock()
	defer r.Unlock()

	newAdapterNode(r.rv, r, adapter, devices)
}

func (r *rootNode) addDevice(device bluetooth.DeviceData) {
	r.Lock()
	defer r.Unlock()

	dlNodeID := newAdapterNodeID(device.AdapterAddress()).appendSubNodeNib(nibDevicesList)

	deviceListNode, _, ok := findTreeNodeByID(r, dlNodeID.String(), nodeTypeDevicesList)
	if !ok {
		return
	}

	newDeviceNode(r.rv, deviceListNode, device)
}

func (r *rootNode) updateAdapter(adapterEvent bluetooth.AdapterEventData, remove bool) {
	r.Lock()
	defer r.Unlock()

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

		return
	}

	adapterNode.noder.SetAdapterEventData(adapterEvent)
}

func (r *rootNode) updateDevice(deviceEvent bluetooth.DeviceEventData, remove bool) {
	r.Lock()
	defer r.Unlock()

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

		return
	}

	deviceNode.noder.SetDeviceEventData(deviceEvent)
}

func (r *rootNode) updateAction(updateMsg actionUpdateMsg) {
	r.Lock()
	defer r.Unlock()

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
