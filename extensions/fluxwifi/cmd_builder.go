package fluxwifi

import (
	"fmt"
	"time"

	"github.com/go-home-iot/connection-pool"
	fluxwifiExt "github.com/go-home-iot/fluxwifi"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type cmdBuilder struct {
	System *gohome.System
}

func getConnAndExecute(d *gohome.Device, f func(*pool.Connection) error) error {
	conn, err := d.Connections.Get(time.Second * 5)
	if err != nil {
		return fmt.Errorf("fluxwifiCmdBuilder - error connecting, no available connections")
	}

	defer func() {
		d.Connections.Release(conn)
	}()
	return f(conn)
}

func (b *cmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneTurnOn:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return getConnAndExecute(d, func(conn *pool.Connection) error {
					return fluxwifiExt.TurnOn(conn)
				})
			},
		}, nil

	case *cmd.ZoneTurnOff:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return getConnAndExecute(d, func(conn *pool.Connection) error {
					return fluxwifiExt.TurnOff(conn)
				})
			},
		}, nil

	case *cmd.ZoneSetLevel:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				var rV, gV, bV byte
				lvl := command.Level.Value
				if lvl == 0 {
					if (command.Level.R == 0) && (command.Level.G == 0) && (command.Level.B == 0) {
						rV = 0
						gV = 0
						bV = 0
					} else {
						rV = command.Level.R
						gV = command.Level.G
						bV = command.Level.B
					}
				} else {
					rV = byte((lvl / 100) * 255)
					gV = rV
					bV = rV
				}

				return getConnAndExecute(d, func(conn *pool.Connection) error {
					return fluxwifiExt.SetLevel(rV, gV, bV, conn)
				})
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
}
