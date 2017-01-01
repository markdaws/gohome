package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/go-home-iot/upnp"
	"github.com/markdaws/gohome/pkg/clock"
	"github.com/markdaws/gohome/pkg/gohome"
	"github.com/markdaws/gohome/pkg/log"
	"github.com/markdaws/gohome/pkg/store"
	"github.com/markdaws/gohome/pkg/www"
)

// This is injected by the build process and read from the VERSION file
var VERSION string

func main() {

	version := flag.Bool(
		"version",
		false,
		"View the version of the ghadmin too")

	configPath := flag.String(
		"config",
		"",
		"Specifies the path and file name to the goHOME config file")

	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	if configPath == nil || *configPath == "" {
		fmt.Println("You must specify the config option\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Load the system from disk
	sys, cfg := loadSystem(*configPath)

	// The event bus is the backbone of the app.  It allows device to post events
	// and other devices can list for events and act upon them.
	log.V("Event Bus - starting")
	eb := evtbus.NewBus(1000, 100)
	sys.Services.EvtBus = eb

	// Processes all commands in the system in an async fashion, init with
	// 3 parallel workers and capacity to store up to 1000 commands to be processed
	cp := gohome.NewCommandProcessor(sys, 3, 1000)
	cp.Start()
	sys.Services.CmdProcessor = cp

	// The UPNP service lets us listen for notifications from UPNP devices
	upnpService := upnp.NewSubServer()
	sys.Services.UPNP = upnpService
	go func() {
		for {
			endPoint := cfg.UPNPNotifyAddr + ":" + cfg.UPNPNotifyPort
			log.V("UPNP Service - listening on %s", endPoint)
			err := upnpService.Start(endPoint)
			log.E("upnp service crashed:" + err.Error())
			time.Sleep(time.Second * 5)
		}
	}()

	// Monitor is responsible for keeping track of all the current state values
	// for zones and sensors.  It listens on the event bus for changes so that
	// it can get the latest values
	monitor := gohome.NewMonitor(sys, sys.Services.EvtBus)
	sys.Services.Monitor = monitor

	// Log all of the events on the bus to the event log
	evtLogger := &gohome.EventLogger{Path: cfg.EventLogPath, Verbose: false}
	eb.AddConsumer(evtLogger)

	log.V("Initing devices...")
	sys.InitDevices()

	// TimeHelper helps fire events like sunrise/sunset that extensions and triggers
	// can use to fire events
	th := &gohome.TimeHelper{
		Time:      clock.SystemTime{},
		System:    sys,
		Latitude:  cfg.Location.Latitude,
		Longitude: cfg.Location.Longitude,
	}
	eb.AddProducer(th)

	sessions := gohome.NewSessions()
	go func() {
		for {
			endPoint := cfg.WWWAddr + ":" + cfg.WWWPort
			log.V("WWW Server starting, listening on %s", endPoint)
			err := www.ListenAndServe(cfg.WebUIPath, endPoint, sys, cfg.SystemPath, sessions, &cfg)
			log.E("error with WWW server, shutting down: %s\n", err)
			time.Sleep(time.Second * 5)
		}
	}()

	// Load all of the automation scripts
	autos, err := gohome.LoadAutomation(sys, cfg.AutomationPath)
	if err != nil {
		log.V("error loading automation scripts: %s", err)
	}
	for _, auto := range autos {
		auto := auto

		// When the automation is triggered, fire off the actions
		auto.Triggered = func(actions *gohome.CommandGroup) {
			sys.Services.EvtBus.Enqueue(&gohome.AutomationTriggeredEvt{
				Name: auto.Name,
			})

			log.V("automation[%s] - trigger fired, enqueuing actions", auto.Name)
			sys.Services.CmdProcessor.Enqueue(*actions)
		}

		sys.AddAutomation(auto)
		if auto.Enabled {
			log.V("automation - starting: %s", auto.Name)
			eb.AddConsumer(auto)
		} else {
			log.V("automation - disabled: %s", auto.Name)
		}
	}

	// Log we started the system
	sys.Services.EvtBus.Enqueue(&gohome.ServerStartedEvt{})

	// Sit forever since we have started all the services
	var done chan bool
	<-done
}

func loadSystem(configPath string) (*gohome.System, gohome.Config) {
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Unable to open config file:", configPath, err)
		os.Exit(1)
	}

	decoder := json.NewDecoder(file)
	var cfg *gohome.Config
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
		os.Exit(1)
	}

	log.V("Config information: %#v", cfg)

	sys, err := store.LoadSystem(cfg.SystemPath)
	if err != nil {
		log.E("System file not found at: %s, run the ghadmin command to initialize an empty system file", cfg.SystemPath)
		os.Exit(1)
	}

	return sys, *cfg
}
