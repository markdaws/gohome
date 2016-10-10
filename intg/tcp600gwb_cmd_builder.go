package intg

import (
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/connectedbytcp"
)

type tcp600gwbCmdBuilder struct {
	System *gohome.System
}

func (b *tcp600gwbCmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneTurnOn:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return connectedbytcp.TurnOn(d.Address(), z.Address, d.Auth().Token)
			},
			Friendly: "tcp600gwbCmdBuilder.ZoneTurnOn",
		}, nil

	case *cmd.ZoneTurnOff:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return connectedbytcp.TurnOff(d.Address(), z.Address, d.Auth().Token)
			},
			Friendly: "tcp600gwbCmdBuilder.ZoneTurnOff",
		}, nil

	case *cmd.ZoneSetLevel:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return connectedbytcp.SetLevel(d.Address(), z.Address, d.Auth().Token, int32(command.Level.Value))
			},
			Friendly: "tcp600gwbCmdBuilder.ZoneSetLevel",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
	return nil, nil
}

func (b *tcp600gwbCmdBuilder) Connections(name, address string) comm.ConnectionPool {
	return nil
}

func (b *tcp600gwbCmdBuilder) ID() string {
	return "tcp600gwb"
}
