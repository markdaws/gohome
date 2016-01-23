package discovery

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome/connectedbytcp"
	"github.com/markdaws/gohome/fluxwifi"
	"github.com/markdaws/gohome/zone"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrUnsupported = errors.New("unsupported model number")

func Discover(modelNumber string) (map[string]string, error) {
	data := make(map[string]string)

	switch modelNumber {
	case "TCP600GWB":
		location, err := connectedbytcp.Discover()
		if err != nil {
			return nil, fmt.Errorf("discover failed: %s", err)
		}
		data["location"] = location
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
	switch modelNumber {
	case "FluxWIFI":
		infos, err := fluxwifi.Scan(5)
		if err != nil {
			return nil, err
		}

		zones := make([]zone.Zone, len(infos))
		for i, info := range infos {
			zones[i] = zone.Zone{
				Address:     info.IP,
				Name:        info.ID,
				Description: "Flux WIFI - " + info.Model,
				Type:        zone.ZTLight,
				Output:      zone.OTContinuous,
				Controller:  zone.ZCFluxWIFI,
			}

			//TODO: remove
			/*
				zones[i*2+1] = zone.Zone{
					Address:     info.IP + "xx",
					Name:        info.ID + " what",
					Description: "Flux WIFI - " + info.Model,
					Type:        zone.ZTShade,
					Output:      zone.OTContinuous,
					Controller:  zone.ZCFluxWIFI,
				}*/
		}
		return zones, nil
	}
	return nil, ErrUnsupported
}

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
