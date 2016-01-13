package main

import (
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/www"
)

type config struct {
	RecipeDirPath string
	StartupFile   string
	WWWPort       string
}

func main() {

	config := config{
		RecipeDirPath: "/Users/mark/code/gohome/recipes/",
		StartupFile:   "main/ip.json",
		WWWPort:       ":8000",
	}

	// Processes all commands in the system
	cp := gohome.NewCommandProcessor()
	go cp.Process()

	//TODO: Remove, load from gohome config file
	var sbpID = "1"
	system, err := gohome.NewImporter().ImportFromFile(config.StartupFile, "L-BDGPRO2-WH", cp)
	if err != nil {
		panic("Failed to import: " + err.Error())
		return
	}

	// Processes events
	eb := gohome.NewEventBroker()
	eb.Init()

	// Event logger used to log event to UI clients via websockets
	wsLogger := gohome.NewWSEventLogger()
	eb.AddConsumer(wsLogger)

	//TODO: Loop through all devices
	sbpDevice := system.Devices[sbpID]
	go func() {
		sbpDevice.InitConnections()
		eb.AddProducer(sbpDevice)
	}()

	// Load all of the recipes from disk, start listening
	rm := &gohome.RecipeManager{System: system}
	rm.Init(eb, config.RecipeDirPath)

	// Start www server
	done := make(chan bool)
	go func() {
		s := www.NewServer("./www", system, rm, wsLogger)
		//TODO
		err := s.ListenAndServe(config.WWWPort)
		if err != nil {
			fmt.Println("error with server")
		}
		close(done)
	}()
	<-done
}
