package belkin

import (
	"fmt"
	"time"

	belkinExt "github.com/go-home-iot/belkin"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/cmd"
)

type cmdBuilder struct {
	System *gohome.System
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

			// If there is an OnOff attribute then set the switch to either on or off
			switch attribute.Type {
			case attr.ATOnOff:
				belkinDev := belkinExt.Device{
					Scan: belkinExt.ScanResponse{
						Location: d.Address,
					},
				}
				return &cmd.Func{
					Func: func() error {
						if attribute.Value.(int32) == attr.OnOffOff {
							return belkinDev.TurnOff(time.Second * 5)
						} else {
							return belkinDev.TurnOn(time.Second * 5)
						}
					},
					Friendly: "belkin.cmdBuilder.onoff",
				}, nil
			}
		}
		return nil, fmt.Errorf("unsupported attribute type")

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
}
