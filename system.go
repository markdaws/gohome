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
	UPNP    *upnp.SubServer
	Monitor *Monitor
}

type System struct {
	Name        string
	Description string
	SavePath    string
	Devices     map[string]*Device
	Scenes      map[string]*Scene
	Zones       map[string]*zone.Zone
	Buttons     map[string]*Button
	Sensors     map[string]*Sensor
	Recipes     map[string]*Recipe
	Extensions  *Extensions
	Services    SystemServices

	//TODO: Remove this, actions should not have execute or pass in cmdproc to Execute
	CmdProcessor CommandProcessor
	//TODO: Make sense here?
	EvtBus *evtbus.Bus

	nextGlobalID int
}

func NewSystem(name, desc string, cmdProcessor CommandProcessor, nextGlobalID int) *System {
	s := &System{
		Name:         name,
		Description:  desc,
		Devices:      make(map[string]*Device),
		Scenes:       make(map[string]*Scene),
		Zones:        make(map[string]*zone.Zone),
		Sensors:      make(map[string]*Sensor),
		Buttons:      make(map[string]*Button),
		Recipes:      make(map[string]*Recipe),
		CmdProcessor: cmdProcessor,
		nextGlobalID: nextGlobalID,
	}

	s.Extensions = NewExtensions()
	return s
}

func (s *System) NextGlobalID() string {
	gid := s.nextGlobalID
	s.nextGlobalID++
	return strconv.Itoa(gid)
}

func (s *System) PeekNextGlobalID() int {
	return s.nextGlobalID
}

func (s *System) InitDevices() {
	for _, d := range s.Devices {
		s.InitDevice(d)
	}
}

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
			s.EvtBus.AddProducer(evts.Producer)
		}
		if evts.Consumer != nil {
			log.V("%s - added event consumer", d)
			s.EvtBus.AddConsumer(evts.Consumer)
		}
	}

	return nil
}

func (s *System) AddButton(b *Button) {
	//TODO: Validate button
	s.Buttons[b.ID] = b
}

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

func (s *System) AddDevice(d *Device) error {
	errors := d.Validate()
	if errors != nil {
		return errors
	}

	//Should we add an id here, like we do in AddZone - consistency...

	//TODO: What about address, allow duplicates?
	s.Devices[d.ID] = d
	return nil
}

func (s *System) AddZone(z *zone.Zone) error {
	errors := z.Validate()
	if errors != nil {
		return errors
	}

	d, ok := s.Devices[z.DeviceID]
	if !ok {
		errors = &validation.Errors{}
		errors.Add("unknown device", "DeviceID")
		return errors
	}

	if z.ID == "" {
		z.ID = s.NextGlobalID()
	}

	//TODO: Don't do this here, confusing or make consistent across AddButton, AddSensor etc
	err := d.AddZone(z)
	if err != nil {
		return err
	}
	s.Zones[z.ID] = z
	return nil
}

func (s *System) AddScene(scn *Scene) error {
	errors := scn.Validate()
	if errors != nil {
		return errors
	}

	scn.ID = s.NextGlobalID()
	s.Scenes[scn.ID] = scn
	return nil
}

func (s *System) DeleteScene(scn *Scene) {
	delete(s.Scenes, scn.ID)
}

func (s *System) DeleteDevice(d *Device) {
	delete(s.Devices, d.ID)
	//TODO: Remove all associated zones + buttons
	//TODO: Need to stop all services, recipes, networking etc to this device
}

func (s *System) AddRecipe(r *Recipe) {
	s.Recipes[r.ID] = r
}

//TODO: Still needed?
func (s *System) FromID(ID string) *Device {
	dev := s.Devices[ID]
	return dev
}
