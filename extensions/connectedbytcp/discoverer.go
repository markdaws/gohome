package connectedbytcp

import (
	"fmt"

	connectedbytcpExt "github.com/go-home-iot/connectedbytcp"
	"github.com/markdaws/gohome"
)

type discoverer struct {
	System *gohome.System
}

func (d *discoverer) Devices(sys *gohome.System, modelNumber string) ([]*gohome.Device, error) {
	infos, err := connectedbytcpExt.Scan(5)
	if err != nil {
		return nil, err
	}

	devices := make([]*gohome.Device, len(infos))
	for i, info := range infos {
		cmdBuilder, ok := sys.Extensions.CmdBuilders[modelNumber]
		if !ok {
			return nil, fmt.Errorf("unsupported command builder ID: %s", modelNumber)
		}

		//TODO: Need to send back a flag indicating this device needs some kind of
		//authentication to work
		dev, _ := gohome.NewDevice(
			modelNumber,
			info.Location,
			"",
			"ConnectedByTcp - ID: "+info.DeviceID,
			"",
			nil,
			false,
			cmdBuilder,
			nil,
			nil,
		)

		/*
			//TODO: Need to get once we have the security information
						z := &zone.Zone{
							Address:     "",
							Name:        dev.Name,
							Description: "",
							DeviceID:    "",
							Type:        zone.ZTLight,
							Output:      zone.OTRGB,
						}
						dev.AddZone(z)*/
		devices[i] = dev
	}
	return devices, nil
}

/*
//TODO: Move into discoverer type
func DiscoverToken(modelNumber, address string) (string, error) {
	switch modelNumber {
	case "TCP600GWB":
		token, err := connectedbytcp.GetToken(address)
		if err == connectedbytcp.ErrUnauthorized {
			return "", ErrUnauthorized
		}
		return token, err
	}
	return "", ErrUnsupported
}

//TODO:
func VerifyConnection(modelNumber, address, token string) error {
	switch modelNumber {
	case "TCP600GWB":
		err := connectedbytcp.VerifyConnection(address, token)
		if err != nil {
			return fmt.Errorf("access check failed: %s", err)
		}
		return nil
	}
	return ErrUnsupported
}
*/
