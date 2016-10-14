package discovery

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/connectedbytcp"
	"github.com/markdaws/gohome/fluxwifi"
	"github.com/markdaws/gohome/zone"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrUnsupported = errors.New("unsupported model number")

func Devices(sys *gohome.System, modelNumber string) ([]gohome.Device, error) {

	discoverer, ok := sys.Extensions.Discoverers[modelNumber]
	if !ok {
		return nil, fmt.Errorf("unsupported model number: %s, no registered extension", modelNumber)
	}

	return discoverer.Devices(sys, modelNumber)

	switch modelNumber {
	case "fluxwifi":
		infos, err := fluxwifi.Scan(5)
		if err != nil {
			return nil, err
		}

		/*
			//TODO: Remove, for testing
			infos := [1]fluxwifi.BulbInfo{fluxwifi.BulbInfo{
				IP:    "192.168.0.1:1234",
				ID:    "thisisanid",
				Model: "modelnumber",
			}}*/

		devices := make([]gohome.Device, len(infos))
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
				info.IP,
				"",
				name,
				"",
				nil,
				false,
				cmdBuilder,
				&comm.ConnectionPoolConfig{
					Name:           name,
					Size:           2,
					ConnectionType: "telnet",
					Address:        info.IP,
					TelnetPingCmd:  "",
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
			devices[i] = *dev
		}
		return devices, nil

	default:
		return nil, ErrUnsupported
	}
}

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
