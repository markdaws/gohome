package gohome

import (
	"sync"

	"github.com/go-home-iot/event-bus"
	"github.com/go-home-iot/upnp"
	"github.com/markdaws/gohome/feature"
	"github.com/markdaws/gohome/log"
	"github.com/nu7hatch/gouuid"
)

// SystemServices is a collection of services that devices can access
// such as UPNP notification and discovery
type SystemServices struct {
	UPNP         *upnp.SubServer
	Monitor      *Monitor
	EvtBus       *evtbus.Bus
	CmdProcessor CommandProcessor
}

// System is a container that holds information such as all the zones and devices
// that have been created.
type System struct {
	Name        string
	Description string
	Area        *Area
	Extensions  *Extensions
	Services    SystemServices

	mutex      sync.RWMutex
	automation map[string]*Automation
	devices    map[string]*Device
	features   map[string]*feature.Feature
	scenes     map[string]*Scene
	users      map[string]*User
}

// NewSystem returns an initial System instance.  It is still up to the caller
// to create all of the services and add them to the system after calling this function
func NewSystem(name string) *System {
	s := &System{
		Name:        name,
		Description: "",
		automation:  make(map[string]*Automation),
		devices:     make(map[string]*Device),
		scenes:      make(map[string]*Scene),
		features:    make(map[string]*feature.Feature),
		users:       make(map[string]*User),
	}

	// Area is the root area which all of the devices and features are contained within
	s.Area = &Area{
		ID:   s.NewID(),
		Name: "Home",
	}

	s.Extensions = NewExtensions()
	return s
}

// NewID returns the next unique global ID that can be used as an identifier
// for an item in the system.
func (s *System) NewID() string {
	ID, err := uuid.NewV4()
	if err != nil {
		//TODO: Fail gracefully from this, keep looping?
		panic("failed to generate unique id in call to NextGlobalID:" + err.Error())
	}
	return ID.String()
}

// InitDevices loops through all of the devices in the system and initializes them.
// This is async so after returning from this function the devices are still
// probably trying to initialize.  A device may try to create network connections
// or other tasks when it is initialized
func (s *System) InitDevices() {
	s.mutex.RLock()
	for _, d := range s.devices {
		s.InitDevice(d)
	}
	s.mutex.RUnlock()
}

// InitDevice initializes a device.  If the device has a connection pool, this is
// when it will be initialized.  Also if the device produces or consumes events
// from the system bus, this is where it will be added to the event bus
func (s *System) InitDevice(d *Device) error {
	log.V("Init Device: %s", d)

	// If the device requires a connection pool, init all of the connections
	var done chan bool
	if d.Connections != nil {
		log.V("%s init connections", d)
		done = d.Connections.Init()
		_ = done
		log.V("%s connected", d)
	}

	evts := s.Extensions.FindEvents(s, d)
	if evts != nil {
		if evts.Producer != nil {
			log.V("%s - added event producer", d)
			s.Services.EvtBus.AddProducer(evts.Producer)
		}
		if evts.Consumer != nil {
			log.V("%s - added event consumer", d)
			s.Services.EvtBus.AddConsumer(evts.Consumer)
		}
	}

	return nil
}

// StopDevice stops the device, closes any network connections and any other services
// associated with the device
func (s *System) StopDevice(d *Device) {
	log.V("Stop Device: %s", d)

	if d.Connections != nil {
		d.Connections.Close()
	}

	evts := s.Extensions.FindEvents(s, d)
	if evts != nil {
		if evts.Producer != nil {
			s.Services.EvtBus.RemoveProducer(evts.Producer)
		}
		if evts.Consumer != nil {
			s.Services.EvtBus.RemoveConsumer(evts.Consumer)
		}
	}
}

// DeviceByID returns the device with the specified ID, nil if not found
func (s *System) DeviceByID(ID string) *Device {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.devices[ID]
}

// DeviceByAddress returns the first device found with the specified address
func (s *System) DeviceByAddress(addr string) *Device {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, dev := range s.devices {
		if dev.Address == addr {
			return dev
		}
	}
	return nil
}

// Devices returns a map of all the devices in the system, keyed by device ID
func (s *System) Devices() map[string]*Device {
	out := make(map[string]*Device)
	s.mutex.RLock()
	for k, v := range s.devices {
		out[k] = v
	}
	s.mutex.RUnlock()
	return out
}

// AddDevice adds a device to the system
func (s *System) AddDevice(d *Device) {
	s.mutex.Lock()
	s.devices[d.ID] = d
	s.mutex.Unlock()
}

// DeleteDevice deletes a device from the system and stops all associated
// services, for all zones and devices this is responsible for
func (s *System) DeleteDevice(d *Device) {
	s.mutex.Lock()
	delete(s.devices, d.ID)
	s.mutex.Unlock()

	//TODO: Remove all associated zones + buttons
	//TODO: Need to stop all services, recipes, networking etc to this device
}

// IsDupeDevice returns true if the device is a dupe of one the system
// already owns.  This check is not based on ID equality, since you could
// have scanned for a device previosuly and a second time, both scans will
// give a different ID for the device, since they are globally unique even
// though they are the same device
func (s *System) IsDupeDevice(x *Device) (*Device, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, y := range s.devices {
		if x.Address == y.Address {
			// Two devices are considered equal if they share the same address, however
			// if they have the same address but different hubs (if they have one) then
			// those are unique since the hub controls the device

			xHasHub := x.Hub != nil
			yHasHub := y.Hub != nil

			// Even though they have the same address, they have different hubs
			// so they are different
			if xHasHub != yHasHub {
				return nil, false
			}

			if !xHasHub && !yHasHub {
				// Both devices aren't under a hub and both have the same address so
				// they are a dupe
				return y, true
			}

			// If we are here both devices have the same address and they both have a hub
			// if the hubs are equal then the devices are equal
			return y, x.Hub.ID == y.Hub.ID
		}
	}
	return nil, false
}

// AddFeature adds a feature to the system. If a feature with the same ID already exists
// it is replaced with the new feature
func (s *System) AddFeature(f *feature.Feature) {
	s.mutex.Lock()
	s.features[f.ID] = f
	s.mutex.Unlock()
}

// FeatureByID returns the feature with the specified ID, nil if not found
func (s *System) FeatureByID(ID string) *feature.Feature {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.features[ID]
}

// FeatureByAID returns the feature with the specified automation ID, nil if not found
func (s *System) FeatureByAID(AID string) *feature.Feature {
	//TODO: Cache
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, f := range s.features {
		if f.AutomationID == AID {
			return f
		}
	}
	return nil
}

// FeaturesByType returns a map of all the features in the system that match the feature type, keyed
// by feature ID
func (s *System) FeaturesByType(ft string) map[string]*feature.Feature {
	//TODO: cache?
	features := make(map[string]*feature.Feature)
	s.mutex.RLock()
	for _, f := range s.features {
		if f.Type == ft {
			features[f.ID] = f
		}
	}
	s.mutex.RUnlock()
	return features
}

// SceneByID returns the scene with the specified ID, nil if not found
func (s *System) SceneByID(ID string) *Scene {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.scenes[ID]
}

// AddScene adds a scene to the system. If a scene with the same ID already exists, it is
// overwritten with the new scene
func (s *System) AddScene(scn *Scene) {
	s.mutex.Lock()
	s.scenes[scn.ID] = scn
	s.mutex.Unlock()
}

// Scenes returns a map of all the scenes in the system keyed by scene ID
func (s *System) Scenes() map[string]*Scene {
	out := make(map[string]*Scene)

	s.mutex.RLock()
	for k, v := range s.scenes {
		out[k] = v
	}
	s.mutex.RUnlock()
	return out
}

// DeleteScene deletes a scene from the system
func (s *System) DeleteScene(scn *Scene) {
	s.mutex.Lock()
	delete(s.scenes, scn.ID)
	s.mutex.Unlock()
}

// Automations reutrns all of the automation scripts keyed by ID
func (s *System) Automations() map[string]*Automation {
	out := make(map[string]*Automation)
	s.mutex.RLock()
	for k, v := range s.automation {
		out[k] = v
	}
	s.mutex.RUnlock()
	return out
}

// AddAutomation adds an automation instance to the system, indexed by name
func (s *System) AddAutomation(a *Automation) {
	s.mutex.Lock()
	s.automation[a.Name] = a
	s.mutex.Unlock()
}

// AutomationByID returns the automation instance with the specified ID, nil if not found
func (s *System) AutomationByID(ID string) *Automation {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.automation[ID]
}

// Users returns a map of all the users, keyed by user ID
func (s *System) Users() map[string]*User {
	out := make(map[string]*User)
	s.mutex.RLock()
	for k, v := range s.users {
		out[k] = v
	}
	s.mutex.RUnlock()
	return out
}

// AddUser adds the user to the system
func (s *System) AddUser(u *User) {
	s.mutex.Lock()
	s.users[u.ID] = u
	s.mutex.Unlock()
}
