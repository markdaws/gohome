package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/intg"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/zone"
)

// ErrFileNotFound is returned when the specified path cannot be found
var ErrFileNotFound = errors.New("file not found")

// LoadSystem loads a gohome data file from the specified path
func LoadSystem(path string, recipeManager *gohome.RecipeManager) (*gohome.System, error) {

	log.V("loading system from %s", path)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, ErrFileNotFound
	}

	var s systemJSON
	err = json.Unmarshal(b, &s)
	if err != nil {
		log.V("failed to unmarshal system json: %s", err)
		return nil, err
	}

	sys := gohome.NewSystem(s.Name, s.Description)
	intg.RegisterExtensions(sys)
	recipeManager.System = sys

	// Load all devices into global device list
	for _, d := range s.Devices {
		var auth *gohome.Auth
		if d.Auth != nil {
			auth = &gohome.Auth{
				Login:    d.Auth.Login,
				Password: d.Auth.Password,
				Token:    d.Auth.Token,
			}
		}

		log.V("loaded Device: ID:%s, Name:%s, Model:%s, Address:%s", d.ID, d.Name, d.ModelNumber, d.Address)

		dev := gohome.NewDevice(
			d.ID,
			d.Name,
			d.Description,
			d.ModelNumber,
			d.ModelName,
			d.SoftwareVersion,
			d.Address,
			nil,
			nil,
			nil,
			auth)

		cmdBuilder := sys.Extensions.FindCmdBuilder(sys, dev)
		dev.CmdBuilder = cmdBuilder

		if d.ConnPool != nil {
			network := sys.Extensions.FindNetwork(sys, dev)
			if network == nil {
				return nil, fmt.Errorf("unsupported model number, no discoverer found: %s", d.ModelNumber)
			}

			connFactory, err := network.NewConnection(sys, dev)
			if err != nil {
				return nil, err
			}

			dev.Connections = pool.NewPool(pool.Config{
				Name:          d.ConnPool.Name,
				Size:          int(d.ConnPool.PoolSize),
				NewConnection: connFactory,

				//TODO: Need to store this in the system file, let imports decide this
				RetryDuration: time.Second * 10,
			})
		}

		sys.AddDevice(dev)
	}

	// Have to go back through patching up devices to point to their child devices
	// since we only store device ID pointers in the JSON
	for _, d := range s.Devices {
		dev := sys.Devices[d.ID]
		for _, dID := range d.DeviceIDs {
			childDev := sys.Devices[dID]
			dev.AddDevice(childDev)
		}

		// If the device has a hub we have to correctly set up that relationship
		if d.HubID != "" {
			hub, ok := sys.Devices[d.HubID]
			if !ok {
				return nil, fmt.Errorf("invalid hub ID: %s", d.HubID)
			}
			dev.Hub = hub
		}

		for _, btn := range d.Buttons {
			b := &gohome.Button{
				Address:     btn.Address,
				ID:          btn.ID,
				Name:        btn.Name,
				Description: btn.Description,
				Device:      dev,
			}

			dev.AddButton(b)
			sys.AddButton(b)
		}

		for _, zn := range d.Zones {
			z := &zone.Zone{
				Address:     zn.Address,
				ID:          zn.ID,
				Name:        zn.Name,
				Description: zn.Description,
				DeviceID:    dev.ID,
				Type:        zone.TypeFromString(zn.Type),
				Output:      zone.OutputFromString(zn.Output),
			}
			err := dev.AddZone(z)
			if err != nil {
				log.V("failed to add zone to device: %s", err)
			}
			sys.AddZone(z)
			log.V("loaded Zone: ID:%s, Name:%s, Address:%s, Type:%s, Output:%s",
				zn.ID, zn.Name, zn.Address, zn.Type, zn.Output,
			)
		}

		for _, sen := range d.Sensors {
			sensor := &gohome.Sensor{
				Address:     sen.Address,
				ID:          sen.ID,
				Name:        sen.Name,
				Description: sen.Description,
				DeviceID:    sen.DeviceID,
				Attr: gohome.SensorAttr{
					Name:     sen.Attr.Name,
					Value:    sen.Attr.Value,
					DataType: gohome.SensorDataType(sen.Attr.DataType),
					States:   sen.Attr.States,
				},
			}

			dev.AddSensor(sensor)
			sys.AddSensor(sensor)
			log.V("loaded Sensor: ID:%s, Name:%s, Address:%s, DeviceID:%s",
				sensor.ID, sensor.Name, sensor.Address, sensor.DeviceID,
			)
		}
	}

	for _, scn := range s.Scenes {
		scene := &gohome.Scene{
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
					ID:          command.ID,
					ZoneAddress: z.Address,
					ZoneID:      z.ID,
					ZoneName:    z.Name,
					Level: cmd.Level{
						Value: float32(command.Attributes["Level"].(float64)),
						R:     byte(command.Attributes["R"].(float64)),
						G:     byte(command.Attributes["G"].(float64)),
						B:     byte(command.Attributes["B"].(float64)),
					},
				}
			case "buttonPress":
				btn := sys.Buttons[command.Attributes["ButtonID"].(string)]
				finalCmd = &cmd.ButtonPress{
					ID:            command.ID,
					ButtonAddress: btn.Address,
					ButtonID:      btn.ID,
					DeviceName:    btn.Device.Name,
					DeviceAddress: btn.Device.Address,
					DeviceID:      btn.Device.ID,
				}
			case "buttonRelease":
				btn := sys.Buttons[command.Attributes["ButtonID"].(string)]
				finalCmd = &cmd.ButtonRelease{
					ID:            command.ID,
					ButtonAddress: btn.Address,
					ButtonID:      btn.ID,
					DeviceName:    btn.Device.Name,
					DeviceAddress: btn.Device.Address,
					DeviceID:      btn.Device.ID,
				}
			case "sceneSet":
				scn := sys.Scenes[command.Attributes["SceneID"].(string)]
				finalCmd = &cmd.SceneSet{
					ID:        command.ID,
					SceneID:   scn.ID,
					SceneName: scn.Name,
				}
			default:
				return nil, fmt.Errorf("unknown command type %s", command.Type)
			}
			scene.Commands[i] = finalCmd
		}
		sys.AddScene(scene)
		log.V("loaded Scene: ID:%s, Name:%s, Address:%s, Managed:%t",
			scene.ID, scene.Name, scene.Address, scene.Managed,
		)
	}

	for _, u := range s.Users {
		user := &gohome.User{
			ID:        u.ID,
			Login:     u.Login,
			HashedPwd: u.HashedPwd,
			Salt:      u.Salt,
		}
		sys.AddUser(user)
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

// SaveSystem saves the specified system to disk
func SaveSystem(savePath string, s *gohome.System, recipeManager *gohome.RecipeManager) error {
	out := systemJSON{
		Version:     "0.1.0",
		Name:        s.Name,
		Description: s.Description,
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
					ID:   xCmd.ID,
					Type: "zoneSetLevel",
					Attributes: map[string]interface{}{
						"ZoneID": xCmd.ZoneID,
						"Level":  xCmd.Level.Value,
						"R":      xCmd.Level.R,
						"G":      xCmd.Level.G,
						"B":      xCmd.Level.B,
					},
				}
			case *cmd.ButtonPress:
				cmds[j] = commandJSON{
					ID:   xCmd.ID,
					Type: "buttonPress",
					Attributes: map[string]interface{}{
						"ButtonID": xCmd.ButtonID,
					},
				}
			case *cmd.ButtonRelease:
				cmds[j] = commandJSON{
					ID:   xCmd.ID,
					Type: "buttonRelease",
					Attributes: map[string]interface{}{
						"ButtonID": xCmd.ButtonID,
					},
				}
			case *cmd.SceneSet:
				cmds[j] = commandJSON{
					ID:   xCmd.ID,
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
		hub := device.Hub
		var hubID = ""
		if hub != nil {
			hubID = hub.ID
		}

		var poolJSON *connPoolJSON
		if device.Connections != nil {
			config := device.Connections.Config
			poolJSON = &connPoolJSON{
				Name:     config.Name,
				PoolSize: int32(config.Size),
			}
		}
		d := deviceJSON{
			ID:              device.ID,
			Address:         device.Address,
			Name:            device.Name,
			Description:     device.Description,
			HubID:           hubID,
			ModelNumber:     device.ModelNumber,
			ModelName:       device.ModelName,
			SoftwareVersion: device.SoftwareVersion,
			ConnPool:        poolJSON,
		}

		if device.Auth != nil {
			auth := device.Auth
			d.Auth = &authJSON{
				Login:    auth.Login,
				Password: auth.Password,
				Token:    auth.Token,
			}
		}

		d.Buttons = make([]buttonJSON, len(device.Buttons))
		bi := 0
		for _, btn := range device.Buttons {
			d.Buttons[bi] = buttonJSON{
				Address:     btn.Address,
				ID:          btn.ID,
				Name:        btn.Name,
				Description: btn.Description,
			}
			bi++
		}

		d.Zones = make([]zoneJSON, len(device.Zones))
		zi := 0
		for _, z := range device.Zones {
			d.Zones[zi] = zoneJSON{
				Address:     z.Address,
				ID:          z.ID,
				Name:        z.Name,
				Description: z.Description,
				DeviceID:    device.ID,
				Type:        z.Type.ToString(),
				Output:      z.Output.ToString(),
			}
			zi++
		}

		d.Sensors = make([]sensorJSON, len(device.Sensors))
		si := 0
		for _, sen := range device.Sensors {
			d.Sensors[si] = sensorJSON{
				Address:     sen.Address,
				ID:          sen.ID,
				Name:        sen.Name,
				Description: sen.Description,
				DeviceID:    sen.DeviceID,
				Attr: sensorAttrJSON{
					Name:     sen.Attr.Name,
					Value:    sen.Attr.Value,
					DataType: string(sen.Attr.DataType),
					States:   sen.Attr.States,
				},
			}
			si++
		}

		d.DeviceIDs = make([]string, len(device.Devices))
		di := 0
		for _, dev := range device.Devices {
			d.DeviceIDs[di] = dev.ID
			di++
		}
		out.Devices[i] = d
		i++
	}

	i = 0
	out.Recipes = make([]gohome.RecipeJSON, len(s.Recipes))
	for _, r := range s.Recipes {
		rec := recipeManager.ToJSON(r)
		out.Recipes[i] = rec
		i++
	}

	i = 0
	out.Users = make([]userJSON, len(s.Users))
	for _, u := range s.Users {
		out.Users[i] = userJSON{
			ID:        u.ID,
			Login:     u.Login,
			HashedPwd: u.HashedPwd,
			Salt:      u.Salt,
		}
		i++
	}

	b, err := json.Marshal(out)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(savePath, b, 0644)
	return err
}
