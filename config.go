package gohome

import (
	"net"

	"github.com/markdaws/gohome/log"
)

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Config struct {
	// SystemPath is a path to the json file containing all of the system information
	SystemPath string `json:"systemPath"`

	// EventLogPath is the path where the event log will be written
	EventLogPath string `json:"eventLogPath"`

	// AutomationPath is the path where all the automation files live
	AutomationPath string `json:"automationPath"`

	// WWWAddr is the IP address for the WWW server
	WWWAddr string `json:"wwwAddr"`

	// WWWPort is the port for the WWW server
	WWWPort string `json:"wwwPort"`

	// UPNPNotifyAddr is the IP of the UPNP Notify server
	UPNPNotifyAddr string `json:"upnpNotifyAddr"`

	// UPNPNotifyPort is the port of the UPNP Notify server
	UPNPNotifyPort string `json:"upnpNotifyPort"`

	// Location specifies the lat/long where the home is physically located. This is needed
	// if you want to get accurate sunrise/sunset times
	Location location `json:"location"`
}

func (c *Config) Merge(cfg Config) {
	if c.SystemPath == "" {
		c.SystemPath = cfg.SystemPath
	}
	if c.EventLogPath == "" {
		c.EventLogPath = cfg.EventLogPath
	}
	if c.AutomationPath == "" {
		c.AutomationPath = cfg.AutomationPath
	}
	if c.WWWAddr == "" {
		c.WWWAddr = cfg.WWWAddr
	}
	if c.WWWPort == "" {
		c.WWWPort = cfg.WWWPort
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

// defaultConfig returns a default Config option with all the values
// populated to some default values
func NewDefaultConfig(systemPath string) *Config {
	addr := "127.0.0.1"
	useLocalhost := false
	if !useLocalhost {
		// Find the first public address we can bind to
		var err error
		addr, err = getIPV4NonLoopbackAddr()
		if err != nil || addr == "" {
			log.E("Could not detect non loopback IP address, falling back to 127.0.0.1")
		}
	}

	cfg := Config{
		SystemPath:     systemPath + "/gohome.json",
		EventLogPath:   systemPath + "/events.json",
		AutomationPath: systemPath + "/automation",
		WWWAddr:        addr,
		WWWPort:        "8000",
		UPNPNotifyAddr: addr,
		UPNPNotifyPort: "8001",
		Location:       location{},
	}

	return &cfg
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
