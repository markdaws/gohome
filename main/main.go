package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

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
	WWWPort           string
	APIAddr           string
	APIPort           string
	UPNPNotifyAddr    string
	UPNPNotifyPort    string
}

func main() {
	//TODO: Don't panic, system should still start but with warning to the user

	useLocalhost := true
	var addr string
	if !useLocalhost {
		// Find the first public address we can bind to
		var err error
		addr, err = getIPV4NonLoopbackAddr()
		if err != nil || addr == "" {
			panic("could not find any address to bind to")
		}
	} else {
		addr = "127.0.0.1"
	}

	// TODO: Should read this from a config file on disk, only if ip addresses are
	// missing should we try to find one automatically
	config := config{
		StartupConfigPath: "/Users/mark/code/gohome/system2.json",
		WWWAddr:           addr,
		WWWPort:           "8000",
		APIAddr:           addr,
		APIPort:           "5000",
		UPNPNotifyAddr:    addr,
		UPNPNotifyPort:    "8001",
	}

	// Recipe manager handles processing all of the recipes the user has created
	rm := gohome.NewRecipeManager()

	//TODO: Remove, simulate user importing lutron information on load
	reset := true
	if reset {
		system := gohome.NewSystem("Lutron Smart Bridge Pro", "Lutron Smart Bridge Pro", 1)
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

		err = store.SaveSystem(config.StartupConfigPath, system, rm)
		if err != nil {
			fmt.Println(err)
		}
	}

	sys, err := store.LoadSystem(config.StartupConfigPath, rm)
	if err == store.ErrFileNotFound {
		log.V("startup file not found at: %s, creating new system", config.StartupConfigPath)

		// First time running the system, create a new blank system, save it
		sys = gohome.NewSystem("My goHOME system", "", 1)
		intg.RegisterExtensions(sys)

		err = store.SaveSystem(config.StartupConfigPath, sys, rm)
		if err != nil {
			panic("Failed to save initial system: " + err.Error())
		}
	} else if err != nil {
		panic("Failed to load system: " + err.Error())
	}

	// The event bus is the backbone of the app.  It allows device to post events
	// and other devices can list for events and act upon them.
	log.V("Event Bus - starting")
	eb := eventExt.NewBus(1000, 100)
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
			endPoint := config.UPNPNotifyAddr + ":" + config.UPNPNotifyPort
			log.V("UPNP Service - listening on %s", endPoint)
			err := upnpService.Start(endPoint)
			log.E("upnp service crashed:" + err.Error())
		}
	}()

	// Monitor is responsible for keeping track of all the current state values
	// for zones and sensors.  It listens on the event bus for changes so that
	// it can get the latest values
	monitor := gohome.NewMonitor(sys, sys.Services.EvtBus, nil, nil)
	sys.Services.Monitor = monitor

	// Log all of the events on the bus to the system log
	// TODO: Remove or json values so we can play back
	lc := &gohome.LogConsumer{}
	eb.AddConsumer(lc)

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
			endPoint := config.WWWAddr + ":" + config.WWWPort
			log.V("WWW Server starting, listening on %s", endPoint)
			err := www.ListenAndServe("./www", endPoint, sys, rm, wsLogger)
			log.E("error with WWW server, shutting down: %s\n", err)
			time.Sleep(time.Second * 5)
		}
	}()

	go func() {
		for {
			endPoint := config.APIAddr + ":" + config.APIPort
			log.V("API Server starting, listening on %s", endPoint)
			err := api.ListenAndServe(config.StartupConfigPath, endPoint, sys, rm, wsLogger)
			log.E("error with API server, shutting down: %s\n", err)
			time.Sleep(time.Second * 5)
		}
	}()

	// Sit forever since we have started all the services
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
