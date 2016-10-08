package intg

import (
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/belkin"
	"github.com/markdaws/gohome/cmd"
)

// Returns a cmd.Builder given the builder ID
func CmdBuilderFromID(system *gohome.System, ID string) (cmd.Builder, error) {
	switch ID {
	case "belkin-wemo-insight":
		return &belkinCmdBuilder{system}, nil
	default:
		return nil, fmt.Errorf("unsupported command builder ID %s", ID)
	}
}

type belkinCmdBuilder struct {
	System *gohome.System
}

func (b *belkinCmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneTurnOn:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return belkin.TurnOn(d.Address())
			},
			Friendly: "belkinCmdBuilder.ZoneTurnOn",
		}, nil

	case *cmd.ZoneTurnOff:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return belkin.TurnOff(d.Address())
			},
			Friendly: "belkinCmdBuilder.ZoneTurnOn",
		}, nil

	default:
		return nil, fmt.Errorf("Unsupported command type")
	}
	return nil, nil
}

func (b *belkinCmdBuilder) ID() string {
	return "belkin-wemo-insight"
}
