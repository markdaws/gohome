package fluxwifi

import (
	"fmt"
	"time"

	"github.com/go-home-iot/connection-pool"
	fluxwifiExt "github.com/go-home-iot/fluxwifi"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/cmd"
)

type cmdBuilder struct {
	System *gohome.System
}

func getConnAndExecute(d *gohome.Device, f func(*pool.Connection) error) error {
	conn, err := d.Connections.Get(time.Second*5, true)
	if err != nil {
		return fmt.Errorf("fluxwifiCmdBuilder - error connecting, no available connections")
	}

	err = f(conn)
	d.Connections.Release(conn, err)
	return err
}

func (b *cmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.FeatureSetAttrs:
		f, ok := b.System.Features[command.FeatureID]
		if !ok {
			return nil, fmt.Errorf("unknown feature ID: %s", command.FeatureID)
		}

		d, ok := b.System.Devices[f.DeviceID]
		if !ok {
			return nil, fmt.Errorf("unknown device ID: %s", f.DeviceID)
		}

		for _, attribute := range command.Attrs {
			attribute := attribute

			switch attribute.Type {
			case attr.ATOnOff:
				if attribute.Value.(int32) == attr.OnOffOff {
					return &cmd.Func{
						Func: func() error {
							return getConnAndExecute(d, func(conn *pool.Connection) error {
								return fluxwifiExt.TurnOff(conn)
							})
						},
					}, nil
				} else {
					return &cmd.Func{
						Func: func() error {
							return getConnAndExecute(d, func(conn *pool.Connection) error {
								return fluxwifiExt.TurnOn(conn)
							})
						},
					}, nil
				}
			case attr.ATBrightness:
				return &cmd.Func{
					Func: func() error {
						return fmt.Errorf("unsupported")
						//TODO: Fix
						/*
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
							})*/
					},
				}, nil
			}
		}

		return nil, fmt.Errorf("unsupported attribute")

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
}
