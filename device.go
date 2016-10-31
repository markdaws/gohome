package gohome

import (
	"fmt"

	"github.com/go-home-iot/connection-pool"
	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

// DeviceType explains the type of a device e.g. Dimmer or Shade
type DeviceType string

const (
	DTDimmer  DeviceType = "dimmer"
	DTSwitch             = "switch"
	DTShade              = "shade"
	DTHub                = "hub"
	DTRemote             = "remote"
	DTUnknown            = "unknown"
)

// Auth contains authentication information such as login/password/security token
type Auth struct {
	Login    string
	Password string
	Token    string
}

type Device struct {
	Address         string
	ID              string
	Name            string
	Description     string
	Type            DeviceType
	ModelName       string
	ModelNumber     string
	SoftwareVersion string
	Buttons         map[string]*Button
	Devices         map[string]*Device
	Zones           map[string]*zone.Zone
	Sensors         map[string]*Sensor
	CmdBuilder      cmd.Builder
	Connections     *pool.ConnectionPool
	Auth            *Auth
	Hub             *Device
}

func NewDevice(
	modelNumber,
	modelName,
	softwareVersion,
	address,
	ID,
	name,
	description string,
	hub *Device,
	cmdBuilder cmd.Builder,
	connPoolCfg *pool.Config,
	auth *Auth) (*Device, error) {

	dev := &Device{
		Address:         address,
		ModelNumber:     modelNumber,
		ModelName:       modelName,
		SoftwareVersion: softwareVersion,
		ID:              ID,
		Name:            name,
		Description:     description,
		Hub:             hub,
		Buttons:         make(map[string]*Button),
		Devices:         make(map[string]*Device),
		Zones:           make(map[string]*zone.Zone),
		Sensors:         make(map[string]*Sensor),
		Auth:            auth,
		CmdBuilder:      cmdBuilder,
	}

	if connPoolCfg != nil {
		dev.SetConnPoolCfg(*connPoolCfg)
	}

	return dev, nil
}

func (d *Device) SetConnPoolCfg(cfg pool.Config) {
	d.Connections = pool.NewPool(cfg)
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

func (d *Device) AddSensor(s *Sensor) error {
	if _, ok := d.Sensors[s.Address]; ok {
		return fmt.Errorf("sensor with address: %s already added to device", s.Address)
	}
	d.Sensors[s.Address] = s
	return nil
}

// ==== evtbus.Producer interface

func (d *Device) ProducerName() string {
	return d.String()
}

func (d *Device) StartProducing(b *evtbus.Bus) {
	//TODO: Raise events like zone added? sensor added, lost connection etc
}

func (d *Device) StopProducing() {
}

// ====================
