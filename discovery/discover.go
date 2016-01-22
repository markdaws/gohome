package discovery

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome/connectedbytcp"
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
