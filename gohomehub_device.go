package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/zone"
)

type GoHomeHubDevice struct {
	device
	pools map[string]comm.ConnectionPool
}

func (d *GoHomeHubDevice) ModelNumber() string {
	return "GoHomeHub"
}

func (d *GoHomeHubDevice) InitConnections() {

	d.pools = make(map[string]comm.ConnectionPool)

	// May have multiple zones + devices we need to talk to
	// set up pools for each one
	for _, z := range d.Zones() {
		switch z.Controller {
		case zone.ZCFluxWIFI:
			ci := comm.TelnetConnectionInfo{
				PoolSize: 2,
				Network:  "tcp",
				Address:  z.Address,
			}
			createConnection := func() comm.Connection {
				conn := comm.NewTelnetConnection(ci)
				/*
					conn.SetPingCallback(func() error {
						if _, err := conn.Write([]byte("#PING\r\n")); err != nil {
							return fmt.Errorf("%s ping failed: %s", d, err)
						}
						return nil
					})*/
				return conn
			}
			ps := ci.PoolSize
			log.V("%s init connections, pool size %d", d, ps)
			d.pools[z.Controller] = comm.NewConnectionPool(d.name, ps, createConnection)
			log.V("%s connected", d)
		}
	}
}

func (d *GoHomeHubDevice) Connect() (comm.Connection, error) {
	return nil, fmt.Errorf("unsupported function connect")
}

func (d *GoHomeHubDevice) ReleaseConnection(c comm.Connection) {
}

func (d *GoHomeHubDevice) StartProducingEvents() (<-chan Event, <-chan bool) {
	return nil, nil
}

func (d *GoHomeHubDevice) Authenticate(c comm.Connection) error {
	return nil
}

func (d *GoHomeHubDevice) BuildCommand(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		return d.buildZoneSetLevelCommand(command)
	case *cmd.ButtonPress:
		//TODO: Phantom buttons?
		return nil, fmt.Errorf("goHomeHubDevice ButtonPressCommand not supported")
	case *cmd.ButtonRelease:
		return nil, fmt.Errorf("goHomeHubDevice ButtonReleaseCommand not supported")
	case *cmd.SceneSet:
		//TODO: Does this make sense, what does a scene mean in terms of this virtual hub?
	default:
		return nil, fmt.Errorf("goHomeHubDevice build commands not supported")
	}

	return nil, fmt.Errorf("goHomeHubDevice unsupported command")
}

//TODO: Level should be a type with value,r,g,b, not just one value
func (d *GoHomeHubDevice) buildZoneSetLevelCommand(c *cmd.ZoneSetLevel) (*cmd.Func, error) {

	z, ok := d.Zones()[c.ZoneAddress]
	if !ok {
		return nil, fmt.Errorf("unknown zone ID %s", c.ZoneID)
	}

	switch z.Controller {
	case zone.ZCFluxWIFI:
		return &cmd.Func{
			Func: func() error {

				var rV, gV, bV byte
				lvl := c.Level.Value
				if lvl == 0 {
					if (c.Level.R == 0) && (c.Level.G == 0) && (c.Level.B == 0) {
						rV = 0
						gV = 0
						bV = 0
					} else {
						rV = c.Level.R
						gV = c.Level.G
						bV = c.Level.B
					}
				} else {
					rV = byte((lvl / 100) * 255)
					gV = rV
					bV = rV
				}

				b := []byte{0x31, rV, gV, bV, 0x00, 0xf0, 0x0f}
				var t int = 0
				for _, v := range b {
					t += int(v)
				}
				cs := t & 0xff
				b = append(b, byte(cs))

				conn := d.pools[z.Controller].Get()
				if conn == nil {
					return fmt.Errorf("gohomehub - error connecting, no available connections")
				}

				defer func() {
					d.pools[z.Controller].Release(conn)
				}()
				_, err := conn.Write(b)
				if err != nil {
					fmt.Printf("ERROR SENDING %s\n", err)
				} else {
				}
				return err
			},
		}, nil
	default:
		fmt.Println(z.Controller)
		return nil, fmt.Errorf("unsupported zone controller")
	}
}
