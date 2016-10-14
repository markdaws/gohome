package connectedbytcp

import (
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	connectedbytcpExt "github.com/markdaws/gohome/connectedbytcp"
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
				return connectedbytcpExt.TurnOn(d.Address, z.Address, d.Auth.Token)
			},
			Friendly: "connectedbytcp.cmdBuilder.ZoneTurnOn",
		}, nil

	case *cmd.ZoneTurnOff:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return connectedbytcpExt.TurnOff(d.Address, z.Address, d.Auth.Token)
			},
			Friendly: "connectedbytcp.cmdBuilder.ZoneTurnOff",
		}, nil

	case *cmd.ZoneSetLevel:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				return connectedbytcpExt.SetLevel(d.Address, z.Address, d.Auth.Token, int32(command.Level.Value))
			},
			Friendly: "connectedbytcp.cmdBuilder.ZoneSetLevel",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
	return nil, nil
}

func (b *cmdBuilder) ID() string {
	return "tcp600gwb"
}
