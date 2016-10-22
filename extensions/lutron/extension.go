package lutron

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type extension struct{}

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
		return &network{System: sys}
	default:
		return nil
	}
}

func (e *extension) ImporterForDevice(sys *gohome.System, d *gohome.Device) gohome.Importer {
	switch d.ModelNumber {
	case "l-bdgpro2-wh":
		return &importer{System: sys}
	default:
		return nil
	}
}

func (e *extension) Name() string {
	return "Lutron"
}

func NewExtension() *extension {
	return &extension{}
}
