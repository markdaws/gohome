package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
)

//TODO: export
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

func (d *genericDevice) BuildCommand(c cmd.Command) (*cmd.Func, error) {
	return nil, fmt.Errorf("genericDevice does not support building commands")
}

func (d *genericDevice) Connect() (comm.Connection, error) {
	return nil, fmt.Errorf("unsupported function connect")
}

func (d *genericDevice) ReleaseConnection(c comm.Connection) {
}
