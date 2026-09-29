package adtree

import (
	"iter"
	"strings"

	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
)

// nodeIDNib is the prefix prepended to a node's identifier.
type nodeIDNib = string

// The different types of node ID nibs.
const (
	nibDevicesList nodeIDNib = "dl"
	nibActionsList nodeIDNib = "al"
	nibAction      nodeIDNib = "ac"
	nibAdapter     nodeIDNib = "ad"
	nibDevice      nodeIDNib = "dv"
)

const (
	lenNibPlusColon = 3
	lenBdAddr       = 12
	lenNodeIDTotal  = lenNibPlusColon + lenBdAddr
)

// nodeID holds the identifier information for a node.
type nodeID struct {
	id string

	nodeNib nodeIDNib
}

// newRootNodeID returns a node ID for the root node.
func newRootNodeID() nodeID {
	return nodeID{id: "root", nodeNib: ""}
}

// newAdapterNodeID returns a node ID from an adapter address.
func newAdapterNodeID(address bluetooth.AdapterAddress) nodeID {
	n := nodeID{
		nodeNib: nibAdapter,
		id: useStringBuffer(lenNodeIDTotal, func(b *strings.Builder) {
			b.WriteString(nibAdapter)
			b.WriteString(":")
			appendMacAddress(b, address.Address)
		}),
	}

	return n
}

// newDeviceNodeID returns a node ID from a device address.
func newDeviceNodeID(address bluetooth.DeviceAddress) nodeID {
	return nodeID{
		nodeNib: nibDevice,
		id: useStringBuffer((lenNodeIDTotal*2)+2+len(nibDevicesList), func(b *strings.Builder) {
			b.WriteString(nibAdapter)
			b.WriteString(":")
			appendMacAddress(b, address.AssociatedAdapter)

			b.WriteString("/")
			b.WriteString(nibDevicesList)
			b.WriteString("/")

			b.WriteString(nibDevice)
			b.WriteString(":")
			appendMacAddress(b, address.Address)
		}),
	}
}

// appendSubNodeNib appends a nib to the current ID.
func (n nodeID) appendSubNodeNib(nib nodeIDNib) nodeID {
	return n.appendSubNodeTextNib(nib, "")
}

// appendSubNodeTextNib appends a nib-value pair to the current ID.
func (n nodeID) appendSubNodeTextNib(nib nodeIDNib, text string) nodeID {
	length := len(n.id) + lenNibPlusColon + len(text)

	n.id = useStringBuffer(length, func(b *strings.Builder) {
		b.WriteString(n.id)
		b.WriteString("/")
		b.WriteString(nib)

		if text != "" {
			b.WriteString(":")
			b.WriteString(text)
		}
	})

	return n
}

// String returns the string representation of the node ID.
func (n nodeID) String() string {
	return n.id
}

// iterNodeID splits a full identifier for a node into individual parts.
func iterNodeID(id string) iter.Seq[string] {
	return func(yield func(string) bool) {
		currIdx := 0

		for currIdx < len(id) {
			idx := strings.Index(id[currIdx:], "/")
			if idx < 0 {
				break
			}

			currIdx += idx + 1
			nodeID := id[:currIdx-1]

			if !yield(nodeID) {
				return
			}
		}

		if !yield(id) {
			return
		}
	}
}

// appendMacAddress appends a MAC address to the string buffer.
func appendMacAddress(b *strings.Builder, mac bluetooth.MacAddress) {
	for i := 5; i >= 0; i-- {
		c := mac[i]
		nibble := c >> 4
		if nibble <= 9 {
			b.WriteByte(nibble + '0')
		} else {
			b.WriteByte(nibble + 'A' - 10)
		}

		nibble = c & 0x0f
		if nibble <= 9 {
			b.WriteByte(nibble + '0')
		} else {
			b.WriteByte(nibble + 'A' - 10)
		}
	}
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

func useStringBuffer(size int, fn func(b *strings.Builder)) string {
	var sb strings.Builder
	sb.Grow(size)

	fn(&sb)

	return sb.String()
}
