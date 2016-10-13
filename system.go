package gohome

import (
	"strconv"

	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

type System struct {
	Name        string
	Description string
	SavePath    string
	Devices     map[string]Device
	Scenes      map[string]*Scene
	Zones       map[string]*zone.Zone
	Buttons     map[string]*Button
	Recipes     map[string]*Recipe

	CmdProcessor CommandProcessor
	EventBroker  event.Broker

	nextGlobalID int
}

func NewSystem(name, desc string, cmdProcessor CommandProcessor, nextGlobalID int) *System {
	s := &System{
		Name:         name,
		Description:  desc,
		Devices:      make(map[string]Device),
		Scenes:       make(map[string]*Scene),
		Zones:        make(map[string]*zone.Zone),
		Buttons:      make(map[string]*Button),
		Recipes:      make(map[string]*Recipe),
		CmdProcessor: cmdProcessor,
		nextGlobalID: nextGlobalID,
	}
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
		d := d
		go s.InitDevice(&d)
	}
}

func (s *System) InitDevice(d *Device) error {
	log.V("Init Device: %s", d)
	// If the device requires a connection pool, init all of the connections
	if d.Connections != nil {
		log.V("%s init connections", d)
		err := d.Connections.Init()
		if err != nil {
			log.E("%s failed to init connection pool: %s", d, err)
			return err
		} else {
			log.V("%s connected", d)
		}
	}

	if s.EventBroker != nil {
		s.EventBroker.AddProducer(d)
	}
	return nil
}

func (s *System) AddButton(b *Button) {
	//TODO: Validate button
	s.Buttons[b.ID] = b
}

func (s *System) AddDevice(d Device) error {
	errors := d.Validate()
	if errors != nil {
		return errors
	}

	//TODO: What about address, allow duplicates?
	//TODO: Add device need to init connections?
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

func (s *System) DeleteDevice(d Device) {
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
	return &dev
}
