package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/comm"
)

type genericDevice struct {
	device
}

func (d *genericDevice) ModelNumber() string {
	return ""
}

func (d *genericDevice) InitConnections() {
}

func (d *genericDevice) StartProducingEvents() (<-chan Event, <-chan bool) {
	return nil, nil
}

func (d *genericDevice) Authenticate(c comm.Connection) error {
	return nil
}

func (d *genericDevice) BuildCommand(c Command) (*FuncCommand, error) {
	return nil, fmt.Errorf("genericDevice does not support building commands")
}
