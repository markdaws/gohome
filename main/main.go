package main

import (
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/www"
)

type config struct {
	RecipeDirPath     string
	StartupFile       string
	StartupConfigPath string
	WWWPort           string
}

func main() {

	config := config{
		RecipeDirPath:     "/Users/mark/code/gohome/recipes/",
		StartupFile:       "main/ip.json",
		StartupConfigPath: "/Users/mark/code/gohome/system.json",
		WWWPort:           ":8000",
	}

	// Processes all commands in the system in an async fashion
	cp := gohome.NewCommandProcessor()
	go cp.Process()

	//TODO: Remove, load from gohome config file
	system, err := gohome.NewImporter().ImportFromFile(config.StartupFile, "L-BDGPRO2-WH", cp)
	if err != nil {
		panic("Failed to import: " + err.Error())
		return
	}

	//TODO: Remove, testing
	err = system.Save(config.StartupConfigPath)
	if err != nil {
		fmt.Println(err)
	}

	// Processes events
	eb := gohome.NewEventBroker()
	eb.Init()

	// Event logger used to log event to UI clients via websockets
	wsLogger := gohome.NewWSEventLogger()
	eb.AddConsumer(wsLogger)

	for _, d := range system.Devices {
		if d.ConnectionInfo() != nil {
			d := d
			go func() {
				d.InitConnections()
				eb.AddProducer(d)
			}()
		}
	}

	// Load all of the recipes from disk, start listening
	rm := &gohome.RecipeManager{System: system}
	rm.Init(eb, config.RecipeDirPath)

	// Start www server
	done := make(chan bool)
	go func() {
		s := www.NewServer("./www", system, rm, wsLogger)
		err := s.ListenAndServe(config.WWWPort)
		if err != nil {
			fmt.Println("error with server")
		}
		close(done)
	}()
	<-done
}

//TODO: Recipes should be stored in the system config information, not in
//a separate file
