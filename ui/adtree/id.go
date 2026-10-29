package adtree

import (
	"encoding/binary"

	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/bluetuith-org/bluetuith/keybindings"
)

// nodeID holds the identifier information for a node.
type nodeID uint64

// IsValid checks if the node ID is valid.
func (n nodeID) IsValid() bool {
	return n != 0
}

// ToUint64 converts the node ID to a uint64 value.
func (n nodeID) ToUint64() uint64 {
	return uint64(n)
}

// nibNodeID holds the nib that is added to the node ID.
type nibNodeID uint8

// The different types of node ID nibs.
const (
	nibRoot nibNodeID = iota
	nibDevicesList
	nibActionsList
	nibAction
)

// Add adds a node ID with the nib and returns the result.
func (n nibNodeID) Add(id nodeID) nodeID {
	return nodeID(n) + id
}

// AddIdentifiers adds a node ID and a set of identifiers with the nib
// and returns the result.
func (n nibNodeID) AddIdentifiers(id nodeID, nums ...uint64) nodeID {
	i := n.Add(id)
	for _, nm := range nums {
		i += nodeID(nm)
	}

	return i
}

// newRootNodeID returns a node ID for the root node.
func newRootNodeID() nodeID {
	return nodeID(nibRoot)
}

// newAdapterNodeID returns a node ID from an adapter address.
func newAdapterNodeID(address bluetooth.AdapterAddress) nodeID {
	return macAddress(address.Address)
}

func queryAdapter(address bluetooth.AdapterAddress) []nodeID {
	return []nodeID{macAddress(address.Address)}
}

// newDeviceNodeID returns a node ID from a device address.
func newDeviceNodeID(dlnode *dlNode, address bluetooth.DeviceAddress) nodeID {
	dv := macAddress(address.Address)

	return dlnode.id + dv
}

func queryDevice(address bluetooth.DeviceAddress) []nodeID {
	ids := queryDlist(bluetooth.NewAdapterAddress(address.AssociatedAdapter))
	ids = append(ids, ids[len(ids)-1]+macAddress(address.Address))

	return ids
}

func newDlNodeID(adapternode *adapterNode) nodeID {
	return nibDevicesList.Add(adapternode.id)
}

func queryDlist(address bluetooth.AdapterAddress) []nodeID {
	ad := macAddress(address.Address)
	return []nodeID{ad, nibDevicesList.Add(ad)}
}

func newAlNodeID(id nodeID) nodeID {
	return nibActionsList.Add(id)
}

func newAcNodeID(alnode *alNode, kn keybindings.KeyID) nodeID {
	return nibAction.AddIdentifiers(alnode.id, uint64(kn))
}

func macAddress(mac bluetooth.MacAddress) nodeID {
	var b [8]byte

	copy(b[2:], mac[:])

	return nodeID(binary.BigEndian.Uint64(b[:]))
}
