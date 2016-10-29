package gohome

import (
	"strconv"
	"sync"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

// Updater is the interface for receiving updates values from the monitor
type Updater interface {
	Update(monitorID string, b *ChangeBatch)
}

// MonitorGroup represents a group of zones and sensors that a client wishes to
// receive updates for.
type MonitorGroup struct {
	Zones            map[string]bool
	Sensors          map[string]bool
	Handler          Updater
	TimeoutInSeconds int
}

// ChangeBatch contains a list of sensors and zones whos values have changed
type ChangeBatch struct {
	Sensors map[string]SensorAttr
	Zones   map[string]cmd.Level
}

// Monitor keeps track of the current zone and sensor values in the system and reports
// updates to clients
type Monitor struct {
	Groups map[string]*MonitorGroup

	system         *System
	nextID         int
	evtBus         *evtbus.Bus
	sensorToGroups map[string]map[string]bool
	zoneToGroups   map[string]map[string]bool
	sensorValues   map[string]SensorAttr
	zoneValues     map[string]cmd.Level
	mutex          sync.RWMutex
}

// NewMonitor returns an initialzed Monitor instance
func NewMonitor(
	sys *System,
	evtBus *evtbus.Bus,
	sensorValues map[string]SensorAttr,
	zoneValues map[string]cmd.Level,
) *Monitor {

	// Callers can pass in initial values if they know what they are
	// at the time of creating the monitor
	if sensorValues == nil {
		sensorValues = make(map[string]SensorAttr)
	}
	if zoneValues == nil {
		zoneValues = make(map[string]cmd.Level)
	}
	m := &Monitor{
		system:         sys,
		nextID:         1,
		Groups:         make(map[string]*MonitorGroup),
		sensorToGroups: make(map[string]map[string]bool),
		zoneToGroups:   make(map[string]map[string]bool),
		sensorValues:   sensorValues,
		zoneValues:     zoneValues,
		evtBus:         evtBus,
	}
	evtBus.AddConsumer(m)
	evtBus.AddProducer(m)
	return m
}

// Refresh causes the monitor to report the current values for any item in the
// monitor group, specified by the monitorID parameter
func (m *Monitor) Refresh(monitorID string) {
	m.mutex.RLock()
	group, ok := m.Groups[monitorID]
	m.mutex.RUnlock()

	if !ok {
		return
	}

	var changeBatch = &ChangeBatch{
		Sensors: make(map[string]SensorAttr),
		Zones:   make(map[string]cmd.Level),
	}

	// Build a list of sensors that need to report their values. If we
	// already have a value for a sensor we can just return that
	var sensorReport = &SensorsReport{}
	for sensorID := range group.Sensors {
		val, ok := m.sensorValues[sensorID]
		if ok {
			changeBatch.Sensors[sensorID] = val
		} else {
			sensorReport.Add(sensorID)
		}
	}

	var zoneReport = &ZonesReport{}
	for zoneID := range group.Zones {
		val, ok := m.zoneValues[zoneID]
		if ok {
			changeBatch.Zones[zoneID] = val
		} else {
			zoneReport.Add(zoneID)
		}
	}

	if len(changeBatch.Sensors) > 0 || len(changeBatch.Zones) > 0 {
		// We have some values already cached for certain items, return
		group.Handler.Update(monitorID, changeBatch)
	}
	if len(sensorReport.SensorIDs) > 0 {
		// Need to request these sensor values
		m.evtBus.Enqueue(sensorReport)
	}
	if len(zoneReport.ZoneIDs) > 0 {
		// Need to request these zone values
		m.evtBus.Enqueue(zoneReport)
	}
}

// Subscribe requests that the monitor keep track of updates for all of the zones
// and sensors listed in the MonitorGroup parameter. If refresh == true, the monitor
// will go and request values for all items in the monitor group, if false it won't
// until the caller calls the Subscribe method.  The function returns a monitorID value
// that can be passed into other functions, such as Unsubscribe and Refresh.
func (m *Monitor) Subscribe(g *MonitorGroup, refresh bool) string {

	m.mutex.Lock()
	monitorID := strconv.Itoa(m.nextID)
	m.nextID++
	m.Groups[monitorID] = g

	// Make sure we map from the zone and sensor ids back to this new group,
	// so that if any zones/snesor change in the future we know that we
	// need to alert this group
	for sensorID := range g.Sensors {
		// Get the monitor groups that are listening to this sensor
		groups, ok := m.sensorToGroups[sensorID]
		if !ok {
			groups = make(map[string]bool)
			m.sensorToGroups[sensorID] = groups
		}
		groups[monitorID] = true
	}
	for zoneID := range g.Zones {
		groups, ok := m.zoneToGroups[zoneID]
		if !ok {
			groups = make(map[string]bool)
			m.zoneToGroups[zoneID] = groups
		}
		groups[monitorID] = true
	}
	m.mutex.Unlock()

	if refresh {
		m.Refresh(monitorID)
	}

	// TODO: Where to set timeout, heap with timeout, set timer for smallest time
	// TODO: expirationHeap
	return monitorID
}

// Unsubscribe removes all references and updates for the specified monitorID
func (m *Monitor) Unsubscribe(monitorID string) {
	if _, ok := m.Groups[monitorID]; !ok {
		return
	}

	m.mutex.Lock()
	delete(m.Groups, monitorID)
	for sensorID, groups := range m.sensorToGroups {
		if _, ok := groups[monitorID]; ok {
			delete(groups, monitorID)
			if len(groups) == 0 {
				delete(m.sensorToGroups, sensorID)
				delete(m.sensorValues, sensorID)
			}
		}
	}
	for zoneID, groups := range m.zoneToGroups {
		if _, ok := groups[monitorID]; ok {
			delete(groups, monitorID)

			// If there are no groups pointed to by the zone, clean up
			// any refs to it
			if len(groups) == 0 {
				delete(m.zoneToGroups, zoneID)
				delete(m.zoneValues, zoneID)
			}
		}
	}
	m.mutex.Unlock()
}

func (m *Monitor) sensorAttrChanged(sensorID string, attr SensorAttr) {
	groups, ok := m.sensorToGroups[sensorID]
	if !ok {
		// Not a sensor we are monitoring, ignore
		return
	}

	// Is this value different to what we already know?
	currentVal, ok := m.sensorValues[sensorID]
	if ok {
		// No change, don't refresh clients
		if currentVal.Value == attr.Value {
			return
		}
	}
	m.sensorValues[sensorID] = attr

	for groupID := range groups {
		group := m.Groups[groupID]
		cb := &ChangeBatch{
			Sensors: make(map[string]SensorAttr),
		}
		cb.Sensors[sensorID] = attr
		group.Handler.Update(groupID, cb)
	}
}

func (m *Monitor) zoneLevelChanged(zoneID string, val cmd.Level) {
	groups, ok := m.zoneToGroups[zoneID]
	if !ok {
		return
	}

	// Is this value different to what we already know?
	currentVal, ok := m.zoneValues[zoneID]
	if ok {
		// No change, don't refresh clients
		if currentVal == val {
			return
		}
	}
	m.zoneValues[zoneID] = val

	for groupID := range groups {
		group := m.Groups[groupID]
		cb := &ChangeBatch{
			Zones: make(map[string]cmd.Level),
		}
		cb.Zones[zoneID] = val
		group.Handler.Update(groupID, cb)
	}
}

// ======= evtbus.Consumer interface

func (m *Monitor) ConsumerName() string {
	return "Monitor"
}

func (m *Monitor) StartConsuming(c chan evtbus.Event) {
	log.V("Monitor - start consuming events")

	go func() {
		for e := range c {
			switch evt := e.(type) {
			case *SensorAttrChanged:
				log.V("Monitor - processing SensorAttrChanged event")
				m.sensorAttrChanged(evt.SensorID, evt.Attr)

			case *SensorsReporting:
				for sensorID, attr := range evt.Sensors {
					m.sensorAttrChanged(sensorID, attr)
				}

			case *ZonesReporting:
				for zoneID, val := range evt.Zones {
					m.zoneLevelChanged(zoneID, val)
				}

			case *ZoneLevelChanged:
				m.zoneLevelChanged(evt.ZoneID, evt.Level)
			}

		}
		log.V("Monitor - event channel has closed")
	}()
}

func (m *Monitor) StopConsuming() {
	//TODO:
}

// =================================

// ======== evtbus.Producer interface
func (m *Monitor) ProducerName() string {
	return "Monitor"
}

func (m *Monitor) StartProducing(evtBus *evtbus.Bus) {
	//TODO: Delete?
}

func (m *Monitor) StopProducing() {
	//TODO: if a producer stops producing, do we need to invalidate all of the
	//values it is responsible for since they will not longer be updated??
}

// ==================================

//TODO: Unsubscribe
