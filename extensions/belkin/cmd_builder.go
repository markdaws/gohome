package belkin

import (
	"fmt"

	"github.com/markdaws/gohome"
	belkinExt "github.com/markdaws/gohome/belkin"
	"github.com/markdaws/gohome/cmd"
)

type cmdBuilder struct {
	System *gohome.System
}

func (b *cmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneTurnOn:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return belkinExt.TurnOn(d.Address)
			},
			Friendly: "belkin.cmdBuilder.ZoneTurnOn",
		}, nil

	case *cmd.ZoneTurnOff:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return belkinExt.TurnOff(d.Address)
			},
			Friendly: "belkin.cmdBuilder.ZoneTurnOff",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
	return nil, nil
}

func (b *cmdBuilder) ID() string {
	return "f7c029v2"
}
