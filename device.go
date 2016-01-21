package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
)

type Device interface {
	Address() string
	ID() string
	Name() string
	Description() string
	ModelNumber() string
	Buttons() map[string]*Button
	Devices() map[string]Device
	Zones() map[string]*Zone
	ConnectionInfo() comm.ConnectionInfo
	InitConnections()
	Connect() (comm.Connection, error)
	ReleaseConnection(comm.Connection)
	Authenticate(comm.Connection) error
	Stream() bool
	BuildCommand(cmd.Command) (*cmd.Func, error)

	EventProducer
	fmt.Stringer
}

type device struct {
	address        string
	id             string
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
}

func NewDevice(modelNumber, address, ID, name, description string, stream bool, ci comm.ConnectionInfo) Device {
	device := device{
		address:        address,
		id:             ID,
		name:           name,
		description:    description,
		stream:         stream,
		buttons:        make(map[string]*Button),
		devices:        make(map[string]Device),
		zones:          make(map[string]*Zone),
		connectionInfo: ci,
	}

	switch modelNumber {
	case "":
		device.producesEvents = false
		return &genericDevice{device: device}
	case "TCP600GWB":
		device.producesEvents = false
		return &Tcp600gwbDevice{device: device}
	case "L-BDGPRO2-WH":
		device.producesEvents = true
		return &Lbdgpro2whDevice{device: device}
	case "GoHomeHub":
		device.producesEvents = true
		return &GoHomeHubDevice{device: device}
	default:
		return nil
	}
}

func (d *device) Address() string {
	return d.address
}

func (d *device) ID() string {
	return d.id
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

func (d *device) ProducesEvents() bool {
	return d.producesEvents
}

func (d *device) String() string {
	return fmt.Sprintf("Device[%s]", d.Name())
}
