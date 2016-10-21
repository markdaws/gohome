package lutron

import (
	"fmt"
	"io"
)

// Device represents an interface to a Lutron device.  Different Lutron devices may
// send different commands, so this interface is used to abstract that from the callers
type Device interface {
	// SetLevel sets the device to the specified level
	SetLevel(level float32, zoneAddress string, w io.Writer) error

	// ButtonPress sends a button press command
	ButtonPress(devAddr, btnAddr string, w io.Writer) error

	// ButtonRelease sends a button release command
	ButtonRelease(devAddr, btnAddr string, w io.Writer) error
}

// DeviceFromModelNumber returns a Lutron device, based on the modelNumber parameter
func DeviceFromModelNumber(modelNumber string) (Device, error) {
	switch modelNumber {
	case "l-bdgpro2-wh":
		return &lbdgpro2whDevice{}, nil
	default:
		return nil, fmt.Errorf("unsupported model number: %s", modelNumber)
	}
}

type lbdgpro2whDevice struct {
}

func (d *lbdgpro2whDevice) SetLevel(level float32, zoneAddress string, w io.Writer) error {
	cmd := fmt.Sprintf("#OUTPUT,"+zoneAddress+",1,%.2f\r\n", level)
	_, err := w.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send command \"%s\" %s\n", cmd, err)
	}
	return nil
}

func (d *lbdgpro2whDevice) ButtonPress(devAddr, btnAddr string, w io.Writer) error {
	cmd := fmt.Sprintf("#DEVICE," + devAddr + "," + btnAddr + ",3\r\n")
	_, err := w.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send command \"%s\" %s\n", cmd, err)
	}
	return nil
}

func (d *lbdgpro2whDevice) ButtonRelease(devAddr, btnAddr string, w io.Writer) error {
	cmd := fmt.Sprintf("#DEVICE," + devAddr + "," + btnAddr + ",4\r\n")
	_, err := w.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send command \"%s\" %s\n", cmd, err)
	}
	return nil
}
