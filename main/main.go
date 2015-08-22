package main

import (
	"fmt"
	"time"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/www"
)

func main() {

	fmt.Println("creating system")
	system := createSystem()

	// TODO: Connection Pool, plus loop through connecting to devices? Only on demand?
	err := system.Devices[0].Connection.Connect()
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

func createSystem() *gohome.System {

	//TODO: Read in from configuration file
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
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,8,3\r\n"}},
			},
			"7": &gohome.Scene{gohome.Identifiable{Id: "7",
				Name:        "Dining Off",
				Description: "Turns dining area lights off"},
				[]gohome.Command{&gohome.StringCommand{Device: sbp, Value: "#DEVICE,1,9,3\r\n"}},
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
