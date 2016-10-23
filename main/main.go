package main

import (
	"fmt"
	"io/ioutil"

	eventExt "github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/api"
	"github.com/markdaws/gohome/intg"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/www"
)

type config struct {
	RecipeDirPath     string
	StartupConfigPath string
	WWWPort           string
	APIPort           string
}

func main() {

	config := config{
		StartupConfigPath: "/Users/mark/code/gohome/system2.json",
		WWWPort:           ":8000",
		APIPort:           ":5000",
	}

	// Processes all commands in the system in an async fashion, init with
	// 3 parallel workers and capacity to store up to 1000 commands to be processed
	cp := gohome.NewCommandProcessor(3, 1000)
	cp.Start()

	// Start the event bus
	eb := eventExt.NewBus(1000, 100)

	// Handles recipe management
	rm := gohome.NewRecipeManager(eb)

	//TODO: Remove, simulate user importing lutron information on load
	reset := true
	if reset {
		system := gohome.NewSystem("Lutron Smart Bridge Pro", "Lutron Smart Bridge Pro", cp, 1)
		intg.RegisterExtensions(system)

		bytes, err := ioutil.ReadFile("main/ip.json")
		if err != nil {
			panic("Could not read ip.json")
		}

		importer := system.Extensions.FindImporter(system, &gohome.Device{ModelNumber: "l-bdgpro2-wh"})
		if importer == nil {
			panic("Failed to import: " + err.Error())
		}
		err = importer.FromString(system, string(bytes[:]), "l-bdgpro2-wh")
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
	if err == store.ErrFileNotFound {
		log.V("startup file not found at: %s, creating new system", config.StartupConfigPath)

		// First time running the system, create a new blank system, save it
		system := gohome.NewSystem("Lutron Smart Bridge Pro", "Lutron Smart Bridge Pro", cp, 1)
		intg.RegisterExtensions(system)

		system.SavePath = config.StartupConfigPath
		err = store.SaveSystem(system, rm)
		if err != nil {
			panic("Failed to save initial system: " + err.Error())
		}
		sys = system
	} else if err != nil {
		panic("Failed to load system: " + err.Error())
	}

	sys.SavePath = config.StartupConfigPath
	cp.SetSystem(sys)
	sys.EvtBus = eb

	log.V("Initing devices...")
	sys.InitDevices()

	log.V("Starting recipes...")
	for _, recipe := range sys.Recipes {
		rm.RegisterAndStart(recipe)
	}

	/* TODO:
	// Event logger used to log event to UI clients via websockets
	wsLogger := gohome.NewWSEventLogger(sys)
	eb.AddConsumer(wsLogger)*/
	var wsLogger gohome.WSEventLogger

	done := make(chan bool)

	//TODO: Restart on fail
	go func() {
		log.V("WWW Server starting, listening on port %s", config.WWWPort)
		err := www.ListenAndServe("./www", config.WWWPort, sys, rm, wsLogger)
		if err != nil {
			fmt.Printf("error with WWW server, shutting down: %s\n", err)
		}
		close(done)
	}()

	//TODO: restart on fail
	go func() {
		log.V("API Server starting, listening on port %s", config.APIPort)
		err := api.ListenAndServe(config.APIPort, sys, rm, wsLogger)
		if err != nil {
			fmt.Printf("error with API server, shutting down: %s\n", err)
		}
	}()

	<-done
}
