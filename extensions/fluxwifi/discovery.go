package fluxwifi

import (
	"errors"

	"github.com/go-home-iot/connection-pool"
	fluxwifiExt "github.com/go-home-iot/fluxwifi"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/zone"
)

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
	return []gohome.DiscovererInfo{gohome.DiscovererInfo{
		ID:          "fluxwifi.bulbs",
		Name:        "FluxWIFI",
		Description: "Discover FluxWIFI bulbs",
		Type:        "ScanDevices",
	}}
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "fluxwifi.bulbs":
		return &discoverer{}
	default:
		return nil
	}
}

type discoverer struct {
}

func (d *discoverer) ScanDevices(sys *gohome.System) (*gohome.DiscoveryResults, error) {
	infos, err := fluxwifiExt.Scan(5)
	if err != nil {
		return nil, err
	}

	devices := make([]*gohome.Device, len(infos))
	for i, info := range infos {
		name := info.ID + ": " + info.Model
		modelNumber := "fluxwifi"

		dev, _ := gohome.NewDevice(
			modelNumber,
			"",
			"",
			info.IP,
			"",
			name,
			"",
			nil,
			nil,
			pool.NewPool(pool.Config{
				Name: name,
				Size: 2,
			}),
			nil,
		)

		z := &zone.Zone{
			Address:     "1",
			Name:        dev.Name,
			Description: "",
			DeviceID:    "",
			Type:        zone.ZTLight,
			Output:      zone.OTRGB,
		}
		dev.AddZone(z)
		devices[i] = dev
	}

	return &gohome.DiscoveryResults{
		Devices: devices,
	}, nil

}
func (d *discoverer) FromString(body string) (*gohome.DiscoveryResults, error) {
	return nil, errors.New("unsupported")
}
