package gohome

import (
	"strconv"

	"github.com/go-home-iot/event-bus"
	"github.com/go-home-iot/upnp"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

// SystemServices is a collection of services that devices can access
// such as UPNP notification and discovery
type SystemServices struct {
	UPNP         *upnp.SubServer
	Monitor      *Monitor
	EvtBus       *evtbus.Bus
	CmdProcessor CommandProcessor
}

// System is a container that holds information such as all the zones and devices
// that have been created.
type System struct {
	Name         string
	Description  string
	Devices      map[string]*Device
	Scenes       map[string]*Scene
	Zones        map[string]*zone.Zone
	Buttons      map[string]*Button
	Sensors      map[string]*Sensor
	Recipes      map[string]*Recipe
	Extensions   *Extensions
	Services     SystemServices
	nextGlobalID int
}

// NewSystem returns an initial System instance.  It is still up to the caller
// to create all of the services and add them to the system after calling this function
func NewSystem(name, desc string, nextGlobalID int) *System {
	s := &System{
		Name:         name,
		Description:  desc,
		Devices:      make(map[string]*Device),
		Scenes:       make(map[string]*Scene),
		Zones:        make(map[string]*zone.Zone),
		Sensors:      make(map[string]*Sensor),
		Buttons:      make(map[string]*Button),
		Recipes:      make(map[string]*Recipe),
		nextGlobalID: nextGlobalID,
	}
	s.Extensions = NewExtensions()
	return s
}

// NextGlobalID returns the next unique global ID that can be used as an identifier
// for an item in the system.
func (s *System) NextGlobalID() string {
	gid := s.nextGlobalID
	s.nextGlobalID++
	return strconv.Itoa(gid)
}

// PeekNextGlobalID returns the next global ID that will be returned, but does not
// increment it
func (s *System) PeekNextGlobalID() int {
	return s.nextGlobalID
}

// InitDevices loops through all of the devices in the system and initializes them.
// This is async so after returning from this function the devices are still
// probably trying to initialize.  A device may try to create network connections
// or other tasks when it is initialized
func (s *System) InitDevices() {
	for _, d := range s.Devices {
		s.InitDevice(d)
	}
}

// InitDevice initializes a device.  If the device has a connection pool, this is
// when it will be initialized.  Also if the device produces or consumes events
// from the system bus, this is where it will be added to the event bus
func (s *System) InitDevice(d *Device) error {
	log.V("Init Device: %s", d)

	// If the device requires a connection pool, init all of the connections
	var done chan bool
	if d.Connections != nil {
		log.V("%s init connections", d)
		done = d.Connections.Init()
		_ = done
		log.V("%s connected", d)
	}

	evts := s.Extensions.FindEvents(s, d)
	if evts != nil {
		if evts.Producer != nil {
			log.V("%s - added event producer", d)
			s.Services.EvtBus.AddProducer(evts.Producer)
		}
		if evts.Consumer != nil {
			log.V("%s - added event consumer", d)
			s.Services.EvtBus.AddConsumer(evts.Consumer)
		}
	}

	return nil
}

// StopDevice stops the device, closes any network connections and any other services
// associated with the device
func (s *System) StopDevice(d *Device) {
	log.V("Stop Device: %s", d)

	if d.Connections != nil {
		d.Connections.Close()
	}

	evts := s.Extensions.FindEvents(s, d)
	if evts != nil {
		if evts.Producer != nil {
			s.Services.EvtBus.RemoveProducer(evts.Producer)
		}
		if evts.Consumer != nil {
			s.Services.EvtBus.RemoveConsumer(evts.Consumer)
		}
	}
}

// AddButton adds the button to the system and gives it a unique
// ID if it already doesn't have one
func (s *System) AddButton(b *Button) error {
	if b.ID == "" {
		b.ID = s.NextGlobalID()
	}
	s.Buttons[b.ID] = b
	return nil
}

// AddSensor adds a sensor to the system
func (s *System) AddSensor(sen *Sensor) error {
	errors := sen.Validate()
	if errors != nil {
		return errors
	}

	if sen.ID == "" {
		sen.ID = s.NextGlobalID()
	}
	s.Sensors[sen.ID] = sen
	return nil
}

// AddDevice adds a device to the system
func (s *System) AddDevice(d *Device) error {
	errors := d.Validate()
	if errors != nil {
		return errors
	}

	if d.ID == "" {
		d.ID = s.NextGlobalID()
	}
	s.Devices[d.ID] = d
	return nil
}

// AddZone adds a zone to the system
func (s *System) AddZone(z *zone.Zone) error {
	errors := z.Validate()
	if errors != nil {
		return errors
	}

	_, ok := s.Devices[z.DeviceID]
	if !ok {
		errors = &validation.Errors{}
		errors.Add("unknown device", "DeviceID")
		return errors
	}

	if z.ID == "" {
		z.ID = s.NextGlobalID()
	}
	s.Zones[z.ID] = z
	return nil
}

// AddScene adds a scene to the system
func (s *System) AddScene(scn *Scene) error {
	errors := scn.Validate()
	if errors != nil {
		return errors
	}

	if scn.ID == "" {
		scn.ID = s.NextGlobalID()
	}
	s.Scenes[scn.ID] = scn
	return nil
}

// DeleteScene deletes a scene from the system
func (s *System) DeleteScene(scn *Scene) {
	delete(s.Scenes, scn.ID)
}

// DeleteDevice deletes a device from the system and stops all associated
// services, for all zones and devices this is responsible for
func (s *System) DeleteDevice(d *Device) {
	delete(s.Devices, d.ID)
	//TODO: Remove all associated zones + buttons
	//TODO: Need to stop all services, recipes, networking etc to this device
}

// AddRecipe adds a recipe to the system
func (s *System) AddRecipe(r *Recipe) {

	if r.ID == "" {
		r.ID = s.NextGlobalID()
	}
	s.Recipes[r.ID] = r
}
