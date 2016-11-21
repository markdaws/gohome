package honeywell

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/feature"

	honeywellExt "github.com/go-home-iot/honeywell"
)

type cmdBuilder struct {
	System *gohome.System
	Device *gohome.Device
}

func (b *cmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {

	command, ok := c.(*cmd.FeatureSetAttrs)
	if !ok {
		return nil, fmt.Errorf("unsupported command type")
	}

	f, ok := b.System.Features[command.FeatureID]
	if !ok {
		return nil, fmt.Errorf("unknown feature ID: %s", command.FeatureID)
	}

	d, ok := b.System.Devices[f.DeviceID]
	if !ok {
		return nil, fmt.Errorf("unknown device ID: %s", f.DeviceID)
	}

	switch f.Type {
	case feature.FTHeatZone:
		return &cmd.Func{
			Func: func() error {
				devID, err := strconv.Atoi(d.Address)
				if err != nil {
					return fmt.Errorf("honeywell device does not have valid device ID in the address field %s, feature ID: %s",
						d.Address, f.ID)
				}

				thermostat := honeywellExt.NewThermostat(devID)
				ctx := context.TODO()
				ctx, cancel := context.WithTimeout(ctx, time.Second*10)
				defer cancel()

				err = thermostat.Connect(ctx, d.Auth.Login, d.Auth.Password)
				if err != nil {
					return fmt.Errorf("failed to connect to honeywell thermostat: %s", err)
				}

				ctx = context.TODO()
				ctx, cancel = context.WithTimeout(ctx, time.Second*10)
				defer cancel()

				//TODO: Allow caller to specify a duration
				targetTemp := command.Attrs[feature.HeatZoneTargetTempLocalID]
				thermostat.HeatMode(ctx, float32(targetTemp.Value.(int32)), 0)

				return nil
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported command type")
}
