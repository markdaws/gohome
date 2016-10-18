package belkin

import (
	"fmt"

	belkinExt "github.com/go-home-iot/belkin"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type cmdBuilder struct {
	System *gohome.System
	id     string
}

func (b *cmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneTurnOn:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]

		belkinDev := belkinExt.Device{Scan: belkinExt.ScanResponse{Location: d.Address}}
		return &cmd.Func{
			Func: func() error {
				return belkinDev.TurnOn()
			},
			Friendly: "belkin.cmdBuilder.ZoneTurnOn",
		}, nil

	case *cmd.ZoneTurnOff:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]

		belkinDev := belkinExt.Device{Scan: belkinExt.ScanResponse{Location: d.Address}}
		return &cmd.Func{
			Func: func() error {
				return belkinDev.TurnOff()
			},
			Friendly: "belkin.cmdBuilder.ZoneTurnOff",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
	return nil, nil
}

func (b *cmdBuilder) ID() string {
	return b.id
}
