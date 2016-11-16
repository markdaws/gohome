package lutron

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-home-iot/event-bus"
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

	Stream(r io.Reader, handler func(Event)) error
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

type lbdgpro2whDevice struct{}

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

// Stream sits watching the stream of messages send by the device and parses them in to
// concrete event types. You provide a handler func that will receive all of the events. This
// function should return immediately and not block otherwise the stream will be blocked
func (d *lbdgpro2whDevice) Stream(r io.Reader, handler func(Event)) error {
	scanner := bufio.NewScanner(r)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		//Match first instance of ~OUTPUT|~DEVICE.*\r\n
		str := string(data[0:])
		indices := regexp.MustCompile("[~|#][OUTPUT|DEVICE].+\r\n").FindStringIndex(str)

		if indices != nil {
			token = []byte(string([]rune(str)[indices[0]:indices[1]]))
			advance = indices[1]
			err = nil
		} else {
			advance = 0
			token = nil
			err = nil
		}
		return
	}

	scanner.Split(split)

	//TODO: Context

	for scanner.Scan() {
		orig := scanner.Text()
		if evt := d.parseCommandString(orig); evt != nil {
			handler(evt)
		}
	}

	return scanner.Err()
}

func (d *lbdgpro2whDevice) parseCommandString(cmd string) Event {
	var evt Event
	switch {
	case strings.HasPrefix(cmd, "~OUTPUT"),
		strings.HasPrefix(cmd, "#OUTPUT"):
		evt = d.parseZoneCommand(cmd)

	case strings.HasPrefix(cmd, "~DEVICE"),
		strings.HasPrefix(cmd, "#DEVICE"):
		evt = d.parseDeviceCommand(cmd)
	}

	if evt == nil {
		return UnknownEvt{Msg: cmd}
	}
	return evt
}

func (d *lbdgpro2whDevice) parseZoneCommand(command string) Event {
	matches := regexp.MustCompile("[~|?]OUTPUT,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(command)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	zoneAddress := matches[1]
	cmdID := matches[2]
	level, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return nil
	}

	switch cmdID {
	case "1":
		return &ZoneLevelEvt{
			Address: zoneAddress,
			Level:   float32(level),
		}
	default:
		return nil
	}
}

func (d *lbdgpro2whDevice) parseDeviceCommand(command string) evtbus.Event {
	matches := regexp.MustCompile("[~|#]DEVICE,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(command)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	deviceAddr := matches[1]
	buttonAddr := matches[2]
	cmdID := matches[3]

	switch cmdID {
	case "3":
		return &BtnPressEvt{
			Address:       buttonAddr,
			DeviceAddress: deviceAddr,
		}
	case "4":
		return &BtnReleaseEvt{
			Address:       buttonAddr,
			DeviceAddress: deviceAddr,
		}
	default:
		return nil
	}
}

func sendString(cmd string, w io.Writer) error {
	_, err := w.Write([]byte(cmd))
	if err != nil {
		return fmt.Errorf("failed to send command \"%s\" %s\n", cmd, err)
	}
	return nil
}
