package honeywell

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type extension struct {
	gohome.NullExtension
}

func (e *extension) Name() string {
	return "honeywell"
}

func (e *extension) BuilderForDevice(sys *gohome.System, d *gohome.Device) cmd.Builder {

	switch d.ModelNumber {
	case "honeywell.redlink.thermostat":
		return &cmdBuilder{Device: d, System: sys}
	default:
		return nil
	}
}

func (e *extension) NetworkForDevice(sys *gohome.System, d *gohome.Device) gohome.Network {
	return nil
}

func (e *extension) EventsForDevice(sys *gohome.System, d *gohome.Device) *gohome.ExtEvents {
	switch d.ModelNumber {
	case "honeywell.redlink.thermostat":
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
