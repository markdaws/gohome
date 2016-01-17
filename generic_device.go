package gohome

import "github.com/markdaws/gohome/comm"

type genericDevice struct {
	device
}

func (d *genericDevice) InitConnections() {
}
func (d *genericDevice) StartProducingEvents() (<-chan Event, <-chan bool) {
	return nil, nil
}
func (d *genericDevice) Authenticate(c comm.Connection) error {
	return nil
}
func (d *genericDevice) ZoneSetLevel(z *Zone, level float32) error {
	return nil
}
