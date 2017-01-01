package lutron

import (
	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/gohome"
)

type extension struct {
	gohome.NullExtension
}

func (e *extension) EventsForDevice(sys *gohome.System, d *gohome.Device) *gohome.ExtEvents {
	switch d.ModelNumber {
	case "l-bdgpro2-wh":
		evts := &gohome.ExtEvents{}

		// event implements both interfaces
		consumerProducer := &event{
			Name:   d.Name,
			System: sys,
			Device: d,
		}
		evts.Producer = consumerProducer
		evts.Consumer = consumerProducer

		return evts
	default:
		return nil
	}
}

func (e *extension) BuilderForDevice(sys *gohome.System, d *gohome.Device) cmd.Builder {
	switch d.ModelNumber {
	case "l-bdgpro2-wh":
		return &cmdBuilder{System: sys}
	default:
		return nil
	}
}

func (e *extension) NetworkForDevice(sys *gohome.System, d *gohome.Device) gohome.Network {
	switch d.ModelNumber {
	case "l-bdgpro2-wh":
		return &network{}
	default:
		return nil
	}
}

func (e *extension) Discovery(sys *gohome.System) gohome.Discovery {
	return &discovery{System: sys}
}

func (e *extension) Name() string {
	return "Lutron"
}

func NewExtension() *extension {
	return &extension{}
}
