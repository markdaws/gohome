package connectedbytcp

import (
	"context"
	"time"

	connectedbytcpExt "github.com/go-home-iot/connectedbytcp"
	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/feature"
	errExt "github.com/pkg/errors"
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
		ctx := context.TODO()
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		token, err := connectedbytcpExt.GetToken(ctx, info.Location)
		if err != nil {
			return nil, errExt.Wrap(err, "call to GetToken failed, make sure you pressed the 'Scan' button on the physical hub")
		}

		dev := gohome.NewDevice(
			sys.NewGlobalID(),
			"ConnectedByTcp - ID: "+info.DeviceID,
			"",
			"tcp600gwb",
			"tcp600gwb",
			"",
			info.Location,
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
		ctx = context.TODO()
		ctx, cancel = context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		resp, err := connectedbytcpExt.RoomGetCarousel(ctx, info.Location, token)
		if err != nil {
			return nil, err
		}

		for _, room := range resp.Rooms {
			for _, roomDev := range room.Devices {
				light := feature.NewLightZone(sys.NewGlobalID(), feature.LightZoneModeContinuous)
				light.Name = roomDev.Name
				light.Address = roomDev.DID
				light.DeviceID = dev.ID
				dev.AddFeature(light)
			}
		}
		devices[i] = dev
	}

	return &gohome.DiscoveryResults{
		Devices: devices,
	}, nil
}
