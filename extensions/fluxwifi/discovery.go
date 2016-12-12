package fluxwifi

import (
	"github.com/go-home-iot/connection-pool"
	fluxwifiExt "github.com/go-home-iot/fluxwifi"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/feature"
)

var infos = []gohome.DiscovererInfo{gohome.DiscovererInfo{
	ID:          "fluxwifi.bulbs",
	Name:        "FluxWIFI Bulbs",
	Description: "Discover FluxWIFI bulbs",
}}

type discovery struct {
	//TODO: Needed?
	System *gohome.System
}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	// Here we return all of the discoverers we support. There may be multiple types
	// of hardware we can import, for each type, or even one piece of hardware with
	// different software versions, which may need different plumbing, we return a
	// discovery info, the Name/Description fields will be displayed to the user, then
	// the ID is then passed back to the DiscovererFromID function, where we can then
	// return the appropriate discoverer instance
	return infos
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "fluxwifi.bulbs":
		return &discoverer{infos[0]}
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
	infos, err := fluxwifiExt.Scan(5)
	if err != nil {
		return nil, err
	}

	devices := make([]*gohome.Device, len(infos))
	for i, info := range infos {
		name := info.ID + ": " + info.Model
		modelNumber := "fluxwifi"

		dev := gohome.NewDevice(
			sys.NewID(),
			name,
			"",
			modelNumber,
			"",
			"",
			info.IP,
			nil,
			nil,
			pool.NewPool(pool.Config{
				Name: name,
				Size: 2,
			}),
			nil,
		)

		light := feature.NewLightZone(sys.NewID(), feature.LightZoneModeHSL)
		light.Name = dev.Name
		light.Address = "1"
		light.DeviceID = dev.ID
		dev.AddFeature(light)
		devices[i] = dev
	}

	return &gohome.DiscoveryResults{
		Devices: devices,
	}, nil

}
