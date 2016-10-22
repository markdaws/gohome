package discovery

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrUnsupported = errors.New("unsupported model number")

func Devices(sys *gohome.System, modelNumber string) ([]*gohome.Device, error) {

	network := sys.Extensions.FindNetwork(sys, &gohome.Device{ModelNumber: modelNumber})
	if network == nil {
		return nil, fmt.Errorf("unsupported model number: %s, no registered extension", modelNumber)
	}

	return network.Devices(sys, modelNumber)
}
