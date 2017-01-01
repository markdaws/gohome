package fluxwifi

import (
	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/gohome"
)

type extension struct {
	gohome.NullExtension
}

func (e *extension) EventsForDevice(sys *gohome.System, d *gohome.Device) *gohome.ExtEvents {
	switch d.ModelNumber {
	case "fluxwifi":
		evts := &gohome.ExtEvents{}
		evts.Producer = &producer{
			Name:   d.Name,
			Device: d,
			System: sys,
		}
		evts.Consumer = &consumer{
			Name:   d.Name,
			System: sys,
			Device: d,
		}
		return evts
	default:
		return nil
	}
}

func (e *extension) BuilderForDevice(sys *gohome.System, d *gohome.Device) cmd.Builder {
	switch d.ModelNumber {
	case "fluxwifi":
		return &cmdBuilder{System: sys}
	default:
		return nil
	}
}

func (e *extension) NetworkForDevice(sys *gohome.System, d *gohome.Device) gohome.Network {
	switch d.ModelNumber {
	case "fluxwifi":
		return &network{}
	default:
		return nil
	}
}

func (e *extension) Discovery(sys *gohome.System) gohome.Discovery {
	return &discovery{System: sys}
}

func (e *extension) Name() string {
	return "fluxwifi"
}

func NewExtension() *extension {
	return &extension{}
}
