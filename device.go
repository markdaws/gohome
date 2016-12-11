package gohome

import (
	"fmt"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/feature"
	"github.com/markdaws/gohome/validation"
)

// DeviceType explains the type of a device e.g. Dimmer or Shade
type DeviceType string

//TODO: Needed?
const (
	// DTDimmer - dimmer
	DTDimmer DeviceType = "dimmer"

	// DTSwitch - switch
	DTSwitch = "switch"

	// DTShade - shade
	DTShade = "shade"

	// DTHub - hub
	DTHub = "hub"

	// DTRemote - remote control
	DTRemote = "remote"

	// DTUnknown - unknown device type
	DTUnknown = "unknown"
)

// Auth contains authentication information such as login/password/security token
type Auth struct {
	Login    string
	Password string
	Token    string
}

// Device is a piece of hardware. It could be a dimmer, a shade, a remote etc
type Device struct {
	// ID a unique system wide ID, this is set when the device is added to the system
	// don't set this manually unless you know what you are doing
	ID string

	// Address can be whatever type of address is needed for the device e.g.
	// an IP address 192.168.0.9:23 or some value that is custom to whatever
	// system we have imported
	Address string

	// Name is a friendly name for the device, it will be shown in the UI
	Name string

	// Description provides more detailed information about the device
	Description string

	// Type describes what type of device this is e.g. Hub, Switch, Shade etc
	Type DeviceType

	// ModelName
	ModelName string

	// Modelnumber
	ModelNumber string

	// SoftwareVersion can store what version of software is installed on the device
	SoftwareVersion string

	// Features is a slice of features the device owns and exports
	// TODO: Hide behind getter/setter, mutex
	Features []*feature.Feature

	// CmdBuilder knows how to take an abstract command like ZoneSetLevel and turn
	// it in to specific commands for this particular piece of hardware.
	CmdBuilder cmd.Builder

	// Connections is optional, if the device needs a connection pool to communicate.
	Connections *pool.ConnectionPool

	// Auth - if authentication information is required to access the device, it is stored here
	Auth *Auth

	// Hub is the device that should be communicated with to talk to this device.  For example
	// you may have a keypad device, you don't talk directly to that but to some hub which
	// has a network address that then knows how to talk to the keypad.  Calling Hub will give
	// you that device.
	// TODO: Rename? Router/Gateway
	Hub *Device
}

// NewDevice returns an initialized device object
func NewDevice(
	ID,
	name,
	description,
	modelNumber,
	modelName,
	softwareVersion,
	address string,
	hub *Device,
	cmdBuilder cmd.Builder,
	connPool *pool.ConnectionPool,
	auth *Auth) *Device {

	dev := &Device{
		Address:         address,
		ModelNumber:     modelNumber,
		ModelName:       modelName,
		SoftwareVersion: softwareVersion,
		ID:              ID,
		Name:            name,
		Description:     description,
		Hub:             hub,
		Auth:            auth,
		CmdBuilder:      cmdBuilder,
		Connections:     connPool,
	}
	return dev
}

// Validate checks that all of the requirements for this to be a valid device are met
func (d *Device) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if d.ID == "" {
		errors.Add("required field", "ID")
	}

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
	return fmt.Sprintf("Device[ID:%s, Address:%s, Name: %s]", d.ID, d.Address, d.Name)
}

func (d *Device) AddFeature(f *feature.Feature) error {
	d.Features = append(d.Features, f)
	return nil
}

// OwnedFeature returns a slice of features that the device owns, where the
// map is keyed by feature.ID
func (d *Device) OwnedFeatures(featureIDs map[string]bool) []*feature.Feature {
	if len(d.Features) == 0 {
		return nil
	}

	var features []*feature.Feature
	for _, feature := range d.Features {
		if _, ok := featureIDs[feature.ID]; ok {
			features = append(features, feature)
		}
	}
	return features
}

// ButtonByAddress returns the button feature which has the matching address, nil
// if no matching button is found
func (d *Device) ButtonByAddress(addr string) *feature.Feature {
	for _, f := range d.Features {
		if f.Type == feature.FTButton && f.Address == addr {
			return f
		}
	}
	return nil
}

// FeatureByAddress returns the first feature that matches the specified address
// nil if no match is found
func (d *Device) FeatureTypeByAddress(t string, addr string) *feature.Feature {
	for _, f := range d.Features {
		if f.Type == t && f.Address == addr {
			return f
		}
	}
	return nil
}

// IsDupeFeature returns true if this is considered a duplicates feature, otherwise false. Features
// are considered equal, if they have the same type and address.
func (d *Device) IsDupeFeature(f *feature.Feature) bool {
	for _, ft := range d.Features {
		if ft.Type == f.Type && ft.Address == f.Address {
			return true
		}
	}
	return false
}
