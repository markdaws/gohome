package connectedbytcp

import (
	"errors"
	"fmt"
	"net"

	connectedbytcpExt "github.com/go-home-iot/connectedbytcp"
	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
)

type network struct {
	System *gohome.System
}

func (d *network) Devices(sys *gohome.System, modelNumber string) ([]*gohome.Device, error) {
	infos, err := connectedbytcpExt.Scan(5)
	if err != nil {
		return nil, err
	}

	devices := make([]*gohome.Device, len(infos))
	for i, info := range infos {
		//TODO: Need to send back a flag indicating this device needs some kind of
		//authentication to work
		dev, _ := gohome.NewDevice(
			modelNumber,
			"",
			"",
			info.Location,
			"",
			"ConnectedByTcp - ID: "+info.DeviceID,
			"",
			nil,
			nil,
			nil,
			nil,
		)

		cmdBuilder := sys.Extensions.FindCmdBuilder(sys, dev)
		if cmdBuilder == nil {
			return nil, fmt.Errorf("unsupported command builder ID: %s", modelNumber)
		}
		dev.CmdBuilder = cmdBuilder

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

func (d *network) NewConnection(sys *gohome.System, dev *gohome.Device) (func(pool.Config) (net.Conn, error), error) {
	return nil, errors.New("unsupported method")
}

/*
//TODO: Move into network type
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
