package gohome

import (
	"strconv"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

type ChangeHandler interface {
	Update(monitorID string, b *ChangeBatch)
}

type ChangeBatch struct {
	Sensors map[string]SensorAttr
	Zones   map[string]cmd.Level
}

type MonitorGroup struct {
	Zones            map[string]bool
	Sensors          map[string]bool
	Handler          ChangeHandler
	TimeoutInSeconds int
}

type Monitor struct {
	system         *System
	nextID         int
	Groups         map[string]*MonitorGroup
	evtBus         *evtbus.Bus
	sensorToGroups map[string]map[string]bool
	zoneToGroups   map[string]map[string]bool
	sensorValues   map[string]SensorAttr
	zoneValues     map[string]cmd.Level
}

func NewMonitor(
	sys *System,
	evtBus *evtbus.Bus,
	sensorValues map[string]SensorAttr,
	zoneValues map[string]cmd.Level,
) *Monitor {

	//Callers can pass in initial values if they know what they are
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

func (m *Monitor) Start() {
}

func (m *Monitor) Refresh(monitorID string) {
	group, ok := m.Groups[monitorID]
	if !ok {
		// bad id, or expired ignore
		return
	}

	var changeBatch = &ChangeBatch{
		Sensors: make(map[string]SensorAttr),
		Zones:   make(map[string]cmd.Level),
	}

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
		group.Handler.Update(monitorID, changeBatch)
	}
	if len(sensorReport.SensorIDs) > 0 {
		m.evtBus.Enqueue(sensorReport)
	}
	if len(zoneReport.ZoneIDs) > 0 {
		m.evtBus.Enqueue(zoneReport)
	}
}

func (m *Monitor) Subscribe(g *MonitorGroup, refresh bool) string {

	//TODO: Called in multiple go routines?, mutex it

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

	if refresh {
		m.Refresh(monitorID)
	}

	// TODO: Where to set timeout, heap with timeout, set timer for smallest time
	// TODO: expirationHeap
	return monitorID
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

	//fmt.Printf("got updated sensor[%s] value: %+v\n", sensorID, attr)
	for groupID := range groups {
		group := m.Groups[groupID]
		//TODO: Don't make blocking...
		//TODO: Batch?
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
		//TODO: Don't make blocking...
		//TODO: Batch?
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
	//TODO:
}

// ==================================

//TODO: Unsubscribe
