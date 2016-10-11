package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

//TODO: Search all refs to gohome.Device or Device see if should be pointers
/*

//TODO: change to struct once lutron device refactored
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

	//TODO: Remove
	Auth() *comm.Auth

	Connections() comm.ConnectionPool
	SetConnections(comm.ConnectionPool)

	//TODO: Remove
	Connect() (comm.Connection, error)
	//TODO: Remove
	ReleaseConnection(comm.Connection)

	Stream() bool

	Hub() Device
	SetHub(Device)
	AddZone(*zone.Zone) error
	AddDevice(Device) error
	Validate() *validation.Errors

	comm.Authenticator
	event.Producer
	fmt.Stringer

	CmdBuilder() cmd.Builder
	SetCmdBuilder(cmd.Builder)
}
*/

type Device struct {
	Address     string
	ID          string
	Name        string
	Description string
	ModelNumber string
	System      *System
	Hub         *Device

	//TODO: delete?
	producesEvents bool
	Auth           *comm.Auth
	Buttons        map[string]*Button
	Devices        map[string]Device
	Zones          map[string]*zone.Zone

	//TODO: Needed? Clean up
	Stream  bool
	evpDone chan bool
	evpFire chan event.Event

	CmdBuilder  cmd.Builder
	Connections comm.ConnectionPool
}

func NewDevice(
	modelNumber,
	address,
	ID,
	name,
	description string,
	hub *Device,
	stream bool,
	auth *comm.Auth) Device {
	device := Device{
		Address:     address,
		ID:          ID,
		Name:        name,
		Description: description,
		Hub:         hub,
		Buttons:     make(map[string]*Button),
		Devices:     make(map[string]Device),
		Zones:       make(map[string]*zone.Zone),
		Stream:      stream,
		Auth:        auth,
	}

	return device
	/*
		switch modelNumber {
		case "":
			device.producesEvents = false
			return &genericDevice{device: device}

				//TODO: Remove
					case "L-BDGPRO2-WH":
						device.producesEvents = true
						return &Lbdgpro2whDevice{device: device}
		default:
			return nil
		}*/
}

//TODO: Delete
func (d *Device) Authenticate(comm.Connection) error {
	return nil
}

//TODO: Delete
func (d *Device) StartProducingEvents() (<-chan event.Event, <-chan bool) {
	return nil, nil
}

func (d *Device) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if d.Name == "" {
		errors.Add("required field", "Name")
	}

	if errors.Has() {
		return errors
	}
	return nil
}

func (d *Device) ProducesEvents() bool {
	return d.producesEvents
}

func (d *Device) String() string {
	return fmt.Sprintf("Device[%s]", d.Name)
}

func (d *Device) AddZone(z *zone.Zone) error {
	errs := &validation.Errors{}

	// Make sure zone doesn't have same address as any other zone
	for _, cz := range d.Zones {
		if cz.Address == z.Address {
			errs.Add(fmt.Sprintf("device already has a zone with the same address [%s], must be unique", z.Address), "Address")
			return errs
		}
	}

	d.Zones[z.Address] = z
	//TODO
	/*
		// Verify controller is supported by the device
		if !interface{}(d).(Device).SupportsController(zone.ControllerFromString(z.Controller)) {
			errs.Add("the device does not support controlling this kind of zone", "Controller")
			return errs
		}*/

	return nil
}

func (d *Device) AddDevice(cd Device) error {
	if _, ok := d.Devices[cd.Address]; ok {
		return fmt.Errorf("device with address: %s already added to parent device", cd.Address)
	}

	d.Devices[cd.Address] = cd
	return nil
}
