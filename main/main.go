package main

import (
	"fmt"
	"time"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/www"
)

func main() {
	//TODO: get from config
	//TODO: Parse buttons and phantom buttons
	//TODO: Get levels for zone
	var sbpID = "1"
	system, err := gohome.NewImporter().ImportFromFile("main/ip.json", "L-BDGPRO2-WH")
	if err != nil {
		panic("Failed to import: " + err.Error())
		return
	}

	// TODO: Connection Pool, plus loop through connecting to devices? Only on demand?
	sbpDevice := system.Devices[sbpID]
	conn, err := sbpDevice.Connect()
	if err != nil {
		panic("Failed to connect to device")
	} else {
		fmt.Println("connected")
	}

	//TOO: Fix this, shouldn't be here
	//TODO: Should be an option to persist connection
	go func() {
		gohome.Stream(sbpDevice, conn)
		//TODO defer Close
	}()

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

	//TODO: RecipeLoader - store/load recipes from a separate file
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
		// These will have to be hand coded - //TODO: recompiled? How to be dynamic
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
					system.Scenes["1:9"].Execute()
				} else {
					// On front door
					system.Scenes["1:8"].Execute()
				}
				on = !on
			}
		}()},
	}
	_ = r2

	//TODO: Recipe - turn light off after certain amount of time e.g. turn off bathroom
	//lights after 30 minutes
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
		}()
	*/

	//<-doneChan
	<-serverDone
}
