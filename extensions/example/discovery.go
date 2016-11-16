package example

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/zone"
)

type discovery struct{}

// Slice explaining all of the different types of discovers we support
var infos = []gohome.DiscovererInfo{gohome.DiscovererInfo{
	// This must be a globally unique id, good to make it
	// packagename.<id> e.g. example.hardwareA so that it won't
	// clash with other extensions
	//
	// This ID is then passed back to the DiscovererFromID function below
	// if the user chooses to import this piece of hardware
	ID: "example.hardware.1",

	// This string is shown in the import UI, give a brief but unique string
	// that will be shown to the user. It should be short, you will get more
	// opportunity in the Description field to put more information.
	Name: "Example Hardware Version 1.0",

	// This string is shown in the import UI, give more details on what hardware
	// this option supports.
	Description: "Discover version 1.0 example hardware",

	// This is a string you can show to the user before the system will try to scan the network
	// for devices, you can give more info. For example the user might have to press a
	// physcial button on the device before we can scan the network, or if this is a config
	// file string, you can give instructions to the user on how to get this string from the
	// existing app.
	PreScanInfo: "Please press the \"Scan\" button on you hardware hub before trying to scan for devices",
}}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	// Here you return a slice of DiscoveryInfo instances, these will be used
	// to show entries in the import UI

	// If you suport more than one piece of hardware, you return on DiscoveryInfo
	// instance per type of hardware

	return infos
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	// When the user has chosen to import some hardware, this function gets called with the
	// ID of the DiscoveryInfo, if this ID is one we own, e.g. like "example.hardware.1" above
	// then we execute our code that scans the network to find devices, otherwise we just return
	// nil

	switch ID {
	case "example.hardware.1":
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
	// This function will be called when the system wants you to scan the
	// network for your hardware

	// In this example we pretend to scan the network, then return some fake device
	// with a zone attached to it.  To see how you can do this using SSDP or other
	// methods, look at other extensions e.g. gohome/extensions/connectedbytcp/discovery.go

	log.V("scanning for example hardware")

	dev := gohome.NewDevice(
		"",
		"fake hardware name",
		"fake hardware description",
		"example.hardware.1",
		"example model name",
		"example softeare version 1.0",
		"some.fake.IP.address",
		nil,
		nil,
		nil,
		nil,
	)

	// You must give each object an ID, which will be unique in the system before returning
	// the data back to the clients
	dev.ID = sys.NextGlobalID()

	// Add one zone to the device
	z := &zone.Zone{
		ID:          sys.NextGlobalID(),
		Address:     "1",
		Name:        "fake zone 1",
		Description: "fake zone 1 desc",
		DeviceID:    dev.ID,

		// Make this a dimmable light zone, there are many other types
		Type:   zone.ZTLight,
		Output: zone.OTContinuous,
	}
	dev.AddZone(z)

	// Add a fake sensor to the device, each sensor can have one attribute that
	// determines the type of the data returned from the device
	sensor := &gohome.Sensor{
		ID:          sys.NextGlobalID(),
		Address:     "1",
		Name:        "fake sensor",
		Description: "",
		DeviceID:    dev.ID,
		Attr: gohome.SensorAttr{
			Name:     "sensor",
			Value:    "-1",
			DataType: gohome.SDTInt,
			States: map[string]string{
				"0": "Closed",
				"1": "Open",
			},
		},
	}
	dev.AddSensor(sensor)

	// Return these devices back to the UI, the user can then choose to import them or not
	return &gohome.DiscoveryResults{
		Devices: []*gohome.Device{dev},
	}, nil

}
