package belkin

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type extension struct{}

func (e *extension) BuilderForDevice(sys *gohome.System, d *gohome.Device) cmd.Builder {
	// Given the device we can return different builds for different devices and even
	// take in to account SoftwareVersion as a field to return a different builder
	switch d.ModelNumber {
	case "f7c043fc":
		//WeMo Maker
		return &cmdBuilder{System: sys}
	case "f7c029v2":
		//WeMo Insight
		return &cmdBuilder{System: sys}
	default:
		return nil
	}
}

func (e *extension) NetworkForDevice(sys *gohome.System, d *gohome.Device) gohome.Network {
	switch d.ModelNumber {
	case "f7c043fc":
		//WeMo Maker
		return &network{System: sys}
	case "f7c029v2":
		//WeMo Insight
		return &network{System: sys}
	default:
		return nil
	}
}

func (e *extension) ImporterForDevice(sys *gohome.System, d *gohome.Device) gohome.Importer {
	return nil
}

func (e *extension) Name() string {
	return "Belkin"
}

func NewExtension() *extension {
	return &extension{}
}
