package gohome

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

// MonitorDelegate is the interface for receiving updates values from the monitor
type MonitorDelegate interface {
	Update(b *ChangeBatch)
	Expired(monitorID string)
}

// MonitorGroup represents a group of zones and sensors that a client wishes to
// receive updates for.
type MonitorGroup struct {
	Zones           map[string]bool
	Sensors         map[string]bool
	Handler         MonitorDelegate
	Timeout         time.Duration
	timeoutAbsolute time.Time
	id              string
}

func (mg *MonitorGroup) String() string {
	return fmt.Sprintf("MonitorGroup[ID:%s, %d zones, %d sensors]", mg.id, len(mg.Zones), len(mg.Sensors))
}

// ChangeBatch contains a list of sensors and zones whos values have changed
type ChangeBatch struct {
	MonitorID string
	Sensors   map[string]SensorAttr
	Zones     map[string]cmd.Level
}

func (cb *ChangeBatch) String() string {
	return fmt.Sprintf("ChangeBatch[monitorID: %s, #zones:%d, #sensors:%d", cb.MonitorID, len(cb.Zones), len(cb.Sensors))
}

// Monitor keeps track of the current zone and sensor values in the system and reports
// updates to clients
type Monitor struct {
	groups         map[string]*MonitorGroup
	system         *System
	nextID         int64
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
		nextID:         time.Now().UnixNano(),
		groups:         make(map[string]*MonitorGroup),
		sensorToGroups: make(map[string]map[string]bool),
		zoneToGroups:   make(map[string]map[string]bool),
		sensorValues:   sensorValues,
		zoneValues:     zoneValues,
		evtBus:         evtBus,
	}

	m.handleTimeouts()
	evtBus.AddConsumer(m)
	evtBus.AddProducer(m)
	return m
}

// Refresh causes the monitor to report the current values for any item in the
// monitor group, specified by the monitorID parameter.  If force is true then
// the current cached values stored in the monitor are ignored and new values are
// requested
func (m *Monitor) Refresh(monitorID string, force bool) {
	m.mutex.RLock()
	group, ok := m.groups[monitorID]

	if !ok {
		m.mutex.RUnlock()
		return
	}

	var changeBatch = &ChangeBatch{
		MonitorID: monitorID,
		Sensors:   make(map[string]SensorAttr),
		Zones:     make(map[string]cmd.Level),
	}

	// Build a list of sensors that need to report their values. If we
	// already have a value for a sensor we can just return that
	var sensorReport = &SensorsReportEvt{}
	for sensorID := range group.Sensors {
		val, ok := m.sensorValues[sensorID]
		if ok && !force {
			changeBatch.Sensors[sensorID] = val
		} else {
			sensorReport.Add(sensorID)
		}
	}

	var zoneReport = &ZonesReportEvt{}
	for zoneID := range group.Zones {
		val, ok := m.zoneValues[zoneID]
		if ok && !force {
			changeBatch.Zones[zoneID] = val
		} else {
			zoneReport.Add(zoneID)
		}
	}

	log.V("Monitor - refreshing: %s, force:%t", group, force)
	log.V("Monitor - refreshing: cached values: [%s], uncached zones: %s, uncached sensors: %s",
		changeBatch, zoneReport, sensorReport)

	m.mutex.RUnlock()

	if len(changeBatch.Sensors) > 0 || len(changeBatch.Zones) > 0 {
		// We have some values already cached for certain items, return
		group.Handler.Update(changeBatch)
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

// InvalidateValues removes any cached values, for zones and sensors listed
// in the monitor group
func (m *Monitor) InvalidateValues(monitorID string) {
	m.mutex.RLock()
	group, ok := m.groups[monitorID]
	m.mutex.RUnlock()

	if !ok {
		return
	}

	log.V("Monitor - invalidate values: monitorID: %s", monitorID)
	m.mutex.Lock()
	for sensorID := range group.Sensors {
		delete(m.sensorValues, sensorID)
	}
	for zoneID := range group.Zones {
		delete(m.zoneValues, zoneID)
	}
	m.mutex.Unlock()
}

// Group returns the group for the specified ID if one exists
func (m *Monitor) Group(monitorID string) (*MonitorGroup, bool) {
	m.mutex.RLock()
	group, ok := m.groups[monitorID]
	m.mutex.RUnlock()
	return group, ok
}

// SubscribeRenew updates the timeout parameter for the group to increment to now() + timeout
// where timeout was specified in the initial call to Subscribe
func (m *Monitor) SubscribeRenew(monitorID string) error {
	m.mutex.RLock()
	group, ok := m.groups[monitorID]
	m.mutex.RUnlock()

	if !ok {
		return fmt.Errorf("invalid monitor ID: %s", monitorID)
	}

	m.mutex.Lock()
	m.setTimeoutOnGroup(group)
	m.mutex.Unlock()

	log.V("Monitor - subscriberenew: monitorID: %s", monitorID)
	return nil
}

// Subscribe requests that the monitor keep track of updates for all of the zones
// and sensors listed in the MonitorGroup parameter. If refresh == true, the monitor
// will go and request values for all items in the monitor group, if false it won't
// until the caller calls the Subscribe method.  The function returns a monitorID value
// that can be passed into other functions, such as Unsubscribe and Refresh.
func (m *Monitor) Subscribe(g *MonitorGroup, refresh bool) (string, error) {

	if len(g.Sensors) == 0 && len(g.Zones) == 0 {
		return "", errors.New("no zones or sensors listed in the monitor group")
	}

	m.mutex.Lock()
	monitorID := strconv.FormatInt(m.nextID, 10)
	m.nextID++
	g.id = monitorID
	m.groups[monitorID] = g

	// store the time that this will expire
	m.setTimeoutOnGroup(g)

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

	log.V("Monitor - subscribe: refresh: %t, monitorID: %s, %s", refresh, monitorID, g)

	if refresh {
		m.Refresh(monitorID, false)
	}

	return monitorID, nil
}

// Unsubscribe removes all references and updates for the specified monitorID
func (m *Monitor) Unsubscribe(monitorID string) {
	if _, ok := m.groups[monitorID]; !ok {
		return
	}

	emptySensorToGroupCount := 0
	emptyZoneToGroupCount := 0

	m.mutex.Lock()
	delete(m.groups, monitorID)
	for sensorID, groups := range m.sensorToGroups {
		if _, ok := groups[monitorID]; ok {
			delete(groups, monitorID)
			if len(groups) == 0 {
				emptySensorToGroupCount++
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
				emptyZoneToGroupCount++
				delete(m.zoneToGroups, zoneID)
				delete(m.zoneValues, zoneID)
			}
		}
	}
	m.mutex.Unlock()

	log.V("Monitor - unsubscribe: monitorID: %s, emptyZoneToGroups: %d, emptySensorToGroups:%d",
		monitorID, emptyZoneToGroupCount, emptySensorToGroupCount)
}

func (m *Monitor) sensorAttrChanged(sensorID string, attr SensorAttr) {
	m.mutex.RLock()
	groups, ok := m.sensorToGroups[sensorID]
	m.mutex.RUnlock()

	if !ok {
		// Not a sensor we are monitoring, ignore
		return
	}

	// Is this value different to what we already know?
	m.mutex.RLock()
	currentVal, ok := m.sensorValues[sensorID]
	m.mutex.RUnlock()
	if ok {
		// No change, don't refresh clients
		if currentVal.Value == attr.Value {
			return
		}
	}

	m.mutex.Lock()
	m.sensorValues[sensorID] = attr
	m.mutex.Unlock()

	for groupID := range groups {
		m.mutex.RLock()
		group := m.groups[groupID]
		cb := &ChangeBatch{
			MonitorID: groupID,
			Sensors:   make(map[string]SensorAttr),
		}
		cb.Sensors[sensorID] = attr
		m.mutex.RUnlock()
		group.Handler.Update(cb)
	}
}

func (m *Monitor) zoneLevelChanged(zoneID string, val cmd.Level) {
	m.mutex.RLock()
	groups, ok := m.zoneToGroups[zoneID]
	m.mutex.RUnlock()
	if !ok {
		return
	}

	// Is this value different to what we already know?
	m.mutex.RLock()
	currentVal, ok := m.zoneValues[zoneID]
	m.mutex.RUnlock()
	if ok {
		// No change, don't refresh clients
		if currentVal == val {
			return
		}
	}

	m.mutex.Lock()
	m.zoneValues[zoneID] = val
	m.mutex.Unlock()

	for groupID := range groups {
		m.mutex.RLock()
		group := m.groups[groupID]
		cb := &ChangeBatch{
			MonitorID: groupID,
			Zones:     make(map[string]cmd.Level),
		}
		cb.Zones[zoneID] = val
		m.mutex.RUnlock()
		group.Handler.Update(cb)
	}
}

// deviceProducing is called when a device start producing events, in the case of
// the monitor we need to see if there are MonitorGroups that require values from
// this device and then request the latest values
func (m *Monitor) deviceProducing(evt *DeviceProducingEvt) {
	groups := make(map[string]bool)

	m.mutex.RLock()
	for _, zone := range evt.Device.Zones {
		grp, ok := m.zoneToGroups[zone.ID]
		if ok {
			for monitorID := range grp {
				groups[monitorID] = true
			}
		}
	}
	for _, sensor := range evt.Device.Sensors {
		grp, ok := m.sensorToGroups[sensor.ID]
		if ok {
			for monitorID := range grp {
				groups[monitorID] = true
			}
		}
	}
	m.mutex.RUnlock()

	log.V("Monitor - %s, found %d group to refresh", evt, len(groups))

	for monitorID := range groups {
		m.Refresh(monitorID, false)
	}
}

// handleTimeouts watches for monitor groups that have expired and purges them
// from the system
func (m *Monitor) handleTimeouts() {
	go func() {
		for {
			now := time.Now()
			var expired []*MonitorGroup
			m.mutex.RLock()
			for _, group := range m.groups {
				if group.timeoutAbsolute.Before(now) {
					expired = append(expired, group)
				}
			}
			m.mutex.RUnlock()

			for _, group := range expired {
				log.V("Monitor - group expired, monitorID: %s", group.id)

				m.Unsubscribe(group.id)
				group.Handler.Expired(group.id)
			}

			// Sleep then wake up and check again for the next expired items
			time.Sleep(time.Second * 5)
		}
	}()
}

// setTimeoutOnGroup sets the time that the group will expire, once a group has
// expired we no longer keep clients updated about changes
func (m *Monitor) setTimeoutOnGroup(group *MonitorGroup) {
	group.timeoutAbsolute = time.Now().Add(group.Timeout)
}

// ======= evtbus.Consumer interface

func (m *Monitor) ConsumerName() string {
	return "Monitor"
}

func (m *Monitor) StartConsuming(c chan evtbus.Event) {
	log.V("Monitor - start consuming events")

	go func() {
		for e := range c {
			showHandled := true

			switch evt := e.(type) {
			case *SensorAttrChangedEvt:
				m.sensorAttrChanged(evt.SensorID, evt.Attr)

			case *SensorsReportingEvt:
				for sensorID, attr := range evt.Sensors {
					m.sensorAttrChanged(sensorID, attr)
				}

			case *ZonesReportingEvt:
				for zoneID, val := range evt.Zones {
					m.zoneLevelChanged(zoneID, val)
				}

			case *ZoneLevelChangedEvt:
				m.zoneLevelChanged(evt.ZoneID, evt.Level)

			case *DeviceProducingEvt:
				m.deviceProducing(evt)
				// Have dedicated logging for this, don't show generic message
				showHandled = false

			default:
				showHandled = false
			}

			if showHandled {
				log.V("Monitor - handled event: %s", e)
			}
		}

		log.V("Monitor - consumer event channel has closed")
	}()
}

func (m *Monitor) StopConsuming() {
	//TODO:
}

// =================================

// ======== evtbus.Producer interface
//TODO: Remove?
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
