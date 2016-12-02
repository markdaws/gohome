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
			case attr.ATHSL:
				return &cmd.Func{
					Func: func() error {
						hsl := attribute.Value.(string)
						r, g, bVal, err := attr.HSLStringToRGB(hsl)
						if err != nil {
							return fmt.Errorf("unable to read HSL value: %s", hsl)
						}

						return getConnAndExecute(d, func(conn *pool.Connection) error {
							err := fluxwifiExt.SetLevel(byte(r), byte(g), byte(bVal), conn)

							// For this particular piece of hardware it takes some time for it
							// to report the correct value after it actually gets set on the
							// hardware, so we inform the system that these values are the
							// new values and we should ignore any FeatureReportingEvt events from this
							// feature for the next 30 seconds, so that we don't report back incorrect
							// values to the users
							if err == nil {
								gohome.SupressFeatureReporting(b.System, command.FeatureID, command.Attrs, time.Second*30)
							}
							return err
						})
					},
				}, nil
			}
		}

		return nil, fmt.Errorf("unsupported attribute")

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
}
