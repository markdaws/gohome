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
		sys.NewGlobalID(),
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
		sys.NewGlobalID(),
		false,
		false)
	light.Address = "1"
	light.Name = "onoff light"
	light.DeviceID = dev.ID
	dev.AddFeature(light)

	dimmableLight := feature.NewLightZone(
		sys.NewGlobalID(),
		true,
		false)
	dimmableLight.Address = "2"
	dimmableLight.Name = "dimmable light"
	dimmableLight.DeviceID = dev.ID
	dev.AddFeature(dimmableLight)

	openClose := attr.NewOpenClose("openclose", nil)
	sensor := feature.NewSensor(
		sys.NewGlobalID(),
		openClose,
	)
	sensor.Address = "3"
	sensor.Name = "test sensor"
	sensor.DeviceID = dev.ID
	dev.AddFeature(sensor)

	swtch := feature.NewSwitch(sys.NewGlobalID())
	swtch.Address = "4"
	swtch.Name = "switch"
	swtch.DeviceID = dev.ID
	dev.AddFeature(swtch)

	heat := feature.NewHeatZone(sys.NewGlobalID())
	heat.Address = "5"
	heat.Name = "heat"
	heat.DeviceID = dev.ID
	dev.AddFeature(heat)

	window := feature.NewWindowTreatment(sys.NewGlobalID())
	window.Address = "6"
	window.Name = "window"
	window.DeviceID = dev.ID
	dev.AddFeature(window)

	hueLight := feature.NewLightZone(
		sys.NewGlobalID(),
		true,
		true)
	hueLight.Address = "7"
	hueLight.Name = "hue light"
	hueLight.DeviceID = dev.ID
	dev.AddFeature(hueLight)

	button := feature.NewButton(sys.NewGlobalID())
	button.Name = "button1"
	button.Address = "8"
	button.DeviceID = dev.ID
	dev.AddFeature(button)

	outlet := feature.NewOutlet(sys.NewGlobalID())
	outlet.Name = "outlet"
	outlet.Address = "9"
	outlet.DeviceID = dev.ID
	dev.AddFeature(outlet)

	return &gohome.DiscoveryResults{
		Devices: []*gohome.Device{dev},
	}, nil
}
