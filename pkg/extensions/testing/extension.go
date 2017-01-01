package testing

import (
	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/gohome"
)

type extension struct {
	gohome.NullExtension
}

func (e *extension) Name() string {
	return "testing"
}

func (e *extension) BuilderForDevice(sys *gohome.System, d *gohome.Device) cmd.Builder {
	switch d.ModelNumber {
	case "testing.hardware":
		return &cmdBuilder{ModelNumber: d.ModelNumber, Device: d}
	default:
		// This device is not one that we know how to control, return nil
		return nil
	}
}

func (e *extension) NetworkForDevice(sys *gohome.System, d *gohome.Device) gohome.Network {
	return nil
}

func (e *extension) EventsForDevice(sys *gohome.System, d *gohome.Device) *gohome.ExtEvents {
	switch d.ModelNumber {
	case "testing.hardware":
		evts := &gohome.ExtEvents{}
		evts.Producer = &producer{
			Device: d,
			System: sys,
		}
		evts.Consumer = &consumer{
			Device: d,
			System: sys,
		}
		return evts
	default:
		return nil
	}
}

func (e *extension) Discovery(sys *gohome.System) gohome.Discovery {
	return &discovery{}
}

func NewExtension() *extension {
	return &extension{}
}
