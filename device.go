package gohome

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/log"
)

//TODO: Change to interface, make this private
type Device struct {
	ID             string
	Name           string
	Description    string
	System         *System
	ConnectionInfo comm.ConnectionInfo
	Buttons        map[string]*Button

	evpDone      chan bool
	evpFire      chan Event
	pool         comm.ConnectionPool
	cmdProcessor CommandProcessor
}

func NewDevice(id, name, description string, s *System, cp CommandProcessor) *Device {
	return &Device{
		ID:           id,
		Name:         name,
		Description:  description,
		System:       s,
		Buttons:      make(map[string]*Button),
		cmdProcessor: cp,
	}
}

func (d *Device) InitConnections() {
	createConnection := func() comm.Connection {
		conn := comm.NewTelnetConnection(d.ConnectionInfo)
		conn.SetPingCallback(func() error {
			if _, err := conn.Write([]byte("#PING\r\n")); err != nil {
				fmt.Printf("ping failed: %s", err.Error())
				return err
			}
			return nil
		})
		return conn
	}
	ps := d.ConnectionInfo.PoolSize
	log.V("%s init connections, pool size %d", d, ps)
	d.pool = comm.NewConnectionPool(d.Name, ps, createConnection)
	log.V("%s connected", d)
}

func (d *Device) Connect() (comm.Connection, error) {
	c := d.pool.Get()
	if c == nil {
		return nil, fmt.Errorf("%s - connect failed, no connection available", d)
	}
	return c, nil
}

func (d *Device) ReleaseConnection(c comm.Connection) {
	d.pool.Release(c)
}

func (d *Device) StartProducingEvents() (<-chan Event, <-chan bool) {
	d.evpDone = make(chan bool)
	d.evpFire = make(chan Event)

	if d.ConnectionInfo.Stream {
		go startStreaming(d)
	}
	return d.evpFire, d.evpDone
}

func (d *Device) Authenticate(c comm.Connection) error {
	r := bufio.NewReader(c)
	_, err := r.ReadString(':')
	if err != nil {
		return fmt.Errorf("authenticate login failed: %s", err)
	}

	info := c.Info()
	_, err = c.Write([]byte(info.Login + "\r\n"))
	if err != nil {
		return fmt.Errorf("authenticate write login failed: %s", err)
	}

	_, err = r.ReadString(':')
	if err != nil {
		return fmt.Errorf("authenticate password failed: %s", err)
	}
	_, err = c.Write([]byte(info.Password + "\r\n"))
	if err != nil {
		return fmt.Errorf("authenticate password failed: %s", err)
	}
	return nil
}

func (d *Device) String() string {
	return fmt.Sprintf("Device[%s]", d.Name)
}

func startStreaming(d *Device) {
	//TODO: Stop?
	for {
		err := stream(d)
		if err != nil {
			log.E("%s streaming failed: %s", d, err)
		}
		time.Sleep(10 * time.Second)
	}
}

func stream(d *Device) error {
	log.V("%s attemping to stream events", d)
	conn, err := d.Connect()
	if err != nil {
		return fmt.Errorf("%s unable to connect to stream events: %s", d, err)
	}

	defer func() {
		d.ReleaseConnection(conn)
	}()

	log.V("%s streaming events", d)
	scanner := bufio.NewScanner(conn)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		//Match first instance of ~OUTPUT|~DEVICE.*\r\n
		str := string(data[0:])
		indices := regexp.MustCompile("[~|#][OUTPUT|DEVICE].+\r\n").FindStringIndex(str)

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
			//TODO: How is ping getting through to here, if we are not scanning for it?
			orig := scanner.Text()
			if cmd := parseCommandString(d, orig); cmd != nil {
				d.evpFire <- NewEvent(d, cmd, orig, ETUnknown)
			}
		}
	}

	log.V("%s stopped streaming events", d)
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%s error streaming events, streaming stopped: %s", d, err)
	}
	return nil

	/*
		//TODO: When?
		if d.evpDone != nil {
			close(d.evpDone)
		}
	*/
}

func parseCommandString(d *Device, cmd string) Command {
	switch {
	case strings.HasPrefix(cmd, "~OUTPUT"),
		strings.HasPrefix(cmd, "#OUTPUT"):
		return parseZoneCommand(d, cmd)

	case strings.HasPrefix(cmd, "~DEVICE"),
		strings.HasPrefix(cmd, "#DEVICE"):
		return parseDeviceCommand(d, cmd)
	default:
		// Ignore commands we don't care about
		return nil
	}
}

type commandBuilderParams struct {
	CommandType  CommandType
	Zone         *Zone
	Intensity    float64
	Device       *Device
	SourceDevice *Device
	ComponentID  string
}

func buildCommand(p commandBuilderParams) Command {
	switch p.CommandType {
	case CTZoneSetLevel:
		return &StringCommand{
			Device:   p.Device,
			Friendly: fmt.Sprintf("Zone \"%s\" set to %.2f%%", p.Zone.Name, p.Intensity),
			Value:    fmt.Sprintf("#OUTPUT,%s,1,%.2f\r\n", p.Zone.ID, p.Intensity),
			Type:     p.CommandType,
		}

	case CTDevicePressButton:
		return &StringCommand{
			Device:   p.Device,
			Friendly: fmt.Sprintf("Device \"%s\" press button %s", p.SourceDevice.Name, p.ComponentID),
			Value:    fmt.Sprintf("#DEVICE,%s,%s,3\r\n", p.SourceDevice.Name, p.ComponentID),
			Type:     p.CommandType,
		}

	case CTDeviceReleaseButton:
		return &StringCommand{
			Device:   p.Device,
			Friendly: fmt.Sprintf("Device \"%s\" release button %s", p.SourceDevice.Name, p.ComponentID),
			Value:    fmt.Sprintf("#DEVICE,%s,%s,4\r\n", p.SourceDevice.Name, p.ComponentID),
			Type:     p.CommandType,
		}

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

	return buildCommand(commandBuilderParams{
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

	return buildCommand(commandBuilderParams{
		Device:      d,
		CommandType: ct,
		Intensity:   intensity,
		Zone:        z,
	})
}
