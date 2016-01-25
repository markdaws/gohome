package gohome

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/zone"
)

//TODO: make lutron package
type Lbdgpro2whDevice struct {
	device
	pool comm.ConnectionPool
}

func (d *Lbdgpro2whDevice) ModelNumber() string {
	return "L-BDGPRO2-WH"
}

func (d *Lbdgpro2whDevice) InitConnections() {
	createConnection := func() comm.Connection {
		conn := comm.NewTelnetConnection(d.Address(), d.Auth().Authenticator)
		conn.SetPingCallback(func() error {
			if _, err := conn.Write([]byte("#PING\r\n")); err != nil {
				return fmt.Errorf("%s ping failed: %s", d, err)
			}
			return nil
		})
		return conn
	}
	ps := 2
	log.V("%s init connections, pool size %d", d, ps)
	d.pool = comm.NewConnectionPool(d.name, ps, createConnection)
	log.V("%s connected", d)
}

func (d *Lbdgpro2whDevice) StartProducingEvents() (<-chan event.Event, <-chan bool) {
	d.evpDone = make(chan bool)
	d.evpFire = make(chan event.Event)

	if d.Stream() {
		go d.startStreaming()
	}
	return d.evpFire, d.evpDone
}

func (d *Lbdgpro2whDevice) Authenticate(c comm.Connection) error {
	r := bufio.NewReader(c)
	_, err := r.ReadString(':')
	if err != nil {
		return fmt.Errorf("authenticate login failed: %s", err)
	}

	_, err = c.Write([]byte(d.auth.Login + "\r\n"))
	if err != nil {
		return fmt.Errorf("authenticate write login failed: %s", err)
	}

	_, err = r.ReadString(':')
	if err != nil {
		return fmt.Errorf("authenticate password failed: %s", err)
	}

	_, err = c.Write([]byte(d.auth.Password + "\r\n"))
	if err != nil {
		return fmt.Errorf("authenticate password failed: %s", err)
	}
	return nil
}

func (d *Lbdgpro2whDevice) Connect() (comm.Connection, error) {
	c := d.pool.Get()
	if c == nil {
		return nil, fmt.Errorf("%s - connect failed, no connection available", d)
	}
	return c, nil
}

func (d *Lbdgpro2whDevice) ReleaseConnection(c comm.Connection) {
	d.pool.Release(c)
}

func (d *Lbdgpro2whDevice) BuildCommand(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		return &cmd.Func{
			Func: func() error {
				newCmd := &StringCommand{
					Device: d,
					Value:  "#OUTPUT," + command.ZoneAddress + ",1,%.2f\r\n",
					Args:   []interface{}{command.Level.Value},
				}
				return newCmd.Execute()
			},
		}, nil
	case *cmd.ZoneTurnOn:
		return &cmd.Func{
			Func: func() error {
				newCmd := &StringCommand{
					Device: d,
					Value:  "#OUTPUT," + command.ZoneAddress + ",1,%.2f\r\n",
					Args:   []interface{}{100.0},
				}
				return newCmd.Execute()
			},
		}, nil
	case *cmd.ZoneTurnOff:
		return &cmd.Func{
			Func: func() error {
				newCmd := &StringCommand{
					Device: d,
					Value:  "#OUTPUT," + command.ZoneAddress + ",1,%.2f\r\n",
					Args:   []interface{}{0.0},
				}
				return newCmd.Execute()
			},
		}, nil
	case *cmd.ButtonPress:
		return &cmd.Func{
			Func: func() error {
				newCmd := &StringCommand{
					Device: d,
					Value:  "#DEVICE," + command.DeviceAddress + "," + command.ButtonAddress + ",3\r\n",
				}
				return newCmd.Execute()
			},
		}, nil

	case *cmd.ButtonRelease:
		return &cmd.Func{
			Func: func() error {
				cmd := &StringCommand{
					Device: d,
					Value:  "#DEVICE," + command.DeviceAddress + "," + command.ButtonAddress + ",4\r\n",
				}
				return cmd.Execute()
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
}

func (d *Lbdgpro2whDevice) startStreaming() {
	//TODO: Stop?
	for {
		err := d.stream()
		if err != nil {
			log.E("%s streaming failed: %s", d, err)
		}
		time.Sleep(10 * time.Second)
	}
}

func (d *Lbdgpro2whDevice) stream() error {
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
			if cmd := d.parseCommandString(orig); cmd != nil {
				d.evpFire <- event.New(d.ID(), cmd, orig)
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

func (d *Lbdgpro2whDevice) parseCommandString(cmd string) cmd.Command {
	switch {
	case strings.HasPrefix(cmd, "~OUTPUT"),
		strings.HasPrefix(cmd, "#OUTPUT"):
		return d.parseZoneCommand(cmd)

	case strings.HasPrefix(cmd, "~DEVICE"),
		strings.HasPrefix(cmd, "#DEVICE"):
		return d.parseDeviceCommand(cmd)
	default:
		// Ignore commands we don't care about
		return nil
	}
}

func (d *Lbdgpro2whDevice) parseDeviceCommand(command string) cmd.Command {
	matches := regexp.MustCompile("[~|#]DEVICE,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(command)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	deviceID := matches[1]
	componentID := matches[2]
	cmdID := matches[3]
	sourceDevice := d.Devices()[deviceID]
	if sourceDevice == nil {
		return nil
	}

	var finalCmd cmd.Command
	switch cmdID {
	case "3":
		if btn := sourceDevice.Buttons()[componentID]; btn != nil {
			finalCmd = &cmd.ButtonPress{
				ButtonAddress: btn.Address,
				ButtonID:      btn.ID,
				DeviceName:    d.Name(),
				DeviceAddress: d.Address(),
				DeviceID:      d.ID(),
			}
		}
	case "4":
		if btn := sourceDevice.Buttons()[componentID]; btn != nil {
			finalCmd = &cmd.ButtonRelease{
				ButtonAddress: btn.Address,
				ButtonID:      btn.ID,
				DeviceName:    d.Name(),
				DeviceAddress: d.Address(),
				DeviceID:      d.ID(),
			}
		}
	default:
		return nil
	}

	return finalCmd
}

func (d *Lbdgpro2whDevice) parseZoneCommand(command string) cmd.Command {
	matches := regexp.MustCompile("[~|?]OUTPUT,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(command)
	if matches == nil || len(matches) != 4 {
		return nil
	}

	zoneID := matches[1]
	cmdID := matches[2]
	level, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return nil
	}

	z := d.Zones()[zoneID]
	if z == nil {
		return nil
	}

	var finalCmd cmd.Command
	switch cmdID {
	case "1":
		finalCmd = &cmd.ZoneSetLevel{
			ZoneAddress: z.Address,
			ZoneID:      z.ID,
			ZoneName:    z.Name,
			Level:       cmd.Level{Value: float32(level)},
		}
	default:
		return nil
	}

	return finalCmd
}

func (d *Lbdgpro2whDevice) SupportsController(c zone.Controller) bool {
	return c == zone.ZCDefault
}
