package discovery

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/connectedbytcp"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrUnsupported = errors.New("unsupported model number")

func Devices(sys *gohome.System, modelNumber string) ([]gohome.Device, error) {

	discoverer, ok := sys.Extensions.Discoverers[modelNumber]
	if !ok {
		return nil, fmt.Errorf("unsupported model number: %s, no registered extension", modelNumber)
	}

	return discoverer.Devices(sys, modelNumber)
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
