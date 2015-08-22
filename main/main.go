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

	fmt.Println("creating system")
	var sbpID = "1"
	system, err := importSystem("main/ip.json", sbpID)
	if err != nil {
		fmt.Println("Failed to import", err)
		return
	}

	//TODO: When to connect, on demand or on start?
	//TODO: Automatic retry for failed commands?
	//TODO: Option to listen for events or not, permanent connection
	//TODO: Log disconnect/reconnect to devicesfor diagnostics purposes

	// TODO: Connection Pool, plus loop through connecting to devices? Only on demand?
	err = system.Devices[sbpID].Connection.Connect()
	if err != nil {
		panic("Failed to connect to device")
	} else {
		fmt.Println("connected")
	}

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
		Id:          "123",
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

func importSystem(integrationReportPath, smartBridgeProID string) (*gohome.System, error) {

	//TODO: dynamic path
	bytes, err := ioutil.ReadFile(integrationReportPath)
	if err != nil {
		return nil, err
	}

	//TODO: Rename x
	var x map[string]interface{}
	if err = json.Unmarshal(bytes, &x); err != nil {
		return nil, err
	}

	system := &gohome.System{
		Identifiable: gohome.Identifiable{
			Id:          "1",
			Name:        "Lutron Smart Bridge Pro",
			Description: "Lutron Smart Bridge Pro - imported //TODO: Date",
		},
		Devices: make(map[string]*gohome.Device),
		Scenes:  make(map[string]*gohome.Scene),
		Zones:   make(map[string]*gohome.Zone),
	}

	root, ok := x["LIPIdList"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Missing LIPIdList key, or value not a map")
	}
	devices, ok := root["Devices"].([]interface{})
	if !ok {
		return nil, errors.New("Missing Devices key, or value not a map")
	}

	fmt.Println("\nDEVICES")

	var makeDevice = func(deviceMap map[string]interface{}) *gohome.Device {
		var deviceID string = strconv.FormatFloat(deviceMap["ID"].(float64), 'f', 0, 64)
		var deviceName string = deviceMap["Name"].(string)

		return &gohome.Device{
			Identifiable: gohome.Identifiable{
				Id:          deviceID,
				Name:        deviceName,
				Description: deviceName},
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
						Id:          uniqueID,
						Name:        buttonName,
						Description: buttonName},
					Commands: []gohome.Command{&gohome.StringCommand{
						Device: sbp,
						Value:  pressCommand + releaseCommand,
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
			sbp = makeDevice(device)
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
		gohomeDevice := makeDevice(device)
		system.Devices[gohomeDevice.Id] = gohomeDevice
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
				Id:          zoneID,
				Name:        zoneName,
				Description: zoneName},
			SetCommand: &gohome.StringCommand{Device: sbp, Value: "#OUTPUT," + zoneID + ",1,%.2f\r\n"},
		}
	}

	return system, nil
}

/*
func createSystem() *gohome.System {

	sbp := &gohome.Device{gohome.Identifiable{
		Id:   "sbp1",
		Name: "Lutron Smart Bridge Pro"},
		&gohome.TelnetConnection{
			Network:  "tcp",
			Address:  "192.168.0.10:23",
			Login:    "lutron",
			Password: "integration",
		}}

	system := &gohome.System{Identifiable: gohome.Identifiable{
		Id:          "Home",
		Name:        "My Home",
		Description: "This is my home"},
		Devices: map[string]*gohome.Device{
			sbp.Id: sbp,
		},
		Zones: map[string]*gohome.Zone{
			"16": &gohome.Zone{Identifiable: gohome.Identifiable{
				Id:          "16",
				Name:        "Dining Area",
				Description: "The dining area light"},
				SetCommand: &gohome.StringCommand{Device: sbp, Value: "#OUTPUT,16,1,%.2f\r\n"},
				//TODO: Get? command how to know what the level will be
			},
			"23": &gohome.Zone{Identifiable: gohome.Identifiable{
				Id:          "23",
				Name:        "Living Room Shade",
				Description: "The shade over the main living room window"},
				SetCommand: &gohome.StringCommand{Device: sbp, Value: "#OUTPUT,23,1,%.2f\r\n"},
				//TODO: Get? command how to know what the level will be
			},
		},
		Scenes: map[string]*gohome.Scene{
			"1": &gohome.Scene{gohome.Identifiable{Id: "1",
				Name:        "All On",
				Description: "Turns on all the lights"},
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,1,3\r\n"}},
			},
			"2": &gohome.Scene{gohome.Identifiable{Id: "2",
				Name:        "All Off",
				Description: "Turns off all of the lights"},
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,2,3\r\n"}},
			},
			"3": &gohome.Scene{gohome.Identifiable{Id: "3",
				Name:        "Movie",
				Description: "Sets up movie mode"},
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,3,3\r\n"}},
			},
			"4": &gohome.Scene{gohome.Identifiable{Id: "4",
				Name:        "Front Door On",
				Description: "Turns front door lights on"},
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,6,3\r\n"}},
			},
			"5": &gohome.Scene{gohome.Identifiable{Id: "5",
				Name:        "Front Door Off",
				Description: "Turns front door lights off"},
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,7,3\r\n"}},
			},
			"6": &gohome.Scene{gohome.Identifiable{Id: "6",
				Name:        "Dining On",
				Description: "Turns dining area lights on"},
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,8,3\r\n#DEVICE,1,8,4\r\n"}},
			},
			"7": &gohome.Scene{gohome.Identifiable{Id: "7",
				Name:        "Dining Off",
				Description: "Turns dining area lights off"},
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,9,3\r\n#DEVICE,1,9,4\r\n"}},
			},
		}}

	//TODO: Commands
	// - Set a scene
	// - Create/Configure a scene
	// - Get list of current active scenes
	// - Set Zone Intensity
	// - Get Zone Intensity
	return system
}
*/
