package gohome

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/validation"
	"github.com/markdaws/gohome/zone"
)

type System struct {
	Name         string
	Description  string
	SavePath     string
	Devices      map[string]Device
	Scenes       map[string]*Scene
	Zones        map[string]*zone.Zone
	Buttons      map[string]*Button
	Recipes      map[string]*Recipe
	CmdProcessor CommandProcessor
	nextGlobalID int
	//TODO: Last modified
	//TODO: Add area concept
}

func NewSystem(name, desc string, cmdProcessor CommandProcessor) *System {
	s := &System{
		Name:         name,
		Description:  desc,
		Devices:      make(map[string]Device),
		Scenes:       make(map[string]*Scene),
		Zones:        make(map[string]*zone.Zone),
		Buttons:      make(map[string]*Button),
		Recipes:      make(map[string]*Recipe),
		CmdProcessor: cmdProcessor,
	}
	return s
}

func (s *System) NextGlobalID() string {
	gid := s.nextGlobalID
	s.nextGlobalID++
	return strconv.Itoa(gid)
}

func (s *System) AddButton(b *Button) {
	s.Buttons[b.ID] = b
}

func (s *System) AddDevice(d Device) error {
	errors := d.Validate()
	if errors != nil {
		return errors
	}

	//TODO: What about address, allow duplicates?
	//TODO: Add device need to init connections?
	s.Devices[d.ID()] = d
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

	// If this is a new zone we need to assign a new id
	if z.ID == "" {
		z.ID = s.NextGlobalID()
	}

	//TODO: When you add a zone, need to then also initconnetions
	//so that the system can talk to the new zone ...
	//TODO:!!
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
	Address     string `json:"address"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type zoneJSON struct {
	Address     string `json:"address"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DeviceID    string `json:"deviceId"`
	Type        string `json:"type"`
	Output      string `json:"output"`
	Controller  string `json:"controller"`
}

type sceneJSON struct {
	Address     string        `json:"address"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Commands    []commandJSON `json:"commands"`
}

type deviceJSON struct {
	Address     string       `json:"address"`
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	ModelNumber string       `json:"modelNumber"`
	Buttons     []buttonJSON `json:"buttons"`
	Zones       []zoneJSON   `json:"zones"`
	DeviceIDs   []string     `json:"deviceIds"`
	Auth        *authJSON    `json:"auth"`
	Stream      bool         `json:"stream"`
}

type authJSON struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type commandJSON struct {
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

func LoadSystem(path string, recipeManager *RecipeManager, cmdProcessor CommandProcessor) (*System, error) {
	//TODO: Verify deviceIds exist etc

	log.V("loading system from %s", path)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.V("failed to read file: %", err)
		return nil, err
	}

	var s systemJSON
	err = json.Unmarshal(b, &s)
	if err != nil {
		log.V("failed to unmarshal system json: %s", err)
		return nil, err
	}

	sys := NewSystem(s.Name, s.Description, cmdProcessor)
	sys.nextGlobalID = s.NextGlobalID

	recipeManager.System = sys

	// Load all devices into global device list
	for _, d := range s.Devices {
		var auth *comm.Auth
		if d.Auth != nil {
			auth = &comm.Auth{
				Login:    d.Auth.Login,
				Password: d.Auth.Password,
				Token:    d.Auth.Token,
			}
		}

		log.V("loaded Device: ID:%s, Name:%s, Model:%s, Address:%s", d.ID, d.Name, d.ModelNumber, d.Address)
		dev := NewDevice(d.ModelNumber, d.Address, d.ID, d.Name, d.Description, d.Stream, auth)
		if auth != nil {
			dev.Auth().Authenticator = dev
		}

		err = sys.AddDevice(dev)
		if err != nil {
			log.V("failed to add device to system: %s", err)
			return nil, err
		}
	}

	// Have to go back through patching up devices to point to their child devices
	// since we only store device ID pointers in the JSON
	for _, d := range s.Devices {
		dev := sys.Devices[d.ID]
		for _, dID := range d.DeviceIDs {
			//TODO: function
			childDev := sys.Devices[dID]
			dev.Devices()[childDev.Address()] = childDev
		}

		for _, btn := range d.Buttons {
			b := &Button{
				Address:     btn.Address,
				ID:          btn.ID,
				Name:        btn.Name,
				Description: btn.Description,
				Device:      dev,
			}

			//TODO: Add button function
			dev.Buttons()[b.Address] = b
			sys.AddButton(b)
		}

		for _, zn := range d.Zones {
			z := &zone.Zone{
				Address:     zn.Address,
				ID:          zn.ID,
				Name:        zn.Name,
				Description: zn.Description,
				DeviceID:    dev.ID(),
				Type:        zone.TypeFromString(zn.Type),
				Output:      zone.OutputFromString(zn.Output),
				Controller:  zn.Controller,
			}

			err := sys.AddZone(z)
			if err != nil {
				log.V("failed to add zone: %s", err)
				return nil, err
			}

			log.V("loaded Zone: ID:%s, Name:%s, Address:%s, Type:%s, Output:%s, Controller:%s",
				zn.ID, zn.Name, zn.Address, zn.Type, zn.Output, zn.Controller,
			)
		}
	}

	for _, scn := range s.Scenes {
		scene := &Scene{
			Address:     scn.Address,
			ID:          scn.ID,
			Name:        scn.Name,
			Description: scn.Description,
		}

		scene.Commands = make([]cmd.Command, len(scn.Commands))
		//TODO: Harden, check map access is ok
		for i, command := range scn.Commands {
			var finalCmd cmd.Command
			switch command.Type {
			case "zoneSetLevel":
				z := sys.Zones[command.Attributes["ZoneID"].(string)]
				finalCmd = &cmd.ZoneSetLevel{
					ZoneAddress: z.Address,
					ZoneID:      z.ID,
					ZoneName:    z.Name,
					Level:       cmd.Level{Value: float32(command.Attributes["Level"].(float64))},
				}
			case "buttonPress":
				btn := sys.Buttons[command.Attributes["ButtonID"].(string)]
				finalCmd = &cmd.ButtonPress{
					ButtonAddress: btn.Address,
					ButtonID:      btn.ID,
					DeviceName:    btn.Device.Name(),
					DeviceAddress: btn.Device.Address(),
					DeviceID:      btn.Device.ID(),
				}
			case "buttonRelease":
				btn := sys.Buttons[command.Attributes["ButtonID"].(string)]
				finalCmd = &cmd.ButtonRelease{
					ButtonAddress: btn.Address,
					ButtonID:      btn.ID,
					DeviceName:    btn.Device.Name(),
					DeviceAddress: btn.Device.Address(),
					DeviceID:      btn.Device.ID(),
				}
			case "sceneSet":
				scn := sys.Scenes[command.Attributes["SceneID"].(string)]
				finalCmd = &cmd.SceneSet{
					SceneID:   scn.ID,
					SceneName: scn.Name,
				}
			default:
				return nil, fmt.Errorf("unknown command type %s", command.Type)
			}
			scene.Commands[i] = finalCmd
		}
		err = sys.AddScene(scene)
		if err != nil {
			log.V("failed to add scene: %s", err)
			return nil, err
		}
		log.V("loaded Scene: ID:%s, Name:%s, Address:%s, Managed:%t",
			scene.ID, scene.Name, scene.Address, scene.Managed,
		)
	}

	for _, r := range s.Recipes {
		rec, err := recipeManager.FromJSON(r)
		if err != nil {
			return nil, err
		}
		sys.AddRecipe(rec)
	}
	return sys, nil
}

func (s *System) Save(recipeManager *RecipeManager) error {
	if s.SavePath == "" {
		return fmt.Errorf("SavePath is not set")
	}

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
			Address:     scene.Address,
			ID:          scene.ID,
			Name:        scene.Name,
			Description: scene.Description,
		}

		cmds := make([]commandJSON, len(scene.Commands))
		//TODO: Put this somewhere common? also in www
		for j, sCmd := range scene.Commands {
			switch xCmd := sCmd.(type) {
			case *cmd.ZoneSetLevel:
				cmds[j] = commandJSON{
					Type: "zoneSetLevel",
					Attributes: map[string]interface{}{
						"ZoneID": xCmd.ZoneID,
						"Level":  xCmd.Level.Value,
					},
				}
			case *cmd.ButtonPress:
				cmds[j] = commandJSON{
					Type: "buttonPress",
					Attributes: map[string]interface{}{
						"ButtonID": xCmd.ButtonID,
					},
				}
			case *cmd.ButtonRelease:
				cmds[j] = commandJSON{
					Type: "buttonRelease",
					Attributes: map[string]interface{}{
						"ButtonID": xCmd.ButtonID,
					},
				}
			case *cmd.SceneSet:
				cmds[j] = commandJSON{
					Type: "sceneSet",
					Attributes: map[string]interface{}{
						"SceneID": xCmd.SceneID,
					},
				}
			default:
				return fmt.Errorf("unknown command type")
			}
		}

		out.Scenes[i].Commands = cmds
		i++
	}

	i = 0
	out.Devices = make([]deviceJSON, len(s.Devices))
	for _, device := range s.Devices {
		d := deviceJSON{
			Address:     device.Address(),
			ID:          device.ID(),
			Name:        device.Name(),
			Description: device.Description(),
			ModelNumber: device.ModelNumber(),
			Stream:      device.Stream(),
		}

		if device.Auth() != nil {
			auth := device.Auth()
			d.Auth = &authJSON{
				Login:    auth.Login,
				Password: auth.Password,
				Token:    auth.Token,
			}
		}

		d.Buttons = make([]buttonJSON, len(device.Buttons()))
		bi := 0
		for _, btn := range device.Buttons() {
			d.Buttons[bi] = buttonJSON{
				Address:     btn.Address,
				ID:          btn.ID,
				Name:        btn.Name,
				Description: btn.Description,
			}
			bi++
		}

		d.Zones = make([]zoneJSON, len(device.Zones()))
		zi := 0
		for _, z := range device.Zones() {
			d.Zones[zi] = zoneJSON{
				Address:     z.Address,
				ID:          z.ID,
				Name:        z.Name,
				Description: z.Description,
				DeviceID:    device.ID(),
				Type:        z.Type.ToString(),
				Output:      z.Output.ToString(),
				Controller:  z.Controller,
			}
			zi++
		}

		d.DeviceIDs = make([]string, len(device.Devices()))
		di := 0
		for _, dev := range device.Devices() {
			d.DeviceIDs[di] = dev.ID()
			di++
		}
		out.Devices[i] = d
		i++
	}

	i = 0
	out.Recipes = make([]recipeJSON, len(s.Recipes))
	for _, r := range s.Recipes {
		rec := recipeManager.ToJSON(r)
		out.Recipes[i] = rec
		i++
	}

	b, err := json.Marshal(out)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.SavePath, b, 0644)
	return err
}

func (s *System) FromID(ID string) Device {
	return s.Devices[ID]
}
