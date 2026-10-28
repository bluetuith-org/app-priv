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

//revive:disable
func nget(n *ctree.Node) adNoder {
	return n.Reference().(adNoder)
}

//revive:disable
func nvget[v any](n *ctree.Node) v {
	return n.Reference().(v)
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
