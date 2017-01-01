package connectedbytcp

import (
	"context"
	"fmt"
	"time"

	connectedbytcpExt "github.com/go-home-iot/connectedbytcp"
	"github.com/markdaws/gohome/pkg/attr"
	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/gohome"
)

type cmdBuilder struct {
	System *gohome.System
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

		for _, attribute := range command.Attrs {
			attribute := attribute

			switch attribute.Type {
			case attr.ATOnOff:
				if attribute.Value.(int32) == attr.OnOffOff {
					return &cmd.Func{
						Func: func() error {
							ctx := context.TODO()
							ctx, cancel := context.WithTimeout(ctx, time.Second*10)
							defer cancel()
							return connectedbytcpExt.TurnOff(ctx, d.Address, f.Address, d.Auth.Token)
						},
						Friendly: "connectedbytcp.cmdBuilder.ZoneTurnOff",
					}, nil
				} else {
					return &cmd.Func{
						Func: func() error {
							ctx := context.TODO()
							ctx, cancel := context.WithTimeout(ctx, time.Second*10)
							defer cancel()
							return connectedbytcpExt.TurnOn(ctx, d.Address, f.Address, d.Auth.Token)
						},
						Friendly: "connectedbytcp.cmdBuilder.ZoneTurnOn",
					}, nil
				}
			case attr.ATBrightness:
				return &cmd.Func{
					Func: func() error {
						ctx := context.TODO()
						ctx, cancel := context.WithTimeout(ctx, time.Second*10)
						defer cancel()
						return connectedbytcpExt.SetLevel(ctx, d.Address, f.Address, d.Auth.Token, int32(attribute.Value.(float32)))
					},
					Friendly: "connectedbytcp.cmdBuilder.ZoneSetLevel",
				}, nil
			}
		}

		return nil, fmt.Errorf("unsupported attribute")
	default:
		return nil, fmt.Errorf("unsupported command type")
	}
}
