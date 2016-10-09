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

	// TODO: Really needed, shouldn't use for logic
	ModelNumber() string
	Buttons() map[string]*Button
	Devices() map[string]Device
	Zones() map[string]*zone.Zone
	Auth() *comm.Auth

	InitConnections()

	Connections() comm.ConnectionPool
	SetConnections(comm.ConnectionPool)

	//TODO: Remove
	Connect() (comm.Connection, error)
	//TODO: Remove
	ReleaseConnection(comm.Connection)
	Stream() bool
	BuildCommand(cmd.Command) (*cmd.Func, error)
	Hub() Device
	SetHub(Device)
	AddZone(*zone.Zone) error
	AddDevice(Device) error
	Validate() *validation.Errors

	comm.Authenticator
	event.Producer
	fmt.Stringer

	Builder() cmd.Builder
	SetBuilder(cmd.Builder)
}

type device struct {
	address     string
	id          string
	name        string
	description string
	system      *System
	hub         Device
	//TODO: delete?
	producesEvents bool
	auth           *comm.Auth
	buttons        map[string]*Button
	devices        map[string]Device
	zones          map[string]*zone.Zone

	//TODO: Needed? Clean up
	stream  bool
	evpDone chan bool
	evpFire chan event.Event

	builder     cmd.Builder
	connections comm.ConnectionPool
}

func NewDevice(
	modelNumber,
	address,
	ID,
	name,
	description string,
	hub Device,
	stream bool,
	auth *comm.Auth) Device {
	device := device{
		address:     address,
		id:          ID,
		name:        name,
		description: description,
		hub:         hub,
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
	default:
		return nil
	}
}

func (d *device) Builder() cmd.Builder {
	return d.builder
}

func (d *device) SetBuilder(b cmd.Builder) {
	d.builder = b
}

func (d *device) Connections() comm.ConnectionPool {
	return d.connections
}

func (d *device) SetConnections(c comm.ConnectionPool) {
	d.connections = c
}

func (d *device) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if d.name == "" {
		errors.Add("required field", "Name")
	}

	if errors.Has() {
		return errors
	}
	return nil
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

func (d *device) Hub() Device {
	return d.hub
}

func (d *device) SetHub(h Device) {
	d.hub = h
}

func (d *device) AddZone(z *zone.Zone) error {
	errs := &validation.Errors{}

	// Make sure zone doesn't have same address as any other zone
	for _, cz := range d.zones {
		if cz.Address == z.Address {
			errs.Add(fmt.Sprintf("device already has a zone with the same address [%s], must be unique", z.Address), "Address")
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

func (d *device) AddDevice(cd Device) error {
	if _, ok := d.devices[cd.Address()]; ok {
		return fmt.Errorf("device with address: %s already added to parent device", cd.Address())
	}

	d.devices[cd.Address()] = cd
	return nil
}
