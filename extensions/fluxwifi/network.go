package fluxwifi

import (
	"fmt"
	"net"

	"github.com/go-home-iot/connection-pool"
	fluxwifiExt "github.com/go-home-iot/fluxwifi"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/zone"
)

type network struct {
	System *gohome.System
}

func (d *network) Devices(sys *gohome.System, modelNumber string) ([]*gohome.Device, error) {
	infos, err := fluxwifiExt.Scan(5)
	if err != nil {
		return nil, err
	}

	devices := make([]*gohome.Device, len(infos))
	for i, info := range infos {
		name := info.ID + ": " + info.Model
		modelNumber := "fluxwifi"
		builderID := modelNumber

		cmdBuilder, ok := sys.Extensions.CmdBuilders[builderID]
		if !ok {
			return nil, fmt.Errorf("unsupported command builder ID: %s", modelNumber)
		}

		dev, _ := gohome.NewDevice(
			modelNumber,
			"",
			"",
			info.IP,
			"",
			name,
			"",
			nil,
			false,
			cmdBuilder,
			&pool.Config{
				Name: name,
				Size: 2,
				//TODO::::!!!!
			},
			nil,
		)

		z := &zone.Zone{
			Address:     "",
			Name:        dev.Name,
			Description: "",
			DeviceID:    "",
			Type:        zone.ZTLight,
			Output:      zone.OTRGB,
		}
		dev.AddZone(z)
		devices[i] = dev
	}
	return devices, nil
}

func (d *network) NewConnection(sys *gohome.System, dev *gohome.Device) (func(pool.Config) (net.Conn, error), error) {
	return func(cfg pool.Config) (net.Conn, error) {
		return net.Dial("tcp", dev.Address)
	}, nil
}
