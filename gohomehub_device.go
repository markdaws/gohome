package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/comm"
)

type GoHomeHubDevice struct {
	device
}

func (d *GoHomeHubDevice) ModelNumber() string {
	return "GoHomeHub"
}

func (d *GoHomeHubDevice) InitConnections() {
}

func (d *GoHomeHubDevice) StartProducingEvents() (<-chan Event, <-chan bool) {
	return nil, nil
}

func (d *GoHomeHubDevice) Authenticate(c comm.Connection) error {
	return nil
}

func (d *GoHomeHubDevice) BuildCommand(c Command) (*FuncCommand, error) {
	switch cmd := c.(type) {
	case *ZoneSetLevelCommand:
		//TODO: Get the zone, get the type of bulb
		//Given the bulb type, then we can figure out how to communicate with it
		return nil, fmt.Errorf("goHomeHubDevice ZoneSetLevelCommand not supported")
	case *ButtonPressCommand:
		//TODO: Phantom buttons?
		return nil, fmt.Errorf("goHomeHubDevice ButtonPressCommand not supported")
	case *ButtonReleaseCommand:
		return nil, fmt.Errorf("goHomeHubDevice ButtonReleaseCommand not supported")
	case *SceneSetCommand:
		//TODO: Does this make sense, what does a scene mean in terms of this virtual hub?
	default:
		_ = cmd
		return nil, fmt.Errorf("goHomeHubDevice build commands not supported")
	}

	return nil, fmt.Errorf("goHomeHubDevice unsupported command")
}

/*
- For wifi bulbs can control directly
- Capabilities supports zigbee etc
- Need to add bulb information
- Is a zone a bulb?

zone is a controllable unit
*/
