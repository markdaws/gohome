package main

import (
	"fmt"
	"time"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/www"
)

func main() {
	fmt.Println("hi")

	//	serverDone := make(chan bool)
	go func() {
		s := www.NewServer("./www")
		err := s.ListenAndServe(":8000")
		if err != nil {
			fmt.Println("error with server")
		}
		//close(serverDone)
	}()

	fmt.Println("creating system")
	dev := createSystem()
	err := dev.Connection.Connect()
	if err != nil {
		fmt.Println("failed to connect")
	} else {
		fmt.Println("connected")
	}

	r := &gohome.Recipe{
		Id:          "123",
		Name:        "Test",
		Description: "Test desc",
		Trigger: &gohome.TimeTrigger{
			Iterations: 5,
			Forever:    true,
			Interval:   time.Second * 10,
			At:         time.Now(),
		},
		Action: &gohome.FuncAction{Func: func() func() {
			var on bool = false
			return func() {
				if on {
					// Off
					dev.SetScene(&dev.Scenes[4])
				} else {
					// On
					dev.SetScene(&dev.Scenes[3])
				}
				on = !on
			}
		}()},
	}
	doneChan := r.Start()

	go func() {
		time.Sleep(time.Second * 60)
		fmt.Println("stopping")
		r.Stop()
	}()

	//What is the lifetime of a recipe? How to know when done?
	<-doneChan
	//	<-serverDone
}

func createSystem() *gohome.Device {
	system := gohome.System{gohome.Identifiable{
		Id:   "Home",
		Name: "My Home"}}

	var sbp gohome.Device
	sbp = gohome.Device{gohome.Identifiable{
		Id:   "sbp1",
		Name: "Lutron Smart Bridge Pro"},
		system,
		&gohome.TelnetConnection{
			Network:  "tcp",
			Address:  "192.168.0.10:23",
			Login:    "lutron",
			Password: "integration",
		},
		[]gohome.Scene{
			gohome.Scene{gohome.Identifiable{Id: "1",
				Name:        "All On",
				Description: "Turns on all the lights"}, &sbp,
				[]gohome.Command{&gohome.StringCommand{Value: "#DEVICE,1,1,3\r\n"}},
			},
			gohome.Scene{gohome.Identifiable{Id: "2",
				Name:        "All Off",
				Description: "Turns off all of the lights"}, &sbp,
				[]gohome.Command{&gohome.StringCommand{Value: "#DEVICE,2,2,3\r\n"}},
			},
			gohome.Scene{gohome.Identifiable{Id: "2",
				Name:        "Movie",
				Description: "Sets up movie mode"}, &sbp,
				[]gohome.Command{&gohome.StringCommand{Value: "#DEVICE,3,3,3\r\n"}},
			},
			gohome.Scene{gohome.Identifiable{Id: "2",
				Name:        "Front Door On",
				Description: "Turns front door lights on"}, &sbp,
				[]gohome.Command{&gohome.StringCommand{Value: "#DEVICE,1,6,3\r\n"}},
			},
			gohome.Scene{gohome.Identifiable{Id: "2",
				Name:        "Front Door Off",
				Description: "Turns front door lights off"}, &sbp,
				[]gohome.Command{&gohome.StringCommand{Value: "#DEVICE,1,7,3\r\n"}},
			},
		}}

	return &sbp
}
