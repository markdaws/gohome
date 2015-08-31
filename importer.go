package gohome

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Importer interface {
	ImportFromFile(path, importerID string) (*System, error)
}

type importer struct {
}

func (i importer) ImportFromFile(path, importerID string) (*System, error) {
	//TODO: Support importing multiple types
	switch importerID {
	case "L-BDGPRO2-WH":
		//TODO: sbpID, how pass this in? Need a bucket of params
		return importL_BDGPRO2_WH(path, "1")
	default:
		return nil, errors.New("Unknown import type: " + importerID)
	}
}

func NewImporter() Importer {
	return importer{}
}

// Used for integration reports from Lutron Smart Bridge Pro
func importL_BDGPRO2_WH(integrationReportPath, smartBridgeProID string) (*System, error) {

	bytes, err := ioutil.ReadFile(integrationReportPath)
	if err != nil {
		return nil, err
	}

	var configJson map[string]interface{}
	if err = json.Unmarshal(bytes, &configJson); err != nil {
		return nil, err
	}

	system := &System{
		Identifiable: Identifiable{
			ID:          "1",
			Name:        "Lutron Smart Bridge Pro",
			Description: "Lutron Smart Bridge Pro - imported //TODO: Date",
		},
		Devices: make(map[string]*Device),
		Scenes:  make(map[string]*Scene),
		Zones:   make(map[string]*Zone),
	}

	root, ok := configJson["LIPIdList"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Missing LIPIdList key, or value not a map")
	}
	devices, ok := root["Devices"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Devices key, or value not a map")
	}

	fmt.Println("\nDEVICES")

	var makeDevice = func(deviceMap map[string]interface{}, sys *System) *Device {
		var deviceID string = strconv.FormatFloat(deviceMap["ID"].(float64), 'f', 0, 64)
		var deviceName string = deviceMap["Name"].(string)

		device := &Device{
			Identifiable: Identifiable{
				ID:          deviceID,
				Name:        deviceName,
				Description: deviceName},
			System:  sys,
			Buttons: make(map[string]*Button),
		}

		for _, buttonMap := range deviceMap["Buttons"].([]interface{}) {
			button := buttonMap.(map[string]interface{})
			btnNumber := strconv.FormatFloat(button["Number"].(float64), 'f', 0, 64)

			var btnName string
			if name, ok := button["Name"]; ok {
				btnName = name.(string)
			} else {
				btnName = "Button " + btnNumber
			}
			device.Buttons[btnNumber] = &Button{
				Identifiable: Identifiable{
					ID:          btnNumber,
					Name:        btnName,
					Description: btnName,
				},
				Device: device,
			}
		}

		return device
	}

	var makeScenes = func(sceneContainer map[string]*Scene, deviceMap map[string]interface{}, sbp *Device) error {
		buttons, ok := deviceMap["Buttons"].([]interface{})
		if !ok {
			return errors.New("Missing Buttons key, or value not array")
		}

		var deviceID string = strconv.FormatFloat(deviceMap["ID"].(float64), 'f', 0, 64)
		for _, buttonMap := range buttons {
			button, ok := buttonMap.(map[string]interface{})
			if !ok {
				return errors.New("Expected Button elements to be objects")
			}
			if name, ok := button["Name"]; ok && !strings.HasPrefix(name.(string), "Button ") {
				fmt.Printf("  Scene %d: %s\n", int(button["Number"].(float64)), name)

				var buttonID string = strconv.FormatFloat(button["Number"].(float64), 'f', 0, 64)
				var buttonName = button["Name"].(string)
				var uniqueID string = deviceID + ":" + buttonID
				var pressCommand string = "#DEVICE," + deviceID + "," + buttonID + ",3\r\n"
				var releaseCommand string = "#DEVICE," + deviceID + "," + buttonID + ",4\r\n"

				sceneContainer[uniqueID] = &Scene{
					Identifiable: Identifiable{
						ID:          uniqueID,
						Name:        buttonName,
						Description: buttonName},
					Commands: []Command{&StringCommand{
						Device:   sbp,
						Value:    pressCommand + releaseCommand,
						Friendly: "//TODO: Friendly",
						Type:     CTSystemSetScene,
					}},
				}
			}
		}

		return nil
	}

	// First need to find the Smart Bridge Pro since it is needed to make scenes and zones
	var sbp *Device
	for _, deviceMap := range devices {
		device, ok := deviceMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("Expected Devices elements to be objects")
		}

		var deviceID string = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			//ModelNumber: L-BDGPRO2-WH
			sbp = makeDevice(device, system)
			//TODO: Shouldn't set here, comes in from user
			sbp.ConnectionInfo = ConnectionInfo{
				Network:  "tcp",
				Address:  "192.168.0.10:23",
				Login:    "lutron",
				Password: "integration",
			}
			makeScenes(system.Scenes, device, sbp)
			break
		}
	}

	if sbp == nil {
		return nil, errors.New("Did not find Smart Bridge Pro with ID:" + smartBridgeProID)
	}
	system.Devices[smartBridgeProID] = sbp

	for _, deviceMap := range devices {
		device, ok := deviceMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("Expected Devices elements to be objects")
		}

		fmt.Printf("%d: %s\n", int(device["ID"].(float64)), device["Name"])

		// Don't want to re-add the SBP
		var deviceID string = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			continue
		}
		gohomeDevice := makeDevice(device, system)
		system.Devices[gohomeDevice.ID] = gohomeDevice
		//Only SBP has scenes that map to buttons, other devices are really buttons
		//makeScenes(system.Scenes, device, sbp)
	}

	zones, ok := root["Zones"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Zones key")
	}

	fmt.Println("\nZONES")
	for _, zoneMap := range zones {
		zone := zoneMap.(map[string]interface{})
		fmt.Printf("%d: %s\n", int(zone["ID"].(float64)), zone["Name"])

		var zoneID string = strconv.FormatFloat(zone["ID"].(float64), 'f', 0, 64)
		var zoneName string = zone["Name"].(string)
		system.Zones[zoneID] = &Zone{
			Identifiable: Identifiable{
				ID:          zoneID,
				Name:        zoneName,
				Description: zoneName},
			Type: ZTLight,
			SetCommand: &StringCommand{
				Device:   sbp,
				Value:    "#OUTPUT," + zoneID + ",1,%.2f\r\n",
				Friendly: "//TODO: Friendly",
				Type:     CTZoneSetLevel,
			},
		}
	}

	return system, nil
}
