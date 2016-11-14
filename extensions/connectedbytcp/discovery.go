package connectedbytcp

import (
	connectedbytcpExt "github.com/go-home-iot/connectedbytcp"
	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/zone"
)

var infos = []gohome.DiscovererInfo{gohome.DiscovererInfo{
	ID:          "connectedbytcp.bulbs",
	Name:        "ConnectedByTCP Bulbs",
	Description: "Discover ConnectedByTCP bulbs",
	PreScanInfo: "IMPORTANT: You must press the \"Scan\" button on your physical hub hardware before trying to discover devices.  If you don't the scan will fail.",
}}

type discovery struct {
	//TODO: Needed?
	System *gohome.System
}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	return infos
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "connectedbytcp.bulbs":
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
		dev.ID = sys.NextGlobalID()

		// Get all of the rooms and all the zones in each room
		resp, err := connectedbytcpExt.RoomGetCarousel(info.Location, token)
		if err != nil {
			return nil, err
		}

		for _, room := range resp.Rooms {
			for _, d := range room.Devices {
				z := &zone.Zone{
					ID:          sys.NextGlobalID(),
					Address:     d.DID,
					Name:        d.Name,
					Description: "",
					DeviceID:    dev.ID,
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
