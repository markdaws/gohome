package connectedbytcp

import (
	"context"
	"fmt"
	"time"

	connectedbytcpExt "github.com/go-home-iot/connectedbytcp"
	"github.com/markdaws/gohome"
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
				ctx := context.TODO()
				ctx, cancel := context.WithTimeout(ctx, time.Second*10)
				defer cancel()
				return connectedbytcpExt.TurnOn(ctx, d.Address, z.Address, d.Auth.Token)
			},
			Friendly: "connectedbytcp.cmdBuilder.ZoneTurnOn",
		}, nil

	case *cmd.ZoneTurnOff:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				ctx := context.TODO()
				ctx, cancel := context.WithTimeout(ctx, time.Second*10)
				defer cancel()
				return connectedbytcpExt.TurnOff(ctx, d.Address, z.Address, d.Auth.Token)
			},
			Friendly: "connectedbytcp.cmdBuilder.ZoneTurnOff",
		}, nil

	case *cmd.ZoneSetLevel:
		z := b.System.Zones[command.ZoneID]
		d := b.System.Devices[z.DeviceID]
		return &cmd.Func{
			Func: func() error {
				ctx := context.TODO()
				ctx, cancel := context.WithTimeout(ctx, time.Second*10)
				defer cancel()
				return connectedbytcpExt.SetLevel(ctx, d.Address, z.Address, d.Auth.Token, int32(command.Level.Value))
			},
			Friendly: "connectedbytcp.cmdBuilder.ZoneSetLevel",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
}
