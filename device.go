package gohome

import (
	"fmt"

	"github.com/go-home-iot/connection-pool"
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

// Device is a piece of hardware. It could be a dimmer, a shade, a remote etc
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

// NewDevice returns an initialized device object
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

// Validate checks that all of the requirements for this to be a valid device are met
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

// String returns a friendly string describing the device that can be useful for debugging
func (d *Device) String() string {
	return fmt.Sprintf("Device[%s]", d.Name)
}

// AddZone adds the zone to the device
func (d *Device) AddZone(z *zone.Zone) error {
	errs := &validation.Errors{}

	// Make sure zone doesn't have same address as any other zone
	if _, ok := d.Zones[z.Address]; ok {
		errs.Add(fmt.Sprintf("device already has a zone with the same address [%s], must be unique", z.Address), "Address")
		return errs
	}

	d.Zones[z.Address] = z
	return nil
}

// Addbuttons adds a button as a child of this device
func (d *Device) AddButton(b *Button) error {
	if _, ok := d.Buttons[b.Address]; ok {
		return fmt.Errorf("button with address: %s already added to parent device", b.Address)
	}
	d.Buttons[b.Address] = b
	return nil
}

// AddDevice adds a device as a child of this device
func (d *Device) AddDevice(cd *Device) error {
	if _, ok := d.Devices[cd.Address]; ok {
		return fmt.Errorf("device with address: %s already added to parent device", cd.Address)
	}
	d.Devices[cd.Address] = cd
	return nil
}

// AddSensor adds a sensor as a child of this device
func (d *Device) AddSensor(s *Sensor) error {
	if _, ok := d.Sensors[s.Address]; ok {
		return fmt.Errorf("sensor with address: %s already added to device", s.Address)
	}
	d.Sensors[s.Address] = s
	return nil
}
