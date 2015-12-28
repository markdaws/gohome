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
	eventBroker := gohome.NewEventBroker()
	eventBroker.AddProducer(sbpDevice)

	// Load all of the recipes from disk, start listening
	recipeManager := &gohome.RecipeManager{System: system}
	recipeManager.Init(eventBroker, config.RecipeDirPath)

	// Start www server
	serverDone := make(chan bool)
	go func() {
		s := www.NewServer("./www", system, recipeManager)
		err := s.ListenAndServe(":8000")
		if err != nil {
			fmt.Println("error with server")
		}
		close(serverDone)
	}()
	<-serverDone

	//TODO: Automatically restat service if there is a crash for any reason
	//TODO: Harden everything
}

//TODO: Why increasing number of prints for a single event - investigate
