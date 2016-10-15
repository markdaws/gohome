package main

import (
	"fmt"
	"io/ioutil"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/api"
	"github.com/markdaws/gohome/event"
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

	// Processes all commands in the system in an async fashion
	cp := gohome.NewCommandProcessor()
	go cp.Process()

	// Processes events
	eb := event.NewBroker()
	eb.Init()

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

		err = system.Extensions.Importers["l-bdgpro2-wh"].FromString(system, string(bytes[:]), "l-bdgpro2-wh")
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
	}

	sys.SavePath = config.StartupConfigPath
	cp.SetSystem(sys)

	log.V("Initing devices...")
	sys.EventBroker = eb
	sys.InitDevices()

	log.V("Starting recipes...")
	for _, recipe := range sys.Recipes {
		rm.RegisterAndStart(recipe)
	}

	// Event logger used to log event to UI clients via websockets
	wsLogger := gohome.NewWSEventLogger(sys)
	eb.AddConsumer(wsLogger)

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
		err := api.ListenAndServe("./www", config.APIPort, sys, rm, wsLogger)
		if err != nil {
			fmt.Printf("error with API server, shutting down: %s\n", err)
		}
	}()

	<-done
}
