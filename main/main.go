package main

import (
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/www"
)

type config struct {
	RecipeDirPath     string
	StartupConfigPath string
	WWWPort           string
}

func main() {

	config := config{
		StartupConfigPath: "/Users/mark/code/gohome/system.json",
		WWWPort:           ":8000",
	}

	// Processes all commands in the system in an async fashion
	cp := gohome.NewCommandProcessor()
	go cp.Process()

	// Processes events
	eb := event.NewBroker()
	eb.Init()

	// Handles recipe management
	rm := gohome.NewRecipeManager(eb)

	//TODO: Remove
	reset := true
	if reset {
		system, err := gohome.NewImporter().ImportFromFile("main/ip.json", "L-BDGPRO2-WH", cp)
		if err != nil {
			panic("Failed to import: " + err.Error())
		}

		system.SavePath = config.StartupConfigPath
		err = system.Save(rm)
		if err != nil {
			fmt.Println(err)
		}
	}

	sys, err := gohome.LoadSystem(config.StartupConfigPath, rm, cp)
	sys.SavePath = config.StartupConfigPath
	if err != nil {
		panic("Failed to load system: " + err.Error())
		//TODO: New systems, should have a blank system, create if not found
	}

	cp.SetSystem(sys)

	for _, d := range sys.Devices {
		d := d
		go func() {
			d.InitConnections()
			eb.AddProducer(d)
		}()
	}

	//Start all the recipes
	log.V("Starting recipes...")
	for _, recipe := range sys.Recipes {
		rm.RegisterAndStart(recipe)
	}

	// Event logger used to log event to UI clients via websockets
	wsLogger := gohome.NewWSEventLogger(sys)
	eb.AddConsumer(wsLogger)

	// Start www server
	done := make(chan bool)
	go func() {
		log.V("WWW Server starting, listening on port %s", config.WWWPort)
		err := www.ListenAndServe("./www", config.WWWPort, sys, rm, wsLogger)
		if err != nil {
			fmt.Println("error with server")
		}
		close(done)
	}()
	<-done
}
