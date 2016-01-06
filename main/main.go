package main

import (
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/www"
)

type config struct {
	RecipeDirPath string
	StartupFile   string
}

func main() {

	//TODO: Read config from a known location
	config := config{
		RecipeDirPath: "/Users/mark/code/gohome/recipes/",
		StartupFile:   "main/ip.json",
	}

	//TODO: Remove, load from gohome config file
	var sbpID = "1"
	system, err := gohome.NewImporter().ImportFromFile(config.StartupFile, "L-BDGPRO2-WH")
	if err != nil {
		panic("Failed to import: " + err.Error())
		return
	}

	// TODO: Connection Pool, plus loop through connecting to all devices
	sbpDevice := system.Devices[sbpID]
	eb := gohome.NewEventBroker()
	eb.Init()
	eb.AddProducer(sbpDevice)

	// Load all of the recipes from disk, start listening
	rm := &gohome.RecipeManager{System: system}
	rm.Init(eb, config.RecipeDirPath)

	// Event logger used to log event to UI clients via websockets
	l := &gohome.EventLogger{}
	eb.AddConsumer(l)

	// Start www server
	done := make(chan bool)
	go func() {
		s := www.NewServer("./www", system, rm, l)
		err := s.ListenAndServe(":8000")
		if err != nil {
			fmt.Println("error with server")
		}
		close(done)
	}()
	<-done
}
