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

type Lbdgpro2whDevice struct {
	device
}

func (d *Lbdgpro2whDevice) ModelNumber() string {
	return "L-BDGPRO2-WH"
}

func (d *Lbdgpro2whDevice) InitConnections() {
	ci := *d.connectionInfo.(*comm.TelnetConnectionInfo)
	createConnection := func() comm.Connection {
		conn := comm.NewTelnetConnection(ci)
		conn.SetPingCallback(func() error {
			//TODO: Should return a command that then gets set on the command queue?
			if _, err := conn.Write([]byte("#PING\r\n")); err != nil {
				return fmt.Errorf("%s ping failed: %s", d, err)
			}
			return nil
		})
		return conn
	}
	ps := ci.PoolSize
	log.V("%s init connections, pool size %d", d, ps)
	d.pool = comm.NewConnectionPool(d.name, ps, createConnection)
	log.V("%s connected", d)
}

func (d *Lbdgpro2whDevice) StartProducingEvents() (<-chan Event, <-chan bool) {
	d.evpDone = make(chan bool)
	d.evpFire = make(chan Event)

	if d.Stream() {
		go startStreaming(d)
	}
	return d.evpFire, d.evpDone
}

func (d *Lbdgpro2whDevice) Authenticate(c comm.Connection) error {
	r := bufio.NewReader(c)
	_, err := r.ReadString(':')
	if err != nil {
		return fmt.Errorf("authenticate login failed: %s", err)
	}

	info := c.Info().(comm.TelnetConnectionInfo)
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

func (d *Lbdgpro2whDevice) BuildCommand(c Command) (*FCommand, error) {
	switch cmd := c.(type) {
	case *ZoneSetLevelCommand:
		return &FCommand{
			Func: func() error {
				cmd := &StringCommand{
					Device:   d,
					Value:    "#OUTPUT," + cmd.Zone.LocalID + ",1,%.2f\r\n",
					Friendly: "//TODO: Friendly",
					Args:     []interface{}{cmd.Level},
				}
				return cmd.Execute()
			},
		}, nil
	case *ButtonPressCommand:
		return &FCommand{
			Func: func() error {
				cmd := &StringCommand{
					Device:   d,
					Value:    "#DEVICE," + cmd.Button.Device.LocalID() + "," + cmd.Button.LocalID + ",3\r\n",
					Friendly: "//TODO: Friendly",
				}
				return cmd.Execute()
			},
		}, nil

	case *ButtonReleaseCommand:
		return &FCommand{
			Func: func() error {
				cmd := &StringCommand{
					Device:   d,
					Value:    "#DEVICE," + cmd.Button.Device.LocalID() + "," + cmd.Button.LocalID + ",4\r\n",
					Friendly: "//TODO: Friendly",
				}
				return cmd.Execute()
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
}

func startStreaming(d *Lbdgpro2whDevice) {
	//TODO: Stop?
	for {
		err := stream(d)
		if err != nil {
			log.E("%s streaming failed: %s", d, err)
		}
		time.Sleep(10 * time.Second)
	}
}

func stream(d *Lbdgpro2whDevice) error {
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
		if d.evpFire != nil {
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

func parseCommandString(d *Lbdgpro2whDevice, cmd string) Command {
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
	Zone         *Zone
	Intensity    float64
	Device       Device
	SourceDevice Device
	Button       *Button
}

func parseDeviceCommand(d *Lbdgpro2whDevice, cmd string) Command {
	matches := regexp.MustCompile("[~|#]DEVICE,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(cmd)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	deviceID := matches[1]
	componentID := matches[2]
	cmdID := matches[3]
	sourceDevice := d.Devices()[deviceID]
	if sourceDevice == nil {
		fmt.Printf("no source device %s\n", deviceID)
		//TODO: Error? Warning?
		return nil
	}

	var finalCmd Command
	switch cmdID {
	case "3":
		btn := sourceDevice.Buttons()[componentID]
		finalCmd = &ButtonPressCommand{
			Button: btn,
		}
	case "4":
		btn := sourceDevice.Buttons()[componentID]
		finalCmd = &ButtonReleaseCommand{
			Button: btn,
		}
	default:
		return nil
	}

	return finalCmd
}

func parseZoneCommand(d *Lbdgpro2whDevice, cmd string) Command {
	matches := regexp.MustCompile("[~|?]OUTPUT,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(cmd)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	zoneID := matches[1]
	cmdID := matches[2]
	level, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		//TODO: Error
		return nil
	}

	z := d.Zones()[zoneID]
	if z == nil {
		//TODO: Error log
		return nil
	}

	var finalCmd Command
	switch cmdID {
	case "1":
		//set level
		finalCmd = &ZoneSetLevelCommand{
			Zone:  z,
			Level: float32(level),
		}
	default:
		return nil
	}

	return finalCmd
}
