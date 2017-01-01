package honeywell

import (
	"github.com/markdaws/gohome/pkg/feature"
	"github.com/markdaws/gohome/pkg/gohome"
	"github.com/markdaws/gohome/pkg/log"
)

type discovery struct{}

var infos = []gohome.DiscovererInfo{gohome.DiscovererInfo{
	ID:          "honeywell.redlink.thermostat",
	Name:        "Honeywell RedLINK Thermostat",
	Description: "Monitor and control Honeywell RedLINK connected thermostats",
	PreScanInfo: "The login/password are the credentials you use on the mytotalconnectcomfort.com website. The Device ID has to be determined manually, log in to the mytotalconnectcomfort website, navigate to your device, then the URL will look something like /portal/Device/CheckDataSession/123456, you need to copy the number that is in place of the 123456 and use that as your device ID.",
	UIFields: []gohome.UIField{
		gohome.UIField{
			ID:          "login",
			Label:       "Login",
			Description: "Login to mytotalconnectcomfort.com website",
			Required:    true,
		},
		gohome.UIField{
			ID:          "password",
			Label:       "Password",
			Description: "Password for mytotalconnectcomfort.com website",
			Required:    true,
		},
		gohome.UIField{
			ID:          "deviceID",
			Label:       "Device ID",
			Description: "The Device ID for the thermostat you wish to control and monitor",
			Required:    true,
		},
	},
}}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	return infos
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "honeywell.redlink.thermostat":
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
	log.V("scanning for honeywell hardware")

	auth := &gohome.Auth{
		Login:    uiFields["login"],
		Password: uiFields["password"],
	}

	dev := gohome.NewDevice(
		sys.NewID(),
		"Honeywell Thermostat",
		"",
		"honeywell.redlink.thermostat",
		"",
		"",
		uiFields["deviceID"],
		nil,
		nil,
		nil,
		auth,
	)

	heat := feature.NewHeatZone(sys.NewID())
	heat.Address = "1"
	heat.Name = "Heat Zone"
	heat.DeviceID = dev.ID
	dev.AddFeature(heat)

	return &gohome.DiscoveryResults{
		Devices: []*gohome.Device{dev},
	}, nil
}
