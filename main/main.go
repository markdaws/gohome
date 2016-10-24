package main

import (
	"fmt"
	"io/ioutil"
	"net"

	eventExt "github.com/go-home-iot/event-bus"
	"github.com/go-home-iot/upnp"
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
	WWWAddr           string
	APIAddr           string
	UPNPNotifyAddr    string
}

func main() {

	addr, err := getIPV4NonLoopbackAddr()
	if err != nil || addr == "" {
		panic("could not find any address to bind to")
	}

	//TODO: Should read this from a config file on disk
	config := config{
		StartupConfigPath: "/Users/mark/code/gohome/system2.json",
		WWWAddr:           addr + ":8000",
		APIAddr:           addr + ":5000",
		UPNPNotifyAddr:    addr + ":8001",
	}

	// Processes all commands in the system in an async fashion, init with
	// 3 parallel workers and capacity to store up to 1000 commands to be processed
	cp := gohome.NewCommandProcessor(3, 1000)
	cp.Start()

	// The event bus is the backbone of the app.  It allows device to post events
	// and other devices can list for events and act upon them.
	log.V("Event Bus - starting")
	eb := eventExt.NewBus(1000, 100)

	// Log all of the events on the bus to the system log
	lc := &gohome.LogConsumer{}
	eb.AddConsumer(lc)

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
		system := gohome.NewSystem("My goHOME system", "", cp, 1)

		//TODO: RegisterExtensions expects that we have valid connections to the device
		// but we are initing afterwards ...
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

	//TODO: Seems janky setting these here, fix
	cp.SetSystem(sys)
	sys.EvtBus = eb

	go func() {
		for {
			//TODO: What happens if this crashes and all devices are waiting
			//for events, need to notify them to resubscribe ...
			upnpService := upnp.NewSubServer()
			sys.Services.UPNP = upnpService
			log.V("UPNP Service - listening on %s", config.UPNPNotifyAddr)
			err := upnpService.Start(config.UPNPNotifyAddr)
			log.E("upnp service crashed:" + err.Error())
		}
	}()

	// Init does things like connecting the gohome server to
	// all of the devices.
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

	go func() {
		for {
			log.V("WWW Server starting, listening on %s", config.WWWAddr)
			err := www.ListenAndServe("./www", config.WWWAddr, sys, rm, wsLogger)
			log.E("error with WWW server, shutting down: %s\n", err)
		}
	}()

	go func() {
		for {
			log.V("API Server starting, listening on %s", config.APIAddr)
			err := api.ListenAndServe(config.APIAddr, sys, rm, wsLogger)
			log.E("error with API server, shutting down: %s\n", err)
		}
	}()

	// Sit forever since we have started all the services
	// TODO: Graceful shutdown of service when receive control signal
	var done chan bool
	<-done
}

func getIPV4NonLoopbackAddr() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.To4() != nil &&
				!ip.IsLoopback() {
				return ip.To4().String(), nil
			}
		}
	}
	return "", nil
}
