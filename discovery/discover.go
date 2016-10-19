package discovery

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome"
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
