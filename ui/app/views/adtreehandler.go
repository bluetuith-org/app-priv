package views

import (
	"strings"

	"github.com/Digital-Shane/treeview/v2"
	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
)

type adTreeNode struct {
	id nodeID

	nodeType adTreeNodeType
	node     *treeview.Node[*adTreeNode]
}

func newAdTreeNode() *adTreeNode {
	data := &adTreeNode{}

	return data
}

func (a *adTreeNode) getID() nodeID {
	return a.id
}

type nodeIDBuilder struct {
	b *strings.Builder
}

// action/adapter/<adaddr>
// action/device/<adaddr>/<dvaddr>
// actionslist/adapter/<adaddr>
// actionslist/device/<adaddr>/<dvaddr>
// deviceslist/adapter/<adaddr>
// adapter/<adaddr>
// device/<adaddr>/<dvaddr>
func (a *adTreeNode) buildID(c ...string) {
	a.id = strings.Join(c, "/")
}

func (a *adTreeNode) withRootNode() *adTreeNode {
	a.buildID("root")

	a.nodeType = nodeTypeRoot

	a.node = treeview.NewNode(a.id, "Adapters", a)
	a.node.SetExpanded(true)

	return a
}

func (a *adTreeNode) withAdapter(adapter bluetooth.AdapterData) *adTreeNode {
	a.buildID(adapterIDNib, adapter.Address.String())

	a.nodeType = nodeTypeAdapter

	a.node = treeview.NewNode(a.id, getAdapterDisplayName(adapter), a)
	a.node.SetExpanded(true)

	return a
}

func (a *adTreeNode) withDevice(device bluetooth.DeviceData) *adTreeNode {
	a.buildID(deviceIDNib, device.AssociatedAdapter.String(), device.Address.String())

	a.nodeType = nodeTypeDevice

	a.node = treeview.NewNode(a.id, getDeviceDisplayName(device.DeviceEventData), a)
	a.node.SetExpanded(true)

	return a
}

func (a *adTreeNode) withDevicesList(id string) *adTreeNode {
	a.buildID(devicesListIDNib, id)

	a.nodeType = nodeTypeDevicesList

	a.node = treeview.NewNode(a.id, "Devices", a)
	a.node.SetExpanded(true)

	return a
}

func (a *adTreeNode) withAction(id, actionName string) *adTreeNode {
	a.buildID(actionIDNib, id)

	a.nodeType = nodeTypeAction

	a.node = treeview.NewNode(a.id, actionName, a)
	a.node.SetExpanded(false)

	return a
}

func (a *adTreeNode) withActionsList(id string) *adTreeNode {
	a.buildID(actionsListIDNib, id)

	a.nodeType = nodeTypeActionsList

	a.node = treeview.NewNode(a.id, "Actions", a)
	a.node.SetExpanded(false)

	return a
}
