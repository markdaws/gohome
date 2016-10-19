// Package belkin provides support for Belkin devices, such as the WeMo Switch
package belkin

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-home-iot/gossdp"
)

// ErrUnsupportedAction is returned when you try to perform an action on a piece of hardware
// that doesn't support it, e.g. calling FetchAttributes on a non Maker device
var ErrUnsupportedAction = errors.New("unsupported action")

// CREDIT: All the knowledge of how to control this product came from:
// https://github.com/timonreinhard/wemo-client

// Scan detects Belkin devices on the network. The devices that are returned have
// limited information in the Scan field, to get more detailed information you will
// have to call Load() on the device
func Scan(dt DeviceType, waitTimeSeconds int) ([]*Device, error) {
	var responses []ScanResponse
	l := belkinListener{
		URN:       string(dt),
		Responses: &responses,
	}

	c, err := gossdp.NewSsdpClientWithLogger(l, l)
	if err != nil {
		return nil, fmt.Errorf("failed to start ssdp discovery client: %s", err)
	}

	defer c.Stop()
	go c.Start()
	err = c.ListenFor(string(dt))
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %s", err)
	}

	time.Sleep(time.Duration(waitTimeSeconds) * time.Second)

	devices := make([]*Device, len(responses))
	for i, response := range responses {
		devices[i] = &Device{Scan: response}
	}
	return devices, nil
}
