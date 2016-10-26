package gohome

import (
	"strconv"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/log"
)

//TODO:Delete
/*
// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}
*/ //TODO: Put back item, item should point to monitorgroup??

type ChangeHandler interface {
	Update(b *ChangeBatch)
}

type ChangeBatch struct {
	Sensors map[string]SensorAttr
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
	groups         map[string]*MonitorGroup
	evtBus         *evtbus.Bus
	sensorToGroups map[string]map[string]bool
	sensorValues   map[string]SensorAttr
}

func NewMonitor(sys *System, evtBus *evtbus.Bus, sensorValues map[string]SensorAttr) *Monitor {

	//Callers can pass in initial values if they know what they are
	if sensorValues == nil {
		sensorValues = make(map[string]SensorAttr)
	}
	m := &Monitor{
		system:         sys,
		nextID:         1,
		groups:         make(map[string]*MonitorGroup),
		sensorToGroups: make(map[string]map[string]bool),
		sensorValues:   sensorValues,
		evtBus:         evtBus,
	}

	evtBus.AddConsumer(m)
	evtBus.AddProducer(m)
	return m
}

func (m *Monitor) Start() {
}

func (m *Monitor) Subscribe(g *MonitorGroup, refresh bool) string {

	//TODO: Called in multiple go routines?, mutex it

	monitorID := strconv.Itoa(m.nextID)
	m.nextID++
	m.groups[monitorID] = g

	//TODO: Zones

	var changeBatch = &ChangeBatch{
		Sensors: make(map[string]SensorAttr),
	}
	var report = &SensorsReport{}
	for sensorID := range g.Sensors {
		// Get the monitor groups that are listening to this sensor
		groups, ok := m.sensorToGroups[sensorID]
		if !ok {
			groups = make(map[string]bool)
			m.sensorToGroups[sensorID] = groups
		}

		groups[monitorID] = true

		// Need to subscribe to the changes if we haven't already ...
		// TODO: How to say subscribe to the sensor events?
		// pass in functions to subscribe to certain objects, so this knows nothing about the system ...

		// Caller wants to get values if we have them, or if not request
		if refresh {
			val, ok := m.sensorValues[sensorID]
			if ok {
				changeBatch.Sensors[sensorID] = val
			} else {
				report.Add(sensorID)
				//m.valueRequester.RequestSensor(sensorID, attr, func(attr SensorAttr) {
				//m.sensorAttrChanged(sensorID, attr)
				//})
			}
		}
	}
	if len(changeBatch.Sensors) > 0 {
		g.Handler.Update(changeBatch)
	}
	if len(report.SensorIDs) > 0 {
		m.evtBus.Enqueue(report)
	}
	// TODO: Where to set timeout, heap with timeout, set timer for smallest time
	// TODO: expirationHeap

	//TODO: Go through each of the inputs, if we have a value for them, send them out, then request updates
	//TODO: How to update, some things can push, others we will have to poll for

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
		group := m.groups[groupID]
		//TODO: Don't make blocking...
		//TODO: Batch?
		cb := &ChangeBatch{
			Sensors: make(map[string]SensorAttr),
		}
		cb.Sensors[sensorID] = attr
		group.Handler.Update(cb)
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
}

func (m *Monitor) StopProducing() {
	//TODO:
}

// ==================================

//TODO: Unsubscribe
