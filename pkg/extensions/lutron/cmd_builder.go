package lutron

import (
	"fmt"
	"io"
	"time"

	lutronExt "github.com/go-home-iot/lutron"
	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/feature"
	"github.com/markdaws/gohome/pkg/gohome"
)

type cmdBuilder struct {
	System *gohome.System
	device lutronExt.Device
}

func (b *cmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {

	switch command := c.(type) {
	case *cmd.FeatureSetAttrs:
		f := b.System.FeatureByID(command.FeatureID)
		if f == nil {
			return nil, fmt.Errorf("unknown feature ID: %s", command.FeatureID)
		}

		d := b.System.DeviceByID(f.DeviceID)
		if d == nil {
			return nil, fmt.Errorf("unknown device ID: %s", f.DeviceID)
		}

		// Lutron supports light zones or window treatments
		switch f.Type {
		case feature.FTLightZone:
			level, err := feature.LightZoneGetBrightness(command.Attrs)
			if err != nil {
				return nil, err
			}

			return &cmd.Func{
				Func: func() error {
					return getWriterAndExec(d, func(d lutronExt.Device, w io.Writer) error {
						return d.SetLevel(level, f.Address, w)
					})
				},
			}, nil

		case feature.FTWindowTreatment:
			level, err := feature.WindowTreatmentGetOffset(command.Attrs)
			if err != nil {
				return nil, err
			}

			return &cmd.Func{
				Func: func() error {
					return getWriterAndExec(d, func(d lutronExt.Device, w io.Writer) error {
						return d.SetLevel(level, f.Address, w)
					})
				},
			}, nil
		}

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
	return nil, nil
}

func getWriterAndExec(d *gohome.Device, f func(lutronExt.Device, io.Writer) error) error {
	var hub *gohome.Device = d
	if d.Hub != nil {
		hub = d.Hub
	}

	conn, err := hub.Connections.Get(time.Second*5, true)
	if err != nil {
		return fmt.Errorf("error connecting, pool returned err: %s", err)
	}

	lDev, err := lutronExt.DeviceFromModelNumber(hub.ModelNumber)
	if err != nil {
		return err
	}

	err = f(lDev, conn)
	hub.Connections.Release(conn, err)
	if err != nil {
		return fmt.Errorf("Failed to send command %s\n", err)
	}
	return nil
}
