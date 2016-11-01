package belkin

import (
	belkinExt "github.com/go-home-iot/belkin"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type extension struct {
	gohome.NullExtension
}

func (e *extension) EventsForDevice(sys *gohome.System, d *gohome.Device) *gohome.ExtEvents {
	var devType belkinExt.DeviceType

	switch d.ModelNumber {
	case "f7c043fc":
		devType = belkinExt.DTMaker
	case "f7c029v2":
		devType = belkinExt.DTInsight
	}

	if devType == "" {
		return nil
	}

	evts := &gohome.ExtEvents{}
	evts.Producer = &producer{
		Name:   d.Name,
		System: sys,
		Device: d,

		// Maker only has one sensor, we just hard code the address to 1 when we create it
		// in the extension scan code
		Sensor: d.Sensors["1"],

		// Maker only has one zone, we set the address to 1 when we did a scan and imported
		// the maker device
		Zone:       d.Zones["1"],
		DeviceType: devType,
	}
	evts.Consumer = &consumer{
		Name:       d.Name,
		System:     sys,
		Device:     d,
		Sensor:     d.Sensors["1"],
		Zone:       d.Zones["1"],
		DeviceType: devType,
	}
	return evts
}

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
