package gohome

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type System struct {
	ID          string
	Name        string
	Description string
	Devices     map[string]Device
	Scenes      map[string]*Scene
	Zones       map[string]*Zone
	Buttons     map[string]*Button

	nextGlobalID int
}

func NewSystem(name, desc string) *System {
	s := &System{
		Name:        name,
		Description: desc,
		Devices:     make(map[string]Device),
		Scenes:      make(map[string]*Scene),
		Zones:       make(map[string]*Zone),
		Buttons:     make(map[string]*Button),
	}
	s.ID = s.NextGlobalID()
	return s
}

func (s *System) NextGlobalID() string {
	gid := s.nextGlobalID
	s.nextGlobalID++
	return strconv.Itoa(gid)
}

func (s *System) AddDevice(d Device) {
	s.Devices[d.GlobalID()] = d
}

func (s *System) AddButton(b *Button) {
	s.Buttons[b.GlobalID] = b
}

func (s *System) AddZone(z *Zone) {
	s.Zones[z.GlobalID] = z
}

type systemJSON struct {
	Version      string       `json:"version"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	NextGlobalID int          `json:"nextGlobalId"`
	Scenes       []sceneJSON  `json:"scenes"`
	Devices      []deviceJSON `json:"devices"`
}

type buttonJSON struct {
	LocalID     string `json:"localId"`
	GlobalID    string `json:"globalId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type zoneJSON struct {
	LocalID     string `json:"localId"`
	GlobalID    string `json:"globalId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DeviceID    string `json:"deviceId"`
	Type        string `json:"type"`
	Output      string `json:"output"`
}

type sceneJSON struct {
	LocalID     string        `json:"localId"`
	GlobalID    string        `json:"globalId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Commands    []commandJSON `json:"commands"`
}

type deviceJSON struct {
	LocalID     string       `json:"localId"`
	GlobalID    string       `json:"globalId"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	ModelNumber string       `json:"modelNumber"`
	Buttons     []buttonJSON `json:"buttons"`
	Zones       []zoneJSON   `json:"zones"`
	DeviceIDs   []string     `json:"deviceIds"`
	//system
	//devices
	//connectioninfo?
	Stream bool `json:"stream"`
}

type commandJSON struct {
}

func (s *System) Save(path string) error {

	out := systemJSON{
		Version:      "1.0.0.0",
		Name:         s.Name,
		Description:  s.Description,
		NextGlobalID: s.nextGlobalID,
	}

	out.Scenes = make([]sceneJSON, len(s.Scenes))
	var i = 0
	for _, scene := range s.Scenes {
		out.Scenes[i] = sceneJSON{
			LocalID:     scene.LocalID,
			GlobalID:    scene.GlobalID,
			Name:        scene.Name,
			Description: scene.Description,
		}

		cmds := make([]commandJSON, len(scene.Commands))
		//TODO: Loop through each command encoding to json
		out.Scenes[i].Commands = cmds
		i++
	}

	i = 0
	out.Devices = make([]deviceJSON, len(s.Devices))
	for _, device := range s.Devices {
		d := deviceJSON{
			LocalID:     device.LocalID(),
			GlobalID:    device.GlobalID(),
			Name:        device.Name(),
			Description: device.Description(),
			ModelNumber: device.ModelNumber(),
			Stream:      device.Stream(),
		}

		d.Buttons = make([]buttonJSON, len(device.Buttons()))
		bi := 0
		for _, btn := range device.Buttons() {
			d.Buttons[bi] = buttonJSON{
				LocalID:     btn.LocalID,
				GlobalID:    btn.GlobalID,
				Name:        btn.Name,
				Description: btn.Description,
			}
			bi++
		}

		d.Zones = make([]zoneJSON, len(device.Zones()))
		zi := 0
		for _, z := range device.Zones() {
			d.Zones[zi] = zoneJSON{
				LocalID:     z.LocalID,
				GlobalID:    z.GlobalID,
				Name:        z.Name,
				Description: z.Description,
				DeviceID:    device.GlobalID(),
				Type:        z.Type.ToString(),
				Output:      z.Output.ToString(),
			}
			zi++
		}

		d.DeviceIDs = make([]string, len(device.Devices()))
		di := 0
		for _, dev := range device.Devices() {
			d.DeviceIDs[di] = dev.GlobalID()
			di++
		}
		out.Devices[i] = d
		i++
	}

	//Recipes

	b, err := json.Marshal(out)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0644)
	return err
}

func LoadSystem(path string) (*System, error) {
	//TODO: Implement
	//Need to make sure global lists in system are populated correctly
	return nil, nil
}
