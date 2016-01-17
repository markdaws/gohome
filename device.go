package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/comm"
)

type Device interface {
	LocalID() string
	GlobalID() string
	Name() string
	Description() string
	System() *System
	Buttons() map[string]*Button
	Devices() map[string]Device
	Zones() map[string]*Zone
	ConnectionInfo() comm.ConnectionInfo
	InitConnections()
	Connect() (comm.Connection, error)
	ReleaseConnection(comm.Connection)
	Authenticate(comm.Connection) error
	Stream() bool

	ZoneSetLevel(z *Zone, level float32) error

	EventProducer
	fmt.Stringer
}

type device struct {
	localID        string
	globalID       string
	name           string
	description    string
	system         *System
	producesEvents bool
	connectionInfo comm.ConnectionInfo
	buttons        map[string]*Button
	devices        map[string]Device
	zones          map[string]*Zone
	stream         bool
	evpDone        chan bool
	evpFire        chan Event
	pool           comm.ConnectionPool
	cmdProcessor   CommandProcessor
}

func NewDevice(modelNumber, localID, globalID, name, description string, producesEvents, stream bool, s *System, cp CommandProcessor, ci comm.ConnectionInfo) Device {
	device := device{
		localID:        localID,
		globalID:       globalID,
		name:           name,
		description:    description,
		producesEvents: producesEvents,
		stream:         stream,
		system:         s,
		buttons:        make(map[string]*Button),
		devices:        make(map[string]Device),
		zones:          make(map[string]*Zone),
		cmdProcessor:   cp,
		connectionInfo: ci,
	}

	switch modelNumber {
	case "":
		return &genericDevice{device: device}
	case "tcphub":
		return &Tcp600gwbDevice{device: device}
	case "L-BDGPRO2-WH":
		return &Lbdgpro2whDevice{device: device}
	default:
		return nil
	}
}

func (d *device) LocalID() string {
	return d.localID
}

func (d *device) GlobalID() string {
	return d.globalID
}

func (d *device) Name() string {
	return d.name
}

func (d *device) Description() string {
	return d.description
}

func (d *device) ConnectionInfo() comm.ConnectionInfo {
	return d.connectionInfo
}

func (d *device) System() *System {
	return d.system
}

func (d *device) Buttons() map[string]*Button {
	return d.buttons
}

func (d *device) Devices() map[string]Device {
	return d.devices
}

func (d *device) Zones() map[string]*Zone {
	return d.zones
}

func (d *device) Stream() bool {
	return d.stream
}

func (d *device) Connect() (comm.Connection, error) {
	c := d.pool.Get()
	if c == nil {
		return nil, fmt.Errorf("%s - connect failed, no connection available", d)
	}
	return c, nil
}

func (d *device) ReleaseConnection(c comm.Connection) {
	d.pool.Release(c)
}

func (d *device) ProducesEvents() bool {
	return d.producesEvents
}

func (d *device) String() string {
	return fmt.Sprintf("Device[%s]", d.Name)
}
