package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

type Device interface {
	Address() string
	ID() string
	Name() string
	Description() string
	ModelNumber() string
	Buttons() map[string]*Button
	Devices() map[string]Device
	Zones() map[string]*zone.Zone
	Auth() *comm.Auth
	InitConnections()
	Connect() (comm.Connection, error)
	ReleaseConnection(comm.Connection)
	Stream() bool
	BuildCommand(cmd.Command) (*cmd.Func, error)
	SupportsController(c zone.Controller) bool

	AddZone(z *zone.Zone) error
	comm.Authenticator
	event.Producer
	fmt.Stringer
}

type device struct {
	address     string
	id          string
	name        string
	description string
	system      *System
	//TODO: delete
	producesEvents bool
	auth           *comm.Auth
	buttons        map[string]*Button
	devices        map[string]Device
	zones          map[string]*zone.Zone
	stream         bool
	evpDone        chan bool
	evpFire        chan event.Event
}

func NewDevice(modelNumber, address, ID, name, description string, stream bool, auth *comm.Auth) Device {
	device := device{
		address:     address,
		id:          ID,
		name:        name,
		description: description,
		stream:      stream,
		buttons:     make(map[string]*Button),
		devices:     make(map[string]Device),
		zones:       make(map[string]*zone.Zone),
		auth:        auth,
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

func (d *device) Auth() *comm.Auth {
	return d.auth
}

func (d *device) Buttons() map[string]*Button {
	return d.buttons
}

func (d *device) Devices() map[string]Device {
	return d.devices
}

func (d *device) Zones() map[string]*zone.Zone {
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

func (d *device) AddZone(z *zone.Zone) error {
	errs := &validation.Errors{}

	// Make sure zone doesn't have same address as any other zone
	for _, cz := range d.zones {
		if cz.Address == z.Address {
			errs.Add("device already has a zone with the same address, must be unique", "Address")
			return errs
		}
	}

	d.zones[z.Address] = z
	//TODO
	/*
		// Verify controller is supported by the device
		if !interface{}(d).(Device).SupportsController(zone.ControllerFromString(z.Controller)) {
			errs.Add("the device does not support controlling this kind of zone", "Controller")
			return errs
		}*/

	return nil
}
