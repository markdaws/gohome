package gohome

import (
	"fmt"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

type Device struct {
	Address     string
	ID          string
	Name        string
	Description string
	//TODO: Add ModelName, FirmwareVersion
	ModelNumber string
	Buttons     map[string]*Button
	Devices     map[string]*Device
	Zones       map[string]*zone.Zone
	CmdBuilder  cmd.Builder
	Connections *pool.ConnectionPool
	Auth        *comm.Auth
	Hub         *Device

	//TODO: delete?
	producesEvents bool

	//TODO: Needed? Clean up
	Stream  bool
	evpDone chan bool
	evpFire chan event.Event
}

func NewDevice(
	modelNumber,
	address,
	ID,
	name,
	description string,
	hub *Device,
	stream bool,
	cmdBuilder cmd.Builder,
	connPoolCfg *pool.Config,
	auth *comm.Auth) (*Device, error) {

	dev := &Device{
		Address:     address,
		ModelNumber: modelNumber,
		ID:          ID,
		Name:        name,
		Description: description,
		Hub:         hub,
		Buttons:     make(map[string]*Button),
		Devices:     make(map[string]*Device),
		Zones:       make(map[string]*zone.Zone),
		Stream:      stream,
		Auth:        auth,
		CmdBuilder:  cmdBuilder,
	}

	if connPoolCfg != nil {
		dev.SetConnPoolCfg(*connPoolCfg)
	}

	return dev, nil
}

func (d *Device) SetConnPoolCfg(cfg pool.Config) {
	d.Connections = pool.NewPool(cfg)
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
	return nil
}

func (d *Device) AddDevice(cd *Device) error {
	if _, ok := d.Devices[cd.Address]; ok {
		return fmt.Errorf("device with address: %s already added to parent device", cd.Address)
	}

	d.Devices[cd.Address] = cd
	return nil
}
