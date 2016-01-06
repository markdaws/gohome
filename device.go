package gohome

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

//TODO: Change to interface
type Device struct {
	Identifiable
	System         *System
	ConnectionInfo ConnectionInfo
	Buttons        map[string]*Button

	evpDone chan bool
	evpFire chan Event
	conn    Connection
}

func (d *Device) Connect() (Connection, error) {
	if d.conn != nil {
		//TODO: Fix real connection pool
		//TODO: Detect closed connections
		return d.conn, nil
	}

	//TODO: What is default timeout of net.Conn, change
	//TODO: When we support more than telnet, will use ConnectionInfo to
	//determine what type of connection we need to make
	conn := &TelnetConnection{}
	conn.SetConnectionInfo(d.ConnectionInfo)
	err := conn.Open()
	if err != nil {
		return nil, err
	}

	d.conn = conn
	return conn, nil
}

func (d *Device) StartProducingEvents() (<-chan Event, <-chan bool) {
	//TODO: When to init these
	d.evpDone = make(chan bool)
	d.evpFire = make(chan Event)

	if conn, err := d.Connect(); err != nil {
		fmt.Println("Failed to connect to device")
	} else {
		fmt.Println("Connected to device")
		go func() {
			stream(d, conn)
		}()
	}
	//TODO: Should return error

	return d.evpFire, d.evpDone
}

//TODO: How to stop this?
func stream(d *Device, r io.Reader) {
	scanner := bufio.NewScanner(r)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		//Match first instance of ~OUTPUT|~DEVICE.*\r\n
		str := string(data[0:])
		//fmt.Println(str)
		indices := regexp.MustCompile("[~|#][OUTPUT|DEVICE].+\r\n").FindStringIndex(str)
		//fmt.Printf("%s === %v\n", str, indices)

		//TODO: Don't let input grow forever - remove beginning chars after reaching max length

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
	for scanner.Scan() {
		//fmt.Printf("scanner: %s\n", scanner.Text())

		if d.evpFire != nil {
			orig := scanner.Text()
			cmd := ParseCommandString(d, orig)
			d.evpFire <- NewEvent(d, cmd, orig)
		}
	}

	if d.evpDone != nil {
		close(d.evpDone)
	}

	//TODO: What about connect/disconnect event
	fmt.Println("Done scanning")
	if err := scanner.Err(); err != nil {
		fmt.Printf("something happened", err)
	}
}

//TODO: Don't export
func ParseCommandString(d *Device, cmd string) Command {
	switch {
	case strings.HasPrefix(cmd, "~OUTPUT"),
		strings.HasPrefix(cmd, "#OUTPUT"):
		return parseZoneCommand(d, cmd)

	case strings.HasPrefix(cmd, "~DEVICE"),
		strings.HasPrefix(cmd, "#DEVICE"):
		return parseDeviceCommand(d, cmd)

	default:
		return nil
	}
}

func parseDeviceCommand(d *Device, cmd string) Command {
	matches := regexp.MustCompile("[~|#]DEVICE,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(cmd)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	deviceID := matches[1]
	componentID := matches[2]
	cmdID := matches[3]
	sourceDevice := d.System.Devices[deviceID]
	if sourceDevice == nil {
		//TODO: Error? Warning?
		return nil
	}

	var ct CommandType
	switch cmdID {
	case "3":
		ct = CTDevicePressButton
	case "4":
		ct = CTDeviceReleaseButton
	default:
		ct = CTUnknown
	}

	return BuildCommand(CommandBuilderParams{
		Device:       d,
		CommandType:  ct,
		SourceDevice: sourceDevice,
		ComponentID:  componentID,
	})
}

func parseZoneCommand(d *Device, cmd string) Command {
	matches := regexp.MustCompile("[~|?]OUTPUT,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(cmd)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	zoneID := matches[1]
	cmdID := matches[2]
	intensity, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		//TODO: Error
		return nil
	}

	//TODO: Get unique id based on device
	z := d.System.Zones[d.ID+":"+zoneID]
	if z == nil {
		//TODO: Error log
		return nil
	}

	var ct CommandType
	switch cmdID {
	case "1":
		ct = CTZoneSetLevel
	default:
		ct = CTUnknown
	}

	return BuildCommand(CommandBuilderParams{
		Device:      d,
		CommandType: ct,
		Intensity:   intensity,
		Zone:        z,
	})
}
