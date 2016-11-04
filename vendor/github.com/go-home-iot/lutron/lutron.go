package lutron

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// Device represents an interface to a Lutron device.  Different Lutron devices may
// send different commands, so this interface is used to abstract that from the callers
type Device interface {
	// SetLevel sets the device to the specified level
	SetLevel(level float32, zoneAddress string, w io.Writer) error

	// RequestLevel requests the current level of the specified zone. You will have to
	// watch the devices stream to parse the response that comes back. It is async so
	// it may take time depending on how fast the lutron hub responds
	RequestLevel(zoneAddr string, w io.Writer) error

	// ButtonPress sends a button press command
	ButtonPress(devAddr, btnAddr string, w io.Writer) error

	// ButtonRelease sends a button release command
	ButtonRelease(devAddr, btnAddr string, w io.Writer) error

	// Ping sends a ping request e.g. #PING, callers can then wait for the
	// ping response e.g. ~PING, if nothing returns after a while you know something
	// is wrong with the connection you have to the hub
	Ping(w net.Conn, wait time.Duration) error
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

// Ping sends a ping command
func (d *lbdgpro2whDevice) Ping(rw net.Conn, wait time.Duration) error {
	err := sendString("#PING\r\n", rw)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(rw)
	endTime := time.Now().Add(wait)
	for {
		rw.SetReadDeadline(time.Now().Add(endTime.Sub(time.Now())))
		bytes, _, err := reader.ReadLine()
		rw.SetReadDeadline(time.Time{})
		if err != nil {
			return err
		}

		resp := string(bytes)
		if strings.Contains(resp, "~PING") {
			return nil
		}

		if time.Now().After(endTime) {
			return errors.New("timeout, no ping received")
		}
	}
}

// SetLevel request to set the level on the specified zone
func (d *lbdgpro2whDevice) SetLevel(level float32, zoneAddr string, w io.Writer) error {
	return sendString(fmt.Sprintf("#OUTPUT,%s,1,%.2f\r\n", zoneAddr, level), w)
}

// RequestLevel sends a level request for the specified zone, you will have to read the stream
// for the response.  e.g this sends ?OUTPUT,2,1 then async there will be a ~OUTPUT,2,1,50.00
// sent back by the lutron hub
func (d *lbdgpro2whDevice) RequestLevel(zoneAddr string, w io.Writer) error {
	return sendString(fmt.Sprintf("?OUTPUT,%s,1\r\n", zoneAddr), w)
}

// ButtonPress sends a button press command
func (d *lbdgpro2whDevice) ButtonPress(devAddr, btnAddr string, w io.Writer) error {
	return sendString(fmt.Sprintf("#DEVICE,%s,%s,3\r\n", devAddr, btnAddr), w)
}

// ButtonRelease sends a button release command
func (d *lbdgpro2whDevice) ButtonRelease(devAddr, btnAddr string, w io.Writer) error {
	return sendString(fmt.Sprintf("#DEVICE,%s,%s,4\r\n", devAddr, btnAddr), w)
}

func sendString(cmd string, w io.Writer) error {
	_, err := w.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send command \"%s\" %s\n", cmd, err)
	}
	return nil
}
