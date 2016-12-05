package testing

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/feature"
)

type discovery struct{}

var infos = []gohome.DiscovererInfo{gohome.DiscovererInfo{
	ID:          "testing.hardware",
	Name:        "Testing Hardware",
	Description: "",
	PreScanInfo: "",
}}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	return infos
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "testing.hardware":
		return &discoverer{info: infos[0]}
	default:
		return nil
	}
}

type discoverer struct {
	info gohome.DiscovererInfo
}

func (d *discoverer) Info() gohome.DiscovererInfo {
	return d.info
}

func (d *discoverer) ScanDevices(sys *gohome.System, uiFields map[string]string) (*gohome.DiscoveryResults, error) {
	dev := gohome.NewDevice(
		sys.NewID(),
		"test device",
		"test description",
		"testing.hardware",
		"test model name",
		"test software version 1.0",
		"some.fake.IP.address",
		nil,
		nil,
		nil,
		nil,
	)

	light := feature.NewLightZone(
		sys.NewID(),
		feature.LightZoneModeBinary)
	light.Address = "1"
	light.Name = "onoff light"
	light.DeviceID = dev.ID
	dev.AddFeature(light)

	dimmableLight := feature.NewLightZone(
		sys.NewID(),
		feature.LightZoneModeContinuous)
	dimmableLight.Address = "2"
	dimmableLight.Name = "dimmable light"
	dimmableLight.DeviceID = dev.ID
	dev.AddFeature(dimmableLight)

	openClose := attr.NewOpenClose("openclose", nil)
	sensor := feature.NewSensor(
		sys.NewID(),
		openClose,
	)
	sensor.Address = "3"
	sensor.Name = "test sensor"
	sensor.DeviceID = dev.ID
	dev.AddFeature(sensor)

	swtch := feature.NewSwitch(sys.NewID())
	swtch.Address = "4"
	swtch.Name = "switch"
	swtch.DeviceID = dev.ID
	dev.AddFeature(swtch)

	heat := feature.NewHeatZone(sys.NewID())
	heat.Address = "5"
	heat.Name = "heat"
	heat.DeviceID = dev.ID
	dev.AddFeature(heat)

	window := feature.NewWindowTreatment(sys.NewID())
	window.Address = "6"
	window.Name = "window"
	window.DeviceID = dev.ID
	dev.AddFeature(window)

	colorLight := feature.NewLightZone(
		sys.NewID(),
		feature.LightZoneModeHSL)
	colorLight.Address = "7"
	colorLight.Name = "colour light"
	colorLight.DeviceID = dev.ID
	dev.AddFeature(colorLight)

	button := feature.NewButton(sys.NewID())
	button.Name = "button1"
	button.Address = "8"
	button.DeviceID = dev.ID
	dev.AddFeature(button)

	outlet := feature.NewOutlet(sys.NewID())
	outlet.Name = "outlet"
	outlet.Address = "9"
	outlet.DeviceID = dev.ID
	dev.AddFeature(outlet)

	return &gohome.DiscoveryResults{
		Devices: []*gohome.Device{dev},
	}, nil
}
