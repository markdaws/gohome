package main

import "github.com/markdaws/gohome/log"

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type config struct {
	// SystemPath is a path to the json file containing all of the system information
	SystemPath string `json:"systemPath"`

	// EventLogPath is the path where the event log will be written
	EventLogPath string `json::eventLogPath"`

	// WWWAddr is the IP address for the WWW server
	WWWAddr string `json:"wwwAddr"`

	// WWWPort is the port for the WWW server
	WWWPort string `json:"wwwPort"`

	// APIAddr is the IP address of the API server
	APIAddr string `json:"apiAddr"`

	// APIPort is the port number of the API server
	APIPort string `json:"apiPort"`

	// UPNPNotifyAddr is the IP of the UPNP Notify server
	UPNPNotifyAddr string `json:"upnpNotifyAddr"`

	// UPNPNotifyPort is the port of the UPNP Notify server
	UPNPNotifyPort string `json:"upnpNotifyPort"`

	// Location specifies the lat/long where the home is physically located. This is needed
	// if you want to get accurate sunrise/sunset times
	Location location `json:"location"`
}

func (c *config) Merge(cfg config) {
	if c.SystemPath == "" {
		c.SystemPath = cfg.SystemPath
	}
	if c.EventLogPath == "" {
		c.EventLogPath = cfg.EventLogPath
	}
	if c.WWWAddr == "" {
		c.WWWAddr = cfg.WWWAddr
	}
	if c.WWWPort == "" {
		c.WWWPort = cfg.WWWPort
	}
	if c.APIAddr == "" {
		c.APIAddr = cfg.APIAddr
	}
	if c.APIPort == "" {
		c.APIPort = cfg.APIPort
	}
	if c.UPNPNotifyAddr == "" {
		c.UPNPNotifyAddr = cfg.UPNPNotifyAddr
	}
	if c.UPNPNotifyPort == "" {
		c.UPNPNotifyPort = cfg.UPNPNotifyPort
	}
	if c.Location.Latitude == 0 && c.Location.Longitude == 0 {
		c.Location.Latitude = cfg.Location.Latitude
		c.Location.Longitude = cfg.Location.Longitude
	}
}

// defaultConfig returns a default config option with all the values
// populated to some default values
func NewDefaultConfig(systemPath string) *config {
	addr := "127.0.0.1"
	useLocalhost := true
	if !useLocalhost {
		// Find the first public address we can bind to
		var err error
		addr, err = getIPV4NonLoopbackAddr()
		if err != nil || addr == "" {
			log.E("Could not detect non loopback IP address, falling back to 127.0.0.1")
		}
	}

	cfg := config{
		SystemPath:     systemPath + "/gohome.json",
		EventLogPath:   systemPath + "/events.json",
		WWWAddr:        addr,
		WWWPort:        "8000",
		APIAddr:        addr,
		APIPort:        "5000",
		UPNPNotifyAddr: addr,
		UPNPNotifyPort: "8001",
		Location:       location{},
	}

	return &cfg
}
