package discovery

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/belkin"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/connectedbytcp"
	"github.com/markdaws/gohome/fluxwifi"
	"github.com/markdaws/gohome/intg"
	"github.com/markdaws/gohome/zone"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrUnsupported = errors.New("unsupported model number")

func Devices(modelNumber string) ([]gohome.Device, error) {
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
			dev := gohome.NewDevice(
				"fluxwifi",
				info.IP,
				"",
				info.ID+": "+info.Model,
				"",
				nil,
				false,
				nil,
			)

			builder, _ := intg.CmdBuilderFromID(nil, dev.ModelNumber)
			dev.CmdBuilder = builder

			pool, _ := comm.NewConnectionPool(comm.ConnectionPoolConfig{
				Name:           dev.Name,
				Size:           2,
				ConnectionType: "telnet",
				Address:        dev.Address,
				TelnetPingCmd:  "",
			})
			dev.Connections = pool

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

	default:
		return nil, ErrUnsupported
	}
}

//TODO: Delete
func Discover(modelNumber string) (map[string]string, error) {
	data := make(map[string]string)

	switch modelNumber {
	case "TCP600GWB":
		responses, err := connectedbytcp.Scan(5)
		if err != nil {
			return nil, fmt.Errorf("discover failed: %s", err)
		}

		if len(responses) > 0 {
			data["location"] = responses[0].Location
		}
		return data, nil
	}
	return nil, ErrUnsupported
}

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

func Zones(modelNumber string) ([]zone.Zone, error) {
	//TODO: Need to also have a device here along with the zone
	//Shouldn't discover zones, only devices
	switch modelNumber {
	case "FluxWIFI":
		infos, err := fluxwifi.Scan(5)
		if err != nil {
			return nil, err
		}

		zones := make([]zone.Zone, len(infos)*2)
		for i, info := range infos {
			zones[i*2] = zone.Zone{
				Address:     info.IP,
				Name:        info.ID,
				Description: "Flux WIFI - " + info.Model,
				Type:        zone.ZTLight,
				Output:      zone.OTContinuous,
			}
		}
		return zones, nil

	case "F7C029V2":
		responses, err := belkin.Scan(belkin.DTInsight, 5)
		_ = responses
		_ = err
		return nil, fmt.Errorf("//TODO:not implemented")
	}
	return nil, ErrUnsupported
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
