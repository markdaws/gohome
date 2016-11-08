package connectedbytcp

import (
	"errors"

	connectedbytcpExt "github.com/go-home-iot/connectedbytcp"
	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/zone"
)

type discovery struct {
	//TODO: Needed?
	System *gohome.System
}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	return []gohome.DiscovererInfo{gohome.DiscovererInfo{
		ID:          "connectedbytcp.bulbs",
		Name:        "ConnectedByTCP Bulbs",
		Description: "Discover ConnectedByTCP bulbs",
		Type:        "ScanDevices",
		PreScanInfo: "IMPORTANT: You must press the \"Scan\" button on your physical hub hardware before trying to discover devices.  If you don't the scan will fail.",
	}}
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "connectedbytcp.bulbs":
		return &discoverer{}
	default:
		return nil
	}
}

type discoverer struct{}

func (d *discoverer) ScanDevices(sys *gohome.System) (*gohome.DiscoveryResults, error) {
	infos, err := connectedbytcpExt.Scan(5)
	if err != nil {
		return nil, err
	}

	devices := make([]*gohome.Device, len(infos))
	for i, info := range infos {
		token, err := connectedbytcpExt.GetToken(info.Location)
		if err != nil {
			return nil, err
		}

		dev := gohome.NewDevice(
			"tcp600gwb",
			"tcp600gwb",
			"",
			info.Location,
			"",
			"ConnectedByTcp - ID: "+info.DeviceID,
			"",
			nil,
			nil,
			pool.NewPool(pool.Config{
				Name: "connectedbytcp",
				Size: 2,
			}),
			&gohome.Auth{
				Token: token,
			},
		)

		// Get all of the rooms and all the zones in each room
		resp, err := connectedbytcpExt.RoomGetCarousel(info.Location, token)
		if err != nil {
			return nil, err
		}

		for _, room := range resp.Rooms {
			for _, d := range room.Devices {
				z := &zone.Zone{
					Address:     d.DID,
					Name:        d.Name,
					Description: "",
					DeviceID:    "",
					Type:        zone.ZTLight,
					Output:      zone.OTContinuous,
				}
				dev.AddZone(z)
			}
		}
		devices[i] = dev
	}

	return &gohome.DiscoveryResults{
		Devices: devices,
	}, nil
}

func (d *discoverer) FromString(body string) (*gohome.DiscoveryResults, error) {
	return nil, errors.New("unsupported")
}
