package gohome

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/markdaws/gohome/comm"
)

type Importer interface {
	ImportFromFile(path, importerID string, cmdProcessor CommandProcessor) (*System, error)
}

type importer struct {
}

func (i importer) ImportFromFile(path, importerID string, cp CommandProcessor) (*System, error) {
	//TODO: Support importing multiple devices
	switch importerID {
	case "L-BDGPRO2-WH":
		//TODO: sbpID, how pass this in? Need a bucket of params
		return importL_BDGPRO2_WH(path, "1", cp)
	default:
		return nil, errors.New("Unknown import type: " + importerID)
	}
}

func NewImporter() Importer {
	return importer{}
}

// Used for integration reports from Lutron Smart Bridge Pro
func importL_BDGPRO2_WH(integrationReportPath, smartBridgeProID string, cmdProcessor CommandProcessor) (*System, error) {

	//TODO: Handle non runtime panic
	bytes, err := ioutil.ReadFile(integrationReportPath)
	if err != nil {
		return nil, err
	}

	var configJson map[string]interface{}
	if err = json.Unmarshal(bytes, &configJson); err != nil {
		return nil, err
	}

	system := NewSystem("Lutron Smart Bridge Pro", "Lutron Smart Bridge Pro")

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

		device := NewDevice(
			deviceID,
			sys.NextGlobalID(),
			deviceName,
			deviceName,
			sys,
			cmdProcessor)

		for _, buttonMap := range deviceMap["Buttons"].([]interface{}) {
			button := buttonMap.(map[string]interface{})
			btnNumber := strconv.FormatFloat(button["Number"].(float64), 'f', 0, 64)

			var btnName string
			if name, ok := button["Name"]; ok {
				btnName = name.(string)
			} else {
				btnName = "Button " + btnNumber
			}

			b := &Button{
				LocalID:     btnNumber,
				GlobalID:    sys.NextGlobalID(),
				Name:        btnName,
				Description: btnName,
				Device:      device,
			}
			device.Buttons[btnNumber] = b
			system.AddButton(b)
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
				var pressCommand string = "#DEVICE," + deviceID + "," + buttonID + ",3\r\n"
				var releaseCommand string = "#DEVICE," + deviceID + "," + buttonID + ",4\r\n"

				var globalID = system.NextGlobalID()
				sceneContainer[globalID] = &Scene{
					LocalID:     buttonID,
					GlobalID:    globalID,
					Name:        buttonName,
					Description: buttonName,
					Commands: []Command{
						&StringCommand{
							Device:   sbp,
							Value:    pressCommand + releaseCommand,
							Friendly: "//TODO: Friendly",
							Type:     CTSystemSetScene,
						},
					},
					cmdProcessor: cmdProcessor,
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
			sbp.ConnectionInfo = comm.ConnectionInfo{
				Network:       "tcp",
				Address:       "192.168.0.10:23",
				Login:         "lutron",
				Password:      "integration",
				Stream:        true,
				PoolSize:      2,
				Authenticator: sbp,
			}
			makeScenes(system.Scenes, device, sbp)
			break
		}
	}

	if sbp == nil {
		return nil, errors.New("Did not find Smart Bridge Pro with ID:" + smartBridgeProID)
	}
	system.AddDevice(sbp)

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
		system.AddDevice(gohomeDevice)
		sbp.Devices[gohomeDevice.LocalID] = gohomeDevice
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
		var zoneTypeFinal ZoneType = ZTLight
		if zoneType, ok := zone["Type"].(string); ok {
			switch zoneType {
			case "light":
				zoneTypeFinal = ZTLight
			case "shade":
				zoneTypeFinal = ZTShade
			}
		}
		var outputTypeFinal OutputType = OTContinuous
		if outputType, ok := zone["Output"].(string); ok {
			switch outputType {
			case "binary":
				outputTypeFinal = OTBinary
			case "continuous":
				outputTypeFinal = OTContinuous
			}
		}
		z := &Zone{
			LocalID:     zoneID,
			GlobalID:    system.NextGlobalID(),
			Name:        zoneName,
			Description: zoneName,
			Type:        zoneTypeFinal,
			Output:      outputTypeFinal,
			setCommand: func(args ...interface{}) Command {
				return &StringCommand{
					Device:   sbp,
					Value:    "#OUTPUT," + zoneID + ",1,%.2f\r\n",
					Friendly: "//TODO: Friendly",
					Type:     CTZoneSetLevel,
					Args:     args,
				}
			},
			cmdProcessor: cmdProcessor,
		}
		system.AddZone(z)
		sbp.Zones[z.LocalID] = z
	}

	return system, nil
}
