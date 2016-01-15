package gohome

import "strconv"

type System struct {
	ID          string
	Name        string
	Description string
	Devices     map[string]*Device
	Scenes      map[string]*Scene
	Zones       map[string]*Zone
	Buttons     map[string]*Button
}

func NewSystem(name, desc string) *System {
	s := &System{
		Name:        name,
		Description: desc,
		Devices:     make(map[string]*Device),
		Scenes:      make(map[string]*Scene),
		Zones:       make(map[string]*Zone),
		Buttons:     make(map[string]*Button),
	}
	s.ID = s.NextGlobalID()
	return s
}

var globalID = 0

func (s *System) NextGlobalID() string {
	globalID++
	return strconv.Itoa(globalID)
}

func (s *System) AddDevice(d *Device) {
	s.Devices[d.GlobalID] = d
}

func (s *System) AddButton(b *Button) {
	s.Buttons[b.GlobalID] = b
}

func (s *System) AddZone(z *Zone) {
	s.Zones[z.GlobalID] = z
}
