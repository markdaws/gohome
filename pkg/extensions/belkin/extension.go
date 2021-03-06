package belkin

import (
	belkinExt "github.com/go-home-iot/belkin"
	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/gohome"
)

type extension struct {
	gohome.NullExtension
}

func (e *extension) EventsForDevice(sys *gohome.System, d *gohome.Device) *gohome.ExtEvents {
	var devType belkinExt.DeviceType

	switch d.ModelName {
	case "Maker":
		devType = belkinExt.DTMaker
	case "Insight":
		devType = belkinExt.DTInsight
	}

	if devType == "" {
		return nil
	}

	evts := &gohome.ExtEvents{}
	evts.Producer = &producer{
		Name:       d.Name,
		System:     sys,
		Device:     d,
		DeviceType: devType,
	}
	evts.Consumer = &consumer{
		Name:       d.Name,
		System:     sys,
		Device:     d,
		DeviceType: devType,
	}
	return evts
}

func (e *extension) BuilderForDevice(sys *gohome.System, d *gohome.Device) cmd.Builder {
	// Given the device we can return different builds for different devices and even
	// take in to account SoftwareVersion as a field to return a different builder
	switch d.ModelName {
	case "Maker":
		//WeMo Maker
		return &cmdBuilder{System: sys}
	case "Insight":
		//WeMo Insight
		return &cmdBuilder{System: sys}
	default:
		return nil
	}
}

func (e *extension) Discovery(sys *gohome.System) gohome.Discovery {
	return &discovery{System: sys}
}

func (e *extension) Name() string {
	return "Belkin"
}

func NewExtension() *extension {
	return &extension{}
}
