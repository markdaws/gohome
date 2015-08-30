package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/www"
)

func main() {
	//TODO: get from config
	//TODO: Parse buttons and phantom buttons
	//TODO: Get levels for zone
	var sbpID = "1"
	system, err := importSystem("main/ip.json", sbpID)
	if err != nil {
		panic("Failed to import: " + err.Error())
		return
	}

	// TODO: Connection Pool, plus loop through connecting to devices? Only on demand?
	sbpDevice := system.Devices[sbpID]
	err = sbpDevice.Connect()
	if err != nil {
		panic("Failed to connect to device")
	} else {
		fmt.Println("connected")
	}

	eventBroker := gohome.NewEventBroker()
	eventBroker.AddProducer(sbpDevice)

	//TODO: Add fmt printer consumer
	//TODO: Consumer that stores on AWS
	//TODO: Users should be able to specify consumers in a config file

	// Start www server
	serverDone := make(chan bool)
	go func() {
		s := www.NewServer("./www", system)
		err := s.ListenAndServe(":8000")
		if err != nil {
			fmt.Println("error with server")
		}
		close(serverDone)
	}()

	// How to codify this?
	// Triggers/Actions ... Entity can support one or more of either
	r := &gohome.Recipe{
		ID:          "123",
		Name:        "Test",
		Description: "Test desc",
		Trigger: &gohome.TimeTrigger{
			Iterations: 5,
			Forever:    true,
			Interval:   time.Second * 2,
			At:         time.Now(),
		},
		Action: &gohome.FuncAction{Func: func() func() {
			var on bool = false
			return func() {
				if on {
					// Off
					system.Scenes["5"].Execute()
				} else {
					// On
					system.Scenes["6"].Execute()
				}
				on = !on
			}
		}()},
	}
	_ = r

	r2 := &gohome.Recipe{
		ID:          "2",
		Name:        "Button Press Set Scene",
		Description: "Test desc",
		Trigger: &gohome.ButtonTrigger{
			MaxDuration: time.Duration(3) * time.Second,
			PressCount:  3,
		},
		Action: &gohome.FuncAction{Func: func() func() {
			var on bool = false
			return func() {
				if on {
					// Off front door
					system.Scenes["1:7"].Execute()
				} else {
					// On front door
					system.Scenes["1:6"].Execute()
				}
				on = !on
			}
		}()},
	}
	_ = r2

	bt, ok := r2.Trigger.(*gohome.ButtonTrigger)
	if ok {
		fmt.Println("got button trigger")
		eventBroker.AddConsumer(bt)
	}
	_ = r2.Start()

	//doneChan := r.Start()

	/*
		go func() {
			time.Sleep(time.Second * 60)
			fmt.Println("stopping")
			r.Stop()
		}()*/

	//<-doneChan
	<-serverDone
}

//TODO: Should be moved into own file
func importSystem(integrationReportPath, smartBridgeProID string) (*gohome.System, error) {

	bytes, err := ioutil.ReadFile(integrationReportPath)
	if err != nil {
		return nil, err
	}

	var configJson map[string]interface{}
	if err = json.Unmarshal(bytes, &configJson); err != nil {
		return nil, err
	}

	system := &gohome.System{
		Identifiable: gohome.Identifiable{
			ID:          "1",
			Name:        "Lutron Smart Bridge Pro",
			Description: "Lutron Smart Bridge Pro - imported //TODO: Date",
		},
		Devices: make(map[string]*gohome.Device),
		Scenes:  make(map[string]*gohome.Scene),
		Zones:   make(map[string]*gohome.Zone),
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

	var makeDevice = func(deviceMap map[string]interface{}, sys *gohome.System) *gohome.Device {
		var deviceID string = strconv.FormatFloat(deviceMap["ID"].(float64), 'f', 0, 64)
		var deviceName string = deviceMap["Name"].(string)

		return &gohome.Device{
			Identifiable: gohome.Identifiable{
				ID:          deviceID,
				Name:        deviceName,
				Description: deviceName},
			System: sys,
			//TODO: Shouldn't set here, comes in from user
			Connection: &gohome.TelnetConnection{
				Network:  "tcp",
				Address:  "192.168.0.10:23",
				Login:    "lutron",
				Password: "integration",
			}}
	}

	var makeScenes = func(sceneContainer map[string]*gohome.Scene, deviceMap map[string]interface{}, sbp *gohome.Device) error {
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

				sceneContainer[uniqueID] = &gohome.Scene{
					Identifiable: gohome.Identifiable{
						ID:          uniqueID,
						Name:        buttonName,
						Description: buttonName},
					Commands: []gohome.Command{&gohome.StringCommand{
						Device:   sbp,
						Value:    pressCommand + releaseCommand,
						Friendly: "//TODO: Friendly",
						Type:     gohome.CTSystemSetScene,
					}},
				}
			}
		}

		return nil
	}

	// First need to find the Smart Bridge Pro since it is needed to make scenes and zones
	var sbp *gohome.Device
	for _, deviceMap := range devices {
		device, ok := deviceMap.(map[string]interface{})
		if !ok {
			return nil, errors.New("Expected Devices elements to be objects")
		}
		var deviceID string = strconv.FormatFloat(device["ID"].(float64), 'f', 0, 64)
		if deviceID == smartBridgeProID {
			//ModelNumber: L-BDGPRO2-WH
			sbp = makeDevice(device, system)
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
		makeScenes(system.Scenes, device, sbp)
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
		system.Zones[zoneID] = &gohome.Zone{
			Identifiable: gohome.Identifiable{
				ID:          zoneID,
				Name:        zoneName,
				Description: zoneName},
			Type: gohome.ZTLight,
			SetCommand: &gohome.StringCommand{
				Device:   sbp,
				Value:    "#OUTPUT," + zoneID + ",1,%.2f\r\n",
				Friendly: "//TODO: Friendly",
				Type:     gohome.CTZoneSetLevel,
			},
		}
	}

	return system, nil
}
