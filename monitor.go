package gohome

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/log"
)

// MonitorDelegate is the interface for receiving updates values from the monitor
type MonitorDelegate interface {
	Update(b *ChangeBatch)
	Expired(monitorID string)
}

// MonitorGroup represents a group of features a client wished to receive updates for.
type MonitorGroup struct {
	//TODO: change name to FeatureIDs
	Features        map[string]bool
	Handler         MonitorDelegate
	Timeout         time.Duration
	timeoutAbsolute time.Time
	id              string
}

func (mg *MonitorGroup) String() string {
	return fmt.Sprintf("MonitorGroup[ID:%s, %d features]", mg.id, len(mg.Features))
}

// ChangeBatch contains a list of features whos values have changed
type ChangeBatch struct {
	MonitorID string
	Features  map[string]map[string]*attr.Attribute
}

func (cb *ChangeBatch) String() string {
	return fmt.Sprintf("ChangeBatch[monitorID: %s, #features:%d]", cb.MonitorID, len(cb.Features))
}

// Monitor keeps track of the current feature attribute values in the system and reports
// updates to clients
type Monitor struct {
	groups          map[string]*MonitorGroup
	system          *System
	nextID          int64
	evtBus          *evtbus.Bus
	featureToGroups map[string]map[string]bool
	featureValues   map[string]map[string]*attr.Attribute
	mutex           sync.RWMutex
}

// NewMonitor returns an initialzed Monitor instance
func NewMonitor(sys *System, evtBus *evtbus.Bus) *Monitor {

	m := &Monitor{
		system:          sys,
		nextID:          time.Now().UnixNano(),
		groups:          make(map[string]*MonitorGroup),
		featureToGroups: make(map[string]map[string]bool),
		featureValues:   make(map[string]map[string]*attr.Attribute),
		evtBus:          evtBus,
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
		Features:  make(map[string]map[string]*attr.Attribute),
	}

	// Build a list of features that need to report their values. If we
	// already have a value for a sensor we can just return that
	var featuresReport = &FeaturesReportEvt{}
	for featureID := range group.Features {
		val, ok := m.featureValues[featureID]
		if ok && !force {
			changeBatch.Features[featureID] = val
		} else {
			featuresReport.Add(featureID)
		}
	}

	log.V("Monitor - refreshing: %s, force:%t", group, force)
	log.V("Monitor - refreshing: cached values: [%s], uncached features: %s", changeBatch, featuresReport)

	m.mutex.RUnlock()

	if len(changeBatch.Features) > 0 {
		// We have some values already cached for certain items, return
		group.Handler.Update(changeBatch)
	}
	if len(featuresReport.FeatureIDs) > 0 {
		// Send event to request features report their current values
		m.evtBus.Enqueue(featuresReport)
	}
}

// InvalidateValues removes any cached values, for any features listed
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
	for featureID := range group.Features {
		delete(m.featureValues, featureID)
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

// Subscribe requests that the monitor keep track of updates for all of the features
// listed in the MonitorGroup parameter. If refresh == true, the monitor
// will go and request values for all items in the monitor group, if false it won't
// until the caller calls the Subscribe method.  The function returns a monitorID value
// that can be passed into other functions, such as Unsubscribe and Refresh.
func (m *Monitor) Subscribe(g *MonitorGroup, refresh bool) (string, error) {

	if len(g.Features) == 0 {
		return "", errors.New("no features listed in the monitor group")
	}

	m.mutex.Lock()
	monitorID := strconv.FormatInt(m.nextID, 10)
	m.nextID++
	g.id = monitorID
	m.groups[monitorID] = g

	// store the time that this will expire
	m.setTimeoutOnGroup(g)

	// Make sure we map from the feature ids back to this new group,
	// so that if any features change in the future we know that we
	// need to alert this group
	for featureID := range g.Features {
		// Get the monitor groups that are listening to this feature
		groups, ok := m.featureToGroups[featureID]
		if !ok {
			groups = make(map[string]bool)
			m.featureToGroups[featureID] = groups
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

	emptyFeatureToGroupCount := 0

	m.mutex.Lock()
	delete(m.groups, monitorID)
	for featureID, groups := range m.featureToGroups {
		if _, ok := groups[monitorID]; ok {
			delete(groups, monitorID)
			if len(groups) == 0 {
				emptyFeatureToGroupCount++
				delete(m.featureToGroups, featureID)
				delete(m.featureValues, featureID)
			}
		}
	}
	m.mutex.Unlock()

	log.V("Monitor - unsubscribe: monitorID: %s, emptyFeatureToGroups: %d",
		monitorID, emptyFeatureToGroupCount)
}

func (m *Monitor) featureReporting(featureID string, attrs map[string]*attr.Attribute) {
	m.mutex.RLock()
	groups, ok := m.featureToGroups[featureID]
	m.mutex.RUnlock()

	if !ok {
		// Not a feature we are monitoring, ignore
		return
	}

	// If not a valid featureID in the system, ignore
	_, ok = m.system.Features[featureID]
	if !ok {
		return
	}

	// Is this value different to what we already know, features can have multiple attributes, so we
	// need to check each one and see if it is different, if any are different then we need to report
	// otherwise we can short circuit
	m.mutex.RLock()
	currentAttrs, ok := m.featureValues[featureID]
	m.mutex.RUnlock()

	// Already have some values, check to see if there are any new ones
	if ok {
		for localID, attr := range attrs {
			currentVal, ok := currentAttrs[localID]
			if !ok {
				// We dont have this attribute value, can't short circuit
				break
			}

			if currentVal != attr.Value {
				// Value is different to what we have, can't short circuit
				break
			}
		}
	} else {
		currentAttrs = make(map[string]*attr.Attribute)
	}

	m.mutex.Lock()
	// Merge new attribute values with the ones we already know about
	for localID, attr := range attrs {
		currentAttrs[localID] = attr
	}
	m.featureValues[featureID] = currentAttrs
	m.mutex.Unlock()

	for groupID := range groups {
		m.mutex.RLock()
		group := m.groups[groupID]
		cb := &ChangeBatch{
			MonitorID: groupID,
			Features:  make(map[string]map[string]*attr.Attribute),
		}
		cb.Features[featureID] = currentAttrs
		m.mutex.RUnlock()
		group.Handler.Update(cb)
	}

	m.system.Services.EvtBus.Enqueue(&FeatureAttrsChangedEvt{
		FeatureID: featureID,
		Attrs:     currentAttrs,
	})

}

// deviceProducing is called when a device start producing events, in the case of
// the monitor we need to see if there are MonitorGroups that require values from
// this device and then request the latest values
func (m *Monitor) deviceProducing(evt *DeviceProducingEvt) {
	groups := make(map[string]bool)

	m.mutex.RLock()
	for _, feature := range evt.Device.Features {
		grp, ok := m.featureToGroups[feature.ID]
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
			switch evt := e.(type) {
			case *FeatureReportingEvt:
				m.featureReporting(evt.FeatureID, evt.Attrs)

			case *DeviceProducingEvt:
				m.deviceProducing(evt)
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
