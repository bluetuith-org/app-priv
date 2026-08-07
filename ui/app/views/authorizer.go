package views

import (
	"errors"

	"github.com/bluetuith-org/bluetooth-classic/api/bluetooth"
	"github.com/google/uuid"
)

// authorizer holds a set of functions used to authenticate pairing and receiving
// file transfer requests. A new instance of this is supposed to be passed to
// [bluetooth.Session.Start] to handle any authorization requests.
type authorizer struct {
	v           *Views
	initialized bool

	alwaysAuthorize bool
}

// newAuthorizer returns a new authorizer.
func newAuthorizer(v *Views) *authorizer {
	return &authorizer{v: v}
}

// setInitialized sets the authorizer to the initialized state.
// This is called after all views have been initialized.
func (a *authorizer) setInitialized() {
	a.initialized = true
}

// AuthorizeTransfer asks the user to authorize a file transfer (Object Push) that is about to be sent
// from the remote device.
func (a *authorizer) AuthorizeTransfer(timeout bluetooth.AuthTimeout, props bluetooth.ObjectPushData) error {
	if !a.initialized {
		return nil
	}

	return errors.New("Cancelled")
}

// DisplayPinCode displays the pincode from the remote device to the user during a pairing authorization session.
func (a *authorizer) DisplayPinCode(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress, pincode string) error {
	if !a.initialized {
		return nil
	}

	// 	modal := a.generateDisplayModal(address, "pincode", "Pin Code", msg)
	// 	modal.display(timeout)

	return nil
}

// DisplayPasskey only displays the passkey from the remote device to the user during a pairing authorization session.
// This can be called multiple times, since each time the user enters a number on the remote device, this function
// is called with the updated 'entered' value.
// TODO: Handle multiple calls/draws when this function is called.
func (a *authorizer) DisplayPasskey(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress, passkey uint32, entered uint16) error {
	if !a.initialized {
		return nil
	}

	// 	modal := a.generateDisplayModal(address, "passkey-display", "Passkey Display", msg)
	// 	modal.display(timeout)

	return nil
}

// ConfirmPasskey asks the user to authorize the pairing request using the provided passkey.
func (a *authorizer) ConfirmPasskey(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress, passkey uint32) error {
	if !a.initialized {
		return nil
	}

	return nil
}

// AuthorizePairing asks the user to authorize a pairing request.
func (a *authorizer) AuthorizePairing(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress) error {
	if !a.initialized {
		return nil
	}

	return nil
}

// AuthorizeService asks the user to authorize whether a specific Bluetooth Profile is allowed to be used.
func (a *authorizer) AuthorizeService(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress, profileUUID uuid.UUID) error {
	if !a.initialized || a.alwaysAuthorize {
		return nil
	}

	return errors.New("Cancelled")
}

// TODO
// generateConfirmModal generates a confirmation modal with the provided parameters.
func (a *authorizer) generateConfirmModal(address bluetooth.MacAddress, name, title, msg string) {
}

// generateDisplayModal generates a display modal with the provided parameters.
func (a *authorizer) generateDisplayModal(address bluetooth.MacAddress, name, title, msg string) {
}
