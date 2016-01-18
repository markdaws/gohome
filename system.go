package gohome

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/markdaws/gohome/comm"
)

type System struct {
	Name         string
	Description  string
	SavePath     string
	Devices      map[string]Device
	Scenes       map[string]*Scene
	Zones        map[string]*Zone
	Buttons      map[string]*Button
	Recipes      map[string]*Recipe
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
		Recipes:     make(map[string]*Recipe),
	}
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

func (s *System) AddScene(scn *Scene) {
	s.Scenes[scn.GlobalID] = scn
}

func (s *System) AddRecipe(r *Recipe) {
	s.Recipes[r.ID] = r
}

type systemJSON struct {
	Version      string       `json:"version"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	NextGlobalID int          `json:"nextGlobalId"`
	Scenes       []sceneJSON  `json:"scenes"`
	Devices      []deviceJSON `json:"devices"`
	Recipes      []recipeJSON `json:"recipes"`
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
	LocalID        string                    `json:"localId"`
	GlobalID       string                    `json:"globalId"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description"`
	ModelNumber    string                    `json:"modelNumber"`
	Buttons        []buttonJSON              `json:"buttons"`
	Zones          []zoneJSON                `json:"zones"`
	DeviceIDs      []string                  `json:"deviceIds"`
	ConnectionInfo *jsonTelnetConnectionInfo `json:"connectionInfo"`
	Stream         bool                      `json:"stream"`
}

type jsonTelnetConnectionInfo struct {
	PoolSize int    `json:"poolSize"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Network  string `json:"network"`
	Address  string `json:"address"`
}

type commandJSON struct {
}

func LoadSystem(path string, recipeManager *RecipeManager, commandProcessor CommandProcessor) (*System, error) {
	//TODO: Verify deviceIds exist etc

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var s systemJSON
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}

	sys := NewSystem(s.Name, s.Description)
	sys.nextGlobalID = s.NextGlobalID

	recipeManager.System = sys

	// Load all devices into global device list
	for _, d := range s.Devices {
		var ci *comm.TelnetConnectionInfo
		if d.ConnectionInfo != nil {
			//TODO: Only support telnet connectioninfo
			ci = &comm.TelnetConnectionInfo{
				Network:  d.ConnectionInfo.Network,
				Address:  d.ConnectionInfo.Address,
				Login:    d.ConnectionInfo.Login,
				Password: d.ConnectionInfo.Password,
				PoolSize: d.ConnectionInfo.PoolSize,
			}
		}

		dev := NewDevice(d.ModelNumber, d.LocalID, d.GlobalID, d.Name, d.Description, d.Stream, sys, commandProcessor, ci)
		if ci != nil {
			dev.ConnectionInfo().(*comm.TelnetConnectionInfo).Authenticator = dev
		}
		sys.AddDevice(dev)
	}

	// Have to go back through patching up devices to point to their child devices
	// since we only store device ID pointers in the JSON
	for _, d := range s.Devices {
		dev := sys.Devices[d.GlobalID]
		for _, dID := range d.DeviceIDs {
			childDev := sys.Devices[dID]
			dev.Devices()[childDev.LocalID()] = childDev
		}

		for _, btn := range d.Buttons {
			b := &Button{
				LocalID:     btn.LocalID,
				GlobalID:    btn.GlobalID,
				Name:        btn.Name,
				Description: btn.Description,
				Device:      dev,
			}
			dev.Buttons()[b.LocalID] = b
			sys.AddButton(b)
		}

		for _, zn := range d.Zones {
			z := &Zone{
				LocalID:     zn.LocalID,
				GlobalID:    zn.GlobalID,
				Name:        zn.Name,
				Description: zn.Description,
				Device:      dev,
				Type:        ZoneTypeFromString(zn.Type),
				Output:      OutputTypeFromString(zn.Output),
			}
			dev.Zones()[z.LocalID] = z
			sys.AddZone(z)
		}
	}

	for _, scn := range s.Scenes {
		scene := &Scene{
			LocalID:      scn.LocalID,
			GlobalID:     scn.GlobalID,
			Name:         scn.Name,
			Description:  scn.Description,
			cmdProcessor: commandProcessor,
			//TODO: commands
		}
		sys.AddScene(scene)
	}

	for _, r := range s.Recipes {
		rec, err := recipeManager.FromJSON(r)
		if err != nil {
			return nil, err
		}
		sys.AddRecipe(rec)
	}

	//TODO: Have to pass all the recipes into recipe manager after loading the system
	fmt.Printf("Loaded system! %+v\n", sys)
	return sys, nil
}

func (s *System) Save(recipeManager *RecipeManager) error {

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

		if ci, ok := device.ConnectionInfo().(*comm.TelnetConnectionInfo); ok && ci != nil {
			d.ConnectionInfo = &jsonTelnetConnectionInfo{
				PoolSize: ci.PoolSize,
				Login:    ci.Login,
				Password: ci.Password,
				Network:  ci.Network,
				Address:  ci.Address,
			}
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

	i = 0
	fmt.Printf("Saving RECIPES: %d\n", len(s.Recipes))
	out.Recipes = make([]recipeJSON, len(s.Recipes))
	for _, r := range s.Recipes {
		rec := recipeManager.ToJSON(r)
		fmt.Printf("JSON recipe %d %+v\n", i, rec)
		out.Recipes[i] = rec
		i++
	}
	fmt.Printf("%+v\n", out.Recipes)

	b, err := json.Marshal(out)
	if err != nil {
		return err
	}

	if s.SavePath == "" {
		return fmt.Errorf("SavePath is not set")
	}
	err = ioutil.WriteFile(s.SavePath, b, 0644)
	return err
}
