package main

import (
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/event"
	"github.com/markdaws/gohome/imports"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/store"
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
		system, err := imports.FromFile("main/ip.json", "L-BDGPRO2-WH", cp)
		if err != nil {
			panic("Failed to import: " + err.Error())
		}

		system.SavePath = config.StartupConfigPath
		err = store.SaveSystem(system, rm)
		if err != nil {
			fmt.Println(err)
		}
	}

	sys, err := store.LoadSystem(config.StartupConfigPath, rm, cp)
	if err != nil {
		panic("Failed to load system: " + err.Error())
		//TODO: New systems, should have a blank system, create if not found
	}
	sys.SavePath = config.StartupConfigPath

	cp.SetSystem(sys)

	for _, d := range sys.Devices {
		d := d
		go func() {
			// If the device requires a connection pool, init all of the connections
			if d.Connections() != nil {
				log.V("%s init connections", d)
				err := d.Connections().Init()
				if err != nil {
					log.E("%s failed to init connection pool: %s", d, err)
				} else {
					log.V("%s connected", d)
				}
			}
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
