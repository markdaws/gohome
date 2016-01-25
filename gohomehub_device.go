package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/fluxwifi"
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
			createConnection := func() comm.Connection {
				conn := comm.NewTelnetConnection(z.Address, nil)
				//TODO: Need to get some ping mechanism
				/*
					conn.SetPingCallback(func() error {
						if _, err := conn.Write([]byte("#PING\r\n")); err != nil {
							return fmt.Errorf("%s ping failed: %s", d, err)
						}
						return nil
					})*/
				return conn
			}
			ps := 2
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

func (d *GoHomeHubDevice) StartProducingEvents() (<-chan event.Event, <-chan bool) {
	return nil, nil
}

func (d *GoHomeHubDevice) Authenticate(c comm.Connection) error {
	return nil
}

func (d *GoHomeHubDevice) BuildCommand(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		return d.buildZoneSetLevelCommand(command)
	case *cmd.ZoneTurnOn:
		return d.buildZoneTurnOnCommand(command)
	case *cmd.ZoneTurnOff:
		return d.buildZoneTurnOffCommand(command)
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

func (d *GoHomeHubDevice) buildZoneTurnOnCommand(c *cmd.ZoneTurnOn) (*cmd.Func, error) {
	z, ok := d.Zones()[c.ZoneAddress]
	if !ok {
		return nil, fmt.Errorf("unknown zone ID %s", c.ZoneID)
	}

	switch z.Controller {
	case zone.ZCFluxWIFI:
		return &cmd.Func{
			Func: func() error {
				pool, ok := d.pools[z.Controller]
				if !ok || pool == nil {
					return fmt.Errorf("gohomehub - connection pool not ready")
				}

				conn := pool.Get()
				if conn == nil {
					return fmt.Errorf("gohomehub - error connecting, no available connections")
				}

				defer func() {
					d.pools[z.Controller].Release(conn)
				}()
				return fluxwifi.TurnOn(conn)
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported zone controller")
	}
}

func (d *GoHomeHubDevice) buildZoneTurnOffCommand(c *cmd.ZoneTurnOff) (*cmd.Func, error) {
	z, ok := d.Zones()[c.ZoneAddress]
	if !ok {
		return nil, fmt.Errorf("unknown zone ID %s", c.ZoneID)
	}

	switch z.Controller {
	case zone.ZCFluxWIFI:
		return &cmd.Func{
			Func: func() error {
				pool, ok := d.pools[z.Controller]
				if !ok || pool == nil {
					return fmt.Errorf("gohomehub - connection pool not ready")
				}

				conn := pool.Get()
				if conn == nil {
					return fmt.Errorf("gohomehub - error connecting, no available connections")
				}

				defer func() {
					d.pools[z.Controller].Release(conn)
				}()
				return fluxwifi.TurnOff(conn)
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported zone controller")
	}
}

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

				pool, ok := d.pools[z.Controller]
				if !ok || pool == nil {
					return fmt.Errorf("gohomehub - connection pool not ready")
				}

				conn := pool.Get()
				if conn == nil {
					return fmt.Errorf("gohomehub - error connecting, no available connections")
				}

				defer func() {
					d.pools[z.Controller].Release(conn)
				}()
				return fluxwifi.SetLevel(rV, gV, bV, conn)
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported zone controller")
	}
}
