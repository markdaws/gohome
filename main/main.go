package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/go-home-iot/upnp"
	"github.com/kardianos/osext"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/api"
	"github.com/markdaws/gohome/clock"
	"github.com/markdaws/gohome/intg"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/store"
	"github.com/markdaws/gohome/www"
)

func main() {

	runServer := flag.Bool("server", false, "run the goHOME server")
	setPassword := flag.Bool("set-password", false, "set the password for a user, creates a user is login nof found. --set-password mark mypwd")

	flag.Parse()

	if *setPassword {
		setPass(flag.Arg(0), flag.Arg(1))
		return
	}

	if *runServer {
		startServer()
		return
	}
}

func setPass(login, password string) {
	addedUser := false

	if login == "" || password == "" {
		fmt.Println("missing values, --set-password <login> <password>")
		return
	}

	// Load the system
	log.Silent = true
	sys, cfg := loadSystem()
	log.Silent = false

	// Add/update user
	var user *gohome.User
	for _, u := range sys.Users() {
		if u.Login == login {
			user = u
			break
		}
	}
	if user == nil {
		user = &gohome.User{
			ID:    sys.NewID(),
			Login: login,
		}
		err := user.Validate()
		if err != nil {
			fmt.Println("failed to add user", err)
			return
		}

		sys.AddUser(user)
		addedUser = true
	}

	err := user.SetPassword(password)
	if err != nil {
		fmt.Println("failed to set the password:", err)
		return
	}

	err = store.SaveSystem(cfg.SystemPath, sys)
	if err != nil {
		fmt.Println("Failed to save the user changes to disk: " + err.Error())
	}

	if addedUser {
		fmt.Println("Successfully added user: ", login)
	} else {
		fmt.Println("Successfully updated password for user: ", login)
	}
}

func startServer() {
	// Load the system from disk
	sys, cfg := loadSystem()

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
			err := www.ListenAndServe("./www", endPoint, sys, sessions)
			log.E("error with WWW server, shutting down: %s\n", err)
			time.Sleep(time.Second * 5)
		}
	}()

	go func() {
		for {
			endPoint := cfg.APIAddr + ":" + cfg.APIPort
			log.V("API Server starting, listening on %s", endPoint)
			err := api.ListenAndServe(cfg.SystemPath, endPoint, sys, sessions)
			log.E("error with API server, shutting down: %s\n", err)
			time.Sleep(time.Second * 5)
		}
	}()

	// Load all of the automation scripts
	autos, err := gohome.LoadAutomation(sys, cfg.AutomationPath)
	if err != nil {
		log.V("error loading automation scripts: %s", err)
	}
	for _, auto := range autos {
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

func loadSystem() (*gohome.System, config) {
	var cfg *config

	// Find the config file, if we can't find one, fall back to defaults
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		panic("Failed to locate the current executable directory")
	}

	// Try to read the config file, if we can't find one use defaults
	file, err := os.Open(folderPath + "/config.json")
	if err != nil {
		log.V("Error trying to open config.json [%s], falling back to defaults", folderPath)
		cfg = NewDefaultConfig(folderPath)
	} else {
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&cfg)
		if err != nil {
			log.V("Failed to parse config.json: %s, generating default config", err)
			cfg = NewDefaultConfig(folderPath)
		} else {
			defaultCfg := NewDefaultConfig(folderPath)

			// Merge the config and default config, user does not have to specify
			// all of the config values
			cfg.Merge(*defaultCfg)
		}
	}
	log.V("Config information: %#v", cfg)

	// Try to load an existing system from disk, if not present we will create
	// a blank one
	sys, err := store.LoadSystem(cfg.SystemPath)
	if err == store.ErrFileNotFound {
		log.V("startup file not found at: %s, creating new system", cfg.SystemPath)

		// First time running the system, create a new blank system, save it
		sys = gohome.NewSystem("My goHOME system")
		intg.RegisterExtensions(sys)

		err = store.SaveSystem(cfg.SystemPath, sys)
		if err != nil {
			panic("Failed to save initial system: " + err.Error())
		}
	} else if err != nil {
		panic("Failed to load system: " + err.Error())
	}

	return sys, *cfg
}

// getIPV4NonLoopbackAddr returns the first ipv4 non loopback address
// we can find. If non can be found an error is returned
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
