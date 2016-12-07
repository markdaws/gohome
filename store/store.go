package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/intg"
	"github.com/markdaws/gohome/log"

	errExt "github.com/pkg/errors"
)

// ErrFileNotFound is returned when the specified path cannot be found
var ErrFileNotFound = errors.New("file not found")

// LoadSystem loads a gohome data file from the specified path
func LoadSystem(path string) (*gohome.System, error) {

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

	sys := gohome.NewSystem(s.Name)
	intg.RegisterExtensions(sys)

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

		dev.Features = d.Features
		for _, f := range d.Features {
			// When deserializing from JSON, the int32 and float32 types are converted
			// to float64 so need to massage them back
			attr.FixJSON(f.Attrs)

			log.V("loaded feature: ID:%s, Name:%s, Address: %s, Type:%s",
				f.ID, f.Name, f.Address, f.Type)
			sys.AddFeature(f)
		}
	}

	// First we have to load each scene, but without the commands, since a scene could have a
	// sceneSet command referencing a scene which hasn't been loaded yet
	for _, scn := range s.Scenes {
		scene := &gohome.Scene{
			Address:     scn.Address,
			ID:          scn.ID,
			Name:        scn.Name,
			Description: scn.Description,
		}
		sys.AddScene(scene)
		log.V("loaded Scene: ID:%s, Name:%s, Address:%s, Managed:%t",
			scene.ID, scene.Name, scene.Address, scene.Managed,
		)
	}

	for _, scn := range s.Scenes {
		scene := sys.Scenes[scn.ID]

		scene.Commands = make([]cmd.Command, len(scn.Commands))
		for i, command := range scn.Commands {
			var finalCmd cmd.Command
			switch command.Type {
			case "sceneSet":
				scn, ok := sys.Scenes[command.Attributes["SceneID"].(string)]
				if !ok {
					return nil, fmt.Errorf("invalid scene ID: %s", command.Attributes["SceneID"].(string))
				}
				finalCmd = &cmd.SceneSet{
					ID:        command.ID,
					SceneID:   scn.ID,
					SceneName: scn.Name,
				}

			case "featureSetAttrs":
				interfaceAttrs := command.Attributes["attrs"].(map[string]interface{})

				// Need to convert these to map[string]*attr.Attribute, don't see an easy way to
				// do that, so just marshal then unmarshal back to the type we want since commands
				// are generic
				b, err := json.Marshal(interfaceAttrs)
				if err != nil {
					return nil, errExt.Wrap(err, "failed to retrieve attrs field")
				}

				attrs := make(map[string]*attr.Attribute)
				err = json.Unmarshal(b, &attrs)
				if err != nil {
					return nil, errExt.Wrap(err, "failed to unmarshal attrs")
				}
				attr.FixJSON(attrs)

				featureID := command.Attributes["featureId"].(string)
				f, ok := sys.Features[featureID]
				if !ok {
					return nil, fmt.Errorf("invalid feature ID: %s", featureID)
				}

				finalCmd = &cmd.FeatureSetAttrs{
					ID:          command.ID,
					FeatureID:   f.ID,
					FeatureName: f.Name,
					FeatureType: f.Type,
					Attrs:       attrs,
				}

			default:
				return nil, fmt.Errorf("unknown command type %s", command.Type)
			}
			scene.Commands[i] = finalCmd
		}
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

	return sys, nil
}

// SaveSystem saves the specified system to disk
func SaveSystem(savePath string, s *gohome.System) error {
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
		for j, sCmd := range scene.Commands {
			switch xCmd := sCmd.(type) {
			case *cmd.SceneSet:
				cmds[j] = commandJSON{
					ID:   xCmd.ID,
					Type: "sceneSet",
					Attributes: map[string]interface{}{
						"SceneID": xCmd.SceneID,
					},
				}

			case *cmd.FeatureSetAttrs:
				cmds[j] = commandJSON{
					ID:   xCmd.ID,
					Type: "featureSetAttrs",
					Attributes: map[string]interface{}{
						"featureId": xCmd.FeatureID,
						"attrs":     xCmd.Attrs,
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

		d.Features = device.Features

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

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(savePath, b, 0644)
	return err
}
