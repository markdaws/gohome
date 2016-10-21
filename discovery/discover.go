package discovery

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrUnsupported = errors.New("unsupported model number")

func Devices(sys *gohome.System, modelNumber string) ([]*gohome.Device, error) {

	network, ok := sys.Extensions.Network[modelNumber]
	if !ok {
		return nil, fmt.Errorf("unsupported model number: %s, no registered extension", modelNumber)
	}

	return network.Devices(sys, modelNumber)
}
