package gohome

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/markdaws/gohome/belkin"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/comm"
	"github.com/markdaws/gohome/zone"
)

type Importer interface {
	ImportFromFile(path, importerID string, cmdProcessor CommandProcessor) (*System, error)
}

type importer struct {
}

func (i importer) ImportFromFile(path, importerID string, cp CommandProcessor) (*System, error) {
	switch importerID {
	case "L-BDGPRO2-WH":
		//TODO: "1" should not be hard coded
		return importL_BDGPRO2_WH(path, "1", cp)
	default:
		return nil, errors.New("unknown import type: " + importerID)
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

	var configJSON map[string]interface{}
	if err = json.Unmarshal(bytes, &configJSON); err != nil {
		return nil, err
	}

	system := NewSystem("Lutron Smart Bridge Pro", "Lutron Smart Bridge Pro", cmdProcessor)

	root, ok := configJSON["LIPIdList"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Missing LIPIdList key, or value not a map")
	}
	devices, ok := root["Devices"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Devices key, or value not a map")
	}

	//TODO: Add logging
	var makeDevice = func(
		modelNumber, name, address string,
		deviceMap map[string]interface{},
		hub Device,
		sys *System,
		stream bool,
		auth *comm.Auth) Device {

		device := NewDevice(
			modelNumber,
			address,
			sys.NextGlobalID(),
			name,
			"",
			hub,
			stream,
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

			b := &Button{
				Address:     btnNumber,
				ID:          sys.NextGlobalID(),
				Name:        btnName,
				Description: btnName,
				Device:      device,
			}
			device.Buttons()[btnNumber] = b
			system.AddButton(b)
		}

		return device
	}

	var makeScenes = func(deviceMap map[string]interface{}, sbp Device) error {
		buttons, ok := deviceMap["Buttons"].([]interface{})
		if !ok {
			return errors.New("Missing Buttons key, or value not array")
		}

		//TODO: remove
		//var deviceID string = strconv.FormatFloat(deviceMap["ID"].(float64), 'f', 0, 64)
		for _, buttonMap := range buttons {
			button, ok := buttonMap.(map[string]interface{})
			if !ok {
				return errors.New("Expected Button elements to be objects")
			}
			if name, ok := button["Name"]; ok && !strings.HasPrefix(name.(string), "Button ") {
				//fmt.Printf("  Scene %d: %s\n", int(button["Number"].(float64)), name)

				var buttonID string = strconv.FormatFloat(button["Number"].(float64), 'f', 0, 64)
				var buttonName = button["Name"].(string)

				var btn = sbp.Buttons()[buttonID]
				scene := &Scene{
					Address:     buttonID,
					Name:        buttonName,
					Description: buttonName,
					Commands: []cmd.Command{
						&cmd.ButtonPress{
							ButtonAddress: btn.Address,
							ButtonID:      btn.ID,
							DeviceName:    sbp.Name(),
							DeviceAddress: sbp.Address(),
							DeviceID:      sbp.ID(),
						},
						&cmd.ButtonRelease{
							ButtonAddress: btn.Address,
							ButtonID:      btn.ID,
							DeviceName:    sbp.Name(),
							DeviceAddress: sbp.Address(),
							DeviceID:      sbp.ID(),
						},
					},
				}
				err := system.AddScene(scene)
				if err != nil {
					fmt.Printf("error adding scene: %s\n", err)
				}
			}
		}

		return nil
	}

	// First need to find the Smart Bridge Pro since it is needed to make scenes and zones
	var sbp Device
	for _, deviceMap := range devices {
		device, ok := deviceMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("Expected Devices elements to be objects")
		}

		var deviceID = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {

			sbp = makeDevice("L-BDGPRO2-WH", "Smart Bridge - Hub", "192.168.0.10:23", device, nil, system, true, &comm.Auth{
				Login:    "lutron",
				Password: "integration",
			})
			sbp.Auth().Authenticator = sbp

			// The smart bridge pro controls scenes by having phantom buttons that can be pressed,
			// each button activates a different scene. This means it really has two addresses, the
			// first is the IP address to talk to it, but then it also have the DeviceID which is needed
			// to press the buttons, so here, we make another device and assign the buttons to this
			// new device and use the lutron hub solely for communicating to.
			fmt.Printf("phantom id: %s\n", deviceID)
			sbpSceneDevice := makeDevice("", "Smart Bridge - Phantom Buttons", deviceID, device, sbp, system, false, nil)
			system.AddDevice(sbpSceneDevice)
			sbp.AddDevice(sbpSceneDevice)
			makeScenes(device, sbpSceneDevice)
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

		//fmt.Printf("%d: %s\n", int(device["ID"].(float64)), device["Name"])

		// Don't want to re-add the SBP
		var deviceID = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			continue
		}
		var deviceName string = device["Name"].(string)
		gohomeDevice := makeDevice("", deviceName, deviceID, device, sbp, system, false, nil)
		//TODO: Errors
		system.AddDevice(gohomeDevice)
		sbp.AddDevice(gohomeDevice)
	}

	zones, ok := root["Zones"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Zones key")
	}

	//fmt.Println("\nZONES")
	for _, zoneMap := range zones {
		z := zoneMap.(map[string]interface{})
		//fmt.Printf("%d: %s\n", int(z["ID"].(float64)), z["Name"])

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
			DeviceID:    sbp.ID(),
			Type:        zoneTypeFinal,
			Output:      outputTypeFinal,
		}
		//TODO: proper logging
		err := system.AddZone(newZone)
		if err != nil {
			fmt.Printf("err add zone: %s\n", err)
		} else {
			//fmt.Printf("added zone %s with ID %s\n", newZone.Name, newZone.ID)
		}
		//sbp.Zones()[newZone.Address] = newZone
	}

	//TODO: Move
	importConnectedByTCP(system)
	ghh := importGoHomeHub(system)
	importBelkin(ghh, system)
	return system, nil
}

//TODO: Temp function - import from UI
func importConnectedByTCP(system *System) {
	tcp := NewDevice(
		"TCP600GWB",
		"https://192.168.0.23",
		system.NextGlobalID(),
		"ConnectedByTcp Hub",
		"Description",
		nil,
		false,
		&comm.Auth{
			Token: "79tz3vbbop9pu5fcen60p97ix3mbvd3sblhjmz21",
		})

	zoneID := "216438039298518643"
	z := &zone.Zone{
		Address:     zoneID,
		Name:        "bulb1",
		Description: "tcp - bulb1",
		DeviceID:    tcp.ID(),
		Type:        zone.ZTLight,
		Output:      zone.OTContinuous,
		Controller:  "TCP - LED A19 11W",
	}
	system.AddDevice(tcp)
	system.AddZone(z)
}

func importGoHomeHub(system *System) Device {
	ghh := NewDevice(
		"GoHomeHub",
		"gohomehub",
		system.NextGlobalID(),
		"GoHome Hub",
		"GoHome Hub Description",
		nil,
		false,
		nil)
	system.AddDevice(ghh)

	z := &zone.Zone{
		Address:     "192.168.0.24:5577",
		Name:        "flux bulb",
		Description: "flux bulb 1",
		DeviceID:    ghh.ID(),
		Type:        zone.ZTLight,
		Output:      zone.OTRGB,
		Controller:  zone.ZCFluxWIFI,
	}
	system.AddZone(z)
	ghh.AddZone(z)
	return ghh
}

func importBelkin(ghh Device, system *System) {
	return
	responses, err := belkin.Scan(belkin.DTInsight, 5)
	fmt.Printf("got responses: %s\n", responses)
	fmt.Println(err)

	belkin.LoadDevice(responses[0])
	err = belkin.TurnOn(strings.Replace(responses[0].Location, "/setup.xml", "", -1))
	fmt.Println(err)
	time.Sleep(10 * time.Second)
	err = belkin.TurnOff(strings.Replace(responses[0].Location, "/setup.xml", "", -1))
	return
	/*
		z := &zone.Zone{
			Address:     location,
			Name:        "Belkin Insight Switch",
			Description: "Belkin",
			DeviceID:    ghh.ID(),
			Type:        zone.ZTOutlet,
			Output:      zone.OTBinary,
			Controller:  zone.ZCWeMoInsightSwitch,
		}
		system.AddZone(z)
		ghh.AddZone(z)
	*/
}
