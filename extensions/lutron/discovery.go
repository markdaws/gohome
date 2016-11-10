package lutron

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/zone"
)

type discovery struct {
	System *gohome.System
}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	return []gohome.DiscovererInfo{gohome.DiscovererInfo{
		ID:          "lutron.l-bdgpro2-wh",
		Name:        "Lutron Smart Bridge Pro",
		Description: "Discover Lutron Smart Bridge Pro hubs",
		Type:        "FromString",
		PreScanInfo: "To get your configuration information, go to the Lutron app, then go: " +
			"Settings -> Advanced -> Integration -> Send Integration Report. Copy and paste the contents " +
			"of the email into the box below.  You also need the IP address of the Smart Bridge device, to find " +
			"that go to Settings -> Advanced -> Integration -> Network Settings. Then take note of the IP Address value " +
			"which you will need when you try to import the device",
	}}
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "lutron.l-bdgpro2-wh":
		return &discoverer{System: d.System}
	default:
		return nil
	}
}

type discoverer struct {
	System *gohome.System
}

func (d *discoverer) ScanDevices(sys *gohome.System) (*gohome.DiscoveryResults, error) {
	return nil, errors.New("unsupported")
}

func (d *discoverer) FromString(body string) (*gohome.DiscoveryResults, error) {
	//TODO: Fix should not suck into a system ...
	//TODO: remove
	//importer := &importer{System: d.System}
	//err := importer.FromString(d.System, body)
	//return nil, err

	result := &gohome.DiscoveryResults{}

	// We need to know which device is the Smart Bridge Pro - it is always ID==1 in the config file
	var smartBridgeProID string = "1"

	var configJSON map[string]interface{}
	if err := json.Unmarshal([]byte(body), &configJSON); err != nil {
		return nil, err
	}

	root, ok := configJSON["LIPIdList"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Missing LIPIdList key, or value not a map")
	}
	devices, ok := root["Devices"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Devices key, or value not a map")
	}

	var makeDevice = func(
		modelNumber,
		name,
		address string,
		deviceMap map[string]interface{},
		hub *gohome.Device,
		auth *gohome.Auth) *gohome.Device {

		device := gohome.NewDevice(
			modelNumber,
			"",
			"",
			address,
			"",
			name,
			"",
			hub,
			nil,
			nil,
			auth)

		for _, buttonMap := range deviceMap["Buttons"].([]interface{}) {
			button := buttonMap.(map[string]interface{})
			btnNumber := strconv.FormatFloat(button["Number"].(float64), 'f', 0, 64)

			var btnName string
			if name, ok := button["Name"]; ok {
				btnName = name.(string)
			} else {
				btnName = "Button " + btnNumber
			}

			b := &gohome.Button{
				Address:     btnNumber,
				Name:        btnName,
				Description: btnName,
				Device:      device,
			}
			b.ID = d.System.NextGlobalID()
			device.AddButton(b)
		}

		device.ID = d.System.NextGlobalID()
		return device
	}

	var makeScenes = func(deviceMap map[string]interface{}, sbp *gohome.Device) ([]*gohome.Scene, error) {
		var scenes = []*gohome.Scene{}
		buttons, ok := deviceMap["Buttons"].([]interface{})
		if !ok {
			return nil, errors.New("Missing Buttons key, or value not array")
		}

		for _, buttonMap := range buttons {
			button, ok := buttonMap.(map[string]interface{})
			if !ok {
				return nil, errors.New("Expected Button elements to be objects")
			}
			if name, ok := button["Name"]; ok && !strings.HasPrefix(name.(string), "Button ") {
				var buttonID string = strconv.FormatFloat(button["Number"].(float64), 'f', 0, 64)
				var buttonName = button["Name"].(string)

				var btn = sbp.Buttons[buttonID]
				scene := &gohome.Scene{
					Address:     buttonID,
					Name:        buttonName,
					Description: buttonName,
					Commands: []cmd.Command{
						&cmd.ButtonPress{
							ID:            d.System.NextGlobalID(),
							ButtonAddress: btn.Address,
							ButtonID:      btn.ID,
							DeviceName:    sbp.Name,
							DeviceAddress: sbp.Address,
							DeviceID:      sbp.ID,
						},
						&cmd.ButtonRelease{
							ID:            d.System.NextGlobalID(),
							ButtonAddress: btn.Address,
							ButtonID:      btn.ID,
							DeviceName:    sbp.Name,
							DeviceAddress: sbp.Address,
							DeviceID:      sbp.ID,
						},
					},
				}
				scene.ID = d.System.NextGlobalID()
				scenes = append(scenes, scene)
			}
		}

		return scenes, nil
	}

	// First need to find the Smart Bridge Pro since it is needed to make scenes and zones
	var sbp *gohome.Device
	for _, deviceMap := range devices {
		device, ok := deviceMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("Expected Devices elements to be objects")
		}

		var deviceID = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			dev := makeDevice(
				"l-bdgpro2-wh",
				"Smart Bridge - Hub",
				"",
				device,
				nil,
				&gohome.Auth{
					Login:    "lutron",
					Password: "integration",
				})

			// We need an address for this device, we can't get it from the Lutron config file so by setting
			// this flag the devices validation will fail when the user tries to import the device
			dev.AddressRequired = true
			sbp = dev

			//TODO: Needed for serialization?
			pool := pool.NewPool(pool.Config{
				Name: sbp.Name,
				Size: 2,
			})
			sbp.Connections = pool

			// The smart bridge pro controls scenes by having phantom buttons that can be pressed,
			// each button activates a different scene. This means it really has two addresses, the
			// first is the IP address to talk to it, but then it also have the DeviceID which is needed
			// to press the buttons, so here, we make another device and assign the buttons to this
			// new device and use the lutron hub solely for communicating to.
			sbpSceneDevice := makeDevice("", "Smart Bridge - Phantom Buttons", deviceID, device, sbp, nil)
			sbp.AddDevice(sbpSceneDevice)
			scenes, err := makeScenes(device, sbpSceneDevice)
			if err != nil {
				return nil, err
			}

			result.Devices = append(result.Devices, sbp, sbpSceneDevice)
			result.Scenes = append(result.Scenes, scenes...)
			break
		}
	}

	if sbp == nil {
		return nil, errors.New("Did not find Smart Bridge Pro with ID:" + smartBridgeProID)
	}

	for _, deviceMap := range devices {
		device, ok := deviceMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("Expected Devices elements to be objects")
		}

		// Don't want to re-add the SBP
		var deviceID = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			continue
		}
		var deviceName string = device["Name"].(string)
		gohomeDevice := makeDevice("", deviceName, deviceID, device, sbp, nil)
		result.Devices = append(result.Devices, gohomeDevice)
	}

	zones, ok := root["Zones"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Zones key")
	}

	for _, zoneMap := range zones {
		z := zoneMap.(map[string]interface{})

		var zoneID = strconv.FormatFloat(z["ID"].(float64), 'f', 0, 64)
		var zoneName = z["Name"].(string)
		var zoneTypeFinal = zone.ZTLight
		if zoneType, ok := z["Type"].(string); ok {
			switch zoneType {
			case "light":
				zoneTypeFinal = zone.ZTLight
			case "shade":
				zoneTypeFinal = zone.ZTShade
			}
		}
		var outputTypeFinal = zone.OTContinuous
		if outputType, ok := z["Output"].(string); ok {
			switch outputType {
			case "binary":
				outputTypeFinal = zone.OTBinary
			case "continuous":
				outputTypeFinal = zone.OTContinuous
			}
		}
		newZone := &zone.Zone{
			Address:     zoneID,
			Name:        zoneName,
			Description: zoneName,
			DeviceID:    sbp.ID,
			Type:        zoneTypeFinal,
			Output:      outputTypeFinal,
		}
		newZone.ID = d.System.NextGlobalID()
		err := sbp.AddZone(newZone)
		if err != nil {
			return nil, fmt.Errorf("err adding zone to device\n", err)
		}
	}

	return result, nil
}
