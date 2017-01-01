package gohome

import "github.com/markdaws/gohome/pkg/cmd"

type NullExtension struct{}

func (e *NullExtension) EventsForDevice(sys *System, d *Device) *ExtEvents {
	return nil
}

func (e *NullExtension) BuilderForDevice(sys *System, d *Device) cmd.Builder {
	return nil
}

func (e *NullExtension) NetworkForDevice(sys *System, d *Device) Network {
	return nil
}

func (e *NullExtension) Discovery(sys *System) Discovery {
	return nil
}
