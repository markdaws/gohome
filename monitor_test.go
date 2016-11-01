package gohome_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/zone"
	"github.com/stretchr/testify/require"
)

type MockChangeHandler struct {
	ChangeBatches []*gohome.ChangeBatch
	ExpiredIDs    []string
}

func (h *MockChangeHandler) Update(cb *gohome.ChangeBatch) {
	h.ChangeBatches = append(h.ChangeBatches, cb)
}
func (h *MockChangeHandler) Expired(monitorID string) {
	h.ExpiredIDs = append(h.ExpiredIDs, monitorID)
}
func (h *MockChangeHandler) Reset() {
	h.ChangeBatches = nil
	h.ExpiredIDs = []string{}
}

func makeSystemWithZonesAndSensors(nZones, nSensors int) (*gohome.System, []*zone.Zone, []*gohome.Sensor) {

	var zones []*zone.Zone
	var sensors []*gohome.Sensor
	sys := gohome.NewSystem("", "", 1)
	dev := &gohome.Device{
		Name:    "dev1",
		ID:      "1234",
		Zones:   make(map[string]*zone.Zone),
		Sensors: make(map[string]*gohome.Sensor),
	}
	sys.AddDevice(dev)

	for i := 0; i < nZones; i++ {
		var strI = strconv.Itoa(i)
		zone := &zone.Zone{
			ID:       strI,
			Name:     "test zone " + strI,
			Address:  strI,
			DeviceID: "1234",
		}
		dev.AddZone(zone)
		sys.AddZone(zone)
		zones = append(zones, zone)
	}

	for i := 0; i < nSensors; i++ {
		var strI = strconv.Itoa(i)
		sensor := &gohome.Sensor{
			Name:     "test sensor " + strI,
			Address:  strI,
			DeviceID: "1234",
			Attr: gohome.SensorAttr{
				Name:     "sensor" + strI,
				DataType: "int",
				Value:    "-1",
			},
		}
		sys.AddSensor(sensor)
		sensors = append(sensors, sensor)
	}
	return sys, zones, sensors
}

type EventConsumer struct {
	SensorsReport *gohome.SensorsReportEvt
	ZonesReport   *gohome.ZonesReportEvt
}

func (ec *EventConsumer) ConsumerName() string {
	return "EventConsumer"
}
func (ec *EventConsumer) StartConsuming(ch chan evtbus.Event) {
	go func() {
		for e := range ch {
			switch evt := e.(type) {
			case *gohome.SensorsReportEvt:
				ec.SensorsReport = evt
			case *gohome.ZonesReportEvt:
				ec.ZonesReport = evt
			}
		}
	}()
}
func (ec *EventConsumer) StopConsuming() {
}

// Test the Subscribe function.  Should make sure that the monitor returns and
// values it already knows about and requests values for ones it doesn't
func TestSubscribeSensors(t *testing.T) {

	// System contains sensors and zones
	sys, _, sensors := makeSystemWithZonesAndSensors(0, 4)
	sensor1 := sensors[0]
	sensor2 := sensors[1]
	sensor3 := sensors[2]
	sensor4 := sensors[3]

	evtBus := evtbus.NewBus(100, 100)
	evtConsumer := &EventConsumer{}
	evtBus.AddConsumer(evtConsumer)

	// Create a monitor, we will pass in some initial state to pretend we know
	// about the value of sensor2, sensor4, this should cause the monitor to not request
	// a value for it and also return the value it knows about to the monitor group
	initialSensorValues := make(map[string]gohome.SensorAttr)
	var attr2 = sensor2.Attr
	attr2.Value = "10"
	initialSensorValues[sensor2.ID] = attr2
	var attr4 = sensor4.Attr
	attr4.Value = "20"
	initialSensorValues[sensor4.ID] = attr4

	m := gohome.NewMonitor(sys, evtBus, initialSensorValues, nil)

	mockHandler := &MockChangeHandler{}

	// Request to monitor certain items
	group := &gohome.MonitorGroup{
		Sensors: make(map[string]bool),
		Handler: mockHandler,
		Timeout: time.Duration(5) * time.Second,
	}

	// Add a sensor to the group, so we monitor it
	group.Sensors[sensor1.ID] = true
	group.Sensors[sensor2.ID] = true
	group.Sensors[sensor3.ID] = true
	group.Sensors[sensor4.ID] = true

	// Begin the subscription, should get back a monitor ID
	mID, _ := m.Subscribe(group, true)
	require.NotEqual(t, "", mID)

	// Processing is async, small delay to let event bus process
	time.Sleep(time.Millisecond * 1000)

	// Should have got an event asking for certain sensors to report their status
	// our sensor should be included in that
	require.NotNil(t, evtConsumer.SensorsReport)
	require.True(t, evtConsumer.SensorsReport.SensorIDs[sensor1.ID])
	require.True(t, evtConsumer.SensorsReport.SensorIDs[sensor3.ID])
	require.False(t, evtConsumer.SensorsReport.SensorIDs[sensor2.ID])
	require.False(t, evtConsumer.SensorsReport.SensorIDs[sensor4.ID])

	// For sensors 2 and 4 we should have got an update callback since we passed in their
	// values to the monitor when we inited it
	require.Equal(t, attr2, mockHandler.ChangeBatches[0].Sensors[sensor2.ID])
	require.Equal(t, attr4, mockHandler.ChangeBatches[0].Sensors[sensor4.ID])

	// Now respond to the request for sensors 1 and 3 to report their values
	reporting := &gohome.SensorsReportingEvt{}
	sensor1Attr := gohome.SensorAttr{
		Name:  "sensor1",
		Value: "111",
	}
	reporting.Add(sensor1.ID, sensor1Attr)
	sensor3Attr := gohome.SensorAttr{
		Name:  "sensor3",
		Value: "333",
	}
	reporting.Add(sensor3.ID, sensor3Attr)

	// Processing is async, small delay to let event bus process
	mockHandler.ChangeBatches = nil
	evtBus.Enqueue(reporting)

	time.Sleep(time.Millisecond * 1000)

	// We should have got updates with the attribute values we are expecting
	require.Equal(t, 2, len(mockHandler.ChangeBatches))
	require.Equal(t, sensor1Attr, mockHandler.ChangeBatches[0].Sensors[sensor1.ID])
	require.Equal(t, sensor3Attr, mockHandler.ChangeBatches[1].Sensors[sensor3.ID])
}

func TestMultipleGroupsOnTheSameSensorAreUpdated(t *testing.T) {
	// If we hav multiple monitor groups looking at the same sensor, then we need to
	// make sure when the sensor updates all of the groups receive notification of the
	// change

	// System contains sensors and zones
	sys, _, sensors := makeSystemWithZonesAndSensors(0, 4)
	sensor1 := sensors[0]
	sensor2 := sensors[1]
	sensor3 := sensors[2]
	sensor4 := sensors[3]

	evtBus := evtbus.NewBus(100, 100)
	evtConsumer := &EventConsumer{}
	evtBus.AddConsumer(evtConsumer)

	m := gohome.NewMonitor(sys, evtBus, nil, nil)

	mockHandler1 := &MockChangeHandler{}
	mockHandler2 := &MockChangeHandler{}

	group1 := &gohome.MonitorGroup{
		Sensors: make(map[string]bool),
		Handler: mockHandler1,
		Timeout: time.Duration(300) * time.Second,
	}
	group1.Sensors[sensor1.ID] = true
	group1.Sensors[sensor2.ID] = true

	group2 := &gohome.MonitorGroup{
		Sensors: make(map[string]bool),
		Handler: mockHandler2,
		Timeout: time.Duration(300) * time.Second,
	}
	group2.Sensors[sensor2.ID] = true
	group2.Sensors[sensor3.ID] = true
	group2.Sensors[sensor4.ID] = true

	mID1, _ := m.Subscribe(group1, false)
	require.NotEqual(t, "", mID1)

	mID2, _ := m.Subscribe(group2, false)
	require.NotEqual(t, "", mID2)

	// Sensor1 update should only update handler1
	attr1 := sensor1.Attr
	attr1.Value = "10"
	evtBus.Enqueue(&gohome.SensorAttrChangedEvt{
		SensorID: sensor1.ID,
		Attr:     attr1,
	})

	time.Sleep(time.Millisecond * 1000)
	require.Equal(t, 1, len(mockHandler1.ChangeBatches))
	require.Equal(t, 1, len(mockHandler1.ChangeBatches[0].Sensors))
	require.Equal(t, attr1, mockHandler1.ChangeBatches[0].Sensors[sensor1.ID])
	require.Equal(t, 0, len(mockHandler2.ChangeBatches))

	// Sensor3 update should only update handler2
	mockHandler1.ChangeBatches = nil
	attr3 := sensor3.Attr
	attr3.Value = "30"
	evtBus.Enqueue(&gohome.SensorAttrChangedEvt{
		SensorID: sensor3.ID,
		Attr:     attr3,
	})

	time.Sleep(time.Millisecond * 1000)
	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
	require.Equal(t, 1, len(mockHandler2.ChangeBatches[0].Sensors))
	require.Equal(t, attr3, mockHandler2.ChangeBatches[0].Sensors[sensor3.ID])
	require.Equal(t, 0, len(mockHandler1.ChangeBatches))

	// Sensor2 update should update handler1 and handler2 since they both subscribe to it
	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil
	attr2 := sensor2.Attr
	attr2.Value = "20"
	evtBus.Enqueue(&gohome.SensorAttrChangedEvt{
		SensorID: sensor2.ID,
		Attr:     attr2,
	})
	time.Sleep(time.Millisecond * 1000)
	require.Equal(t, 1, len(mockHandler1.ChangeBatches))
	require.Equal(t, 1, len(mockHandler1.ChangeBatches[0].Sensors))
	require.Equal(t, attr2, mockHandler1.ChangeBatches[0].Sensors[sensor2.ID])
	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
	require.Equal(t, 1, len(mockHandler2.ChangeBatches[0].Sensors))
	require.Equal(t, attr2, mockHandler2.ChangeBatches[0].Sensors[sensor2.ID])
}

func TestSubscribeZones(t *testing.T) {

	// System contains sensors and zones
	sys, zones, _ := makeSystemWithZonesAndSensors(4, 0)
	zone1 := zones[0]
	zone2 := zones[1]
	zone3 := zones[2]
	zone4 := zones[3]

	evtBus := evtbus.NewBus(100, 100)
	evtConsumer := &EventConsumer{}
	evtBus.AddConsumer(evtConsumer)

	// Create a monitor, we will pass in some initial state to pretend we know
	// about the value of zone2, zone4, this should cause the monitor to not request
	// a value for it and also return the value it knows about to the monitor group
	initialZoneValues := make(map[string]cmd.Level)
	var lvl2 = cmd.Level{}
	lvl2.Value = 10
	initialZoneValues[zone2.ID] = lvl2
	var lvl4 = cmd.Level{}
	lvl4.Value = 20
	initialZoneValues[zone4.ID] = lvl4

	m := gohome.NewMonitor(sys, evtBus, nil, initialZoneValues)

	mockHandler := &MockChangeHandler{}

	// Request to monitor certain items
	group := &gohome.MonitorGroup{
		Sensors: make(map[string]bool),
		Zones:   make(map[string]bool),
		Handler: mockHandler,
		Timeout: time.Duration(300) * time.Second,
	}

	// Add a zone to the group, so we monitor it
	group.Zones[zone1.ID] = true
	group.Zones[zone2.ID] = true
	group.Zones[zone3.ID] = true
	group.Zones[zone4.ID] = true

	// Begin the subscription, should get back a monitor ID
	mID, _ := m.Subscribe(group, true)
	require.NotEqual(t, "", mID)

	// Processing is async, small delay to let event bus process
	time.Sleep(time.Millisecond * 1000)

	// Should have got an event asking for certain zones to report their status
	// our zone should be included in that
	require.NotNil(t, evtConsumer.ZonesReport)
	require.True(t, evtConsumer.ZonesReport.ZoneIDs[zone1.ID])
	require.True(t, evtConsumer.ZonesReport.ZoneIDs[zone3.ID])
	require.False(t, evtConsumer.ZonesReport.ZoneIDs[zone2.ID])
	require.False(t, evtConsumer.ZonesReport.ZoneIDs[zone4.ID])

	// For zones 2 and 4 we should have got an update callback since we passed in their
	// values to the monitor when we inited it
	require.Equal(t, lvl2, mockHandler.ChangeBatches[0].Zones[zone2.ID])
	require.Equal(t, lvl4, mockHandler.ChangeBatches[0].Zones[zone4.ID])

	// Now respond to the request for zones 1 and 3 to report their values
	reporting := &gohome.ZonesReportingEvt{}
	zone1Lvl := cmd.Level{
		Value: 11,
	}
	reporting.Add(zone1.ID, zone1Lvl)
	zone3Lvl := cmd.Level{
		Value: 22,
	}
	reporting.Add(zone3.ID, zone3Lvl)

	// Processing is async, small delay to let event bus process
	mockHandler.ChangeBatches = nil
	evtBus.Enqueue(reporting)

	time.Sleep(time.Millisecond * 1000)

	// We should have got updates with the attribute values we are expecting
	require.Equal(t, 2, len(mockHandler.ChangeBatches))
	require.Equal(t, zone1Lvl, mockHandler.ChangeBatches[0].Zones[zone1.ID])
	require.Equal(t, zone3Lvl, mockHandler.ChangeBatches[1].Zones[zone3.ID])
}

func TestMultipleGroupsOnTheSameZoneAreUpdated(t *testing.T) {
	// If we hav multiple monitor groups looking at the same sensor, then we need to
	// make sure when the sensor updates all of the groups receive notification of the
	// change

	// System contains zones and zones
	sys, zones, _ := makeSystemWithZonesAndSensors(4, 0)
	zone1 := zones[0]
	zone2 := zones[1]
	zone3 := zones[2]
	zone4 := zones[3]

	evtBus := evtbus.NewBus(100, 100)
	evtConsumer := &EventConsumer{}
	evtBus.AddConsumer(evtConsumer)

	m := gohome.NewMonitor(sys, evtBus, nil, nil)

	mockHandler1 := &MockChangeHandler{}
	mockHandler2 := &MockChangeHandler{}

	group1 := &gohome.MonitorGroup{
		Zones:   make(map[string]bool),
		Handler: mockHandler1,
		Timeout: time.Duration(300) * time.Second,
	}
	group1.Zones[zone1.ID] = true
	group1.Zones[zone2.ID] = true

	group2 := &gohome.MonitorGroup{
		Zones:   make(map[string]bool),
		Handler: mockHandler2,
		Timeout: time.Duration(300) * time.Second,
	}
	group2.Zones[zone2.ID] = true
	group2.Zones[zone3.ID] = true
	group2.Zones[zone4.ID] = true

	mID1, _ := m.Subscribe(group1, false)
	require.NotEqual(t, "", mID1)

	mID2, _ := m.Subscribe(group2, false)
	require.NotEqual(t, "", mID2)

	// Zone1 update should only update handler1
	lvl1 := cmd.Level{Value: 10}
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone1.ID,
		Level:  lvl1,
	})

	time.Sleep(time.Millisecond * 1000)
	require.Equal(t, 1, len(mockHandler1.ChangeBatches))
	require.Equal(t, 1, len(mockHandler1.ChangeBatches[0].Zones))
	require.Equal(t, lvl1, mockHandler1.ChangeBatches[0].Zones[zone1.ID])
	require.Equal(t, 0, len(mockHandler2.ChangeBatches))

	// Zone3 update should only update handler2
	mockHandler1.ChangeBatches = nil
	lvl3 := cmd.Level{Value: 30}
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone3.ID,
		Level:  lvl3,
	})

	time.Sleep(time.Millisecond * 1000)
	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
	require.Equal(t, 1, len(mockHandler2.ChangeBatches[0].Zones))
	require.Equal(t, lvl3, mockHandler2.ChangeBatches[0].Zones[zone3.ID])
	require.Equal(t, 0, len(mockHandler1.ChangeBatches))

	// Zone2 update should update handler1 and handler2 since they both subscribe to it
	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil

	lvl2 := cmd.Level{Value: 20}
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone2.ID,
		Level:  lvl2,
	})

	time.Sleep(time.Millisecond * 1000)
	require.Equal(t, 1, len(mockHandler1.ChangeBatches))
	require.Equal(t, 1, len(mockHandler1.ChangeBatches[0].Zones))
	require.Equal(t, lvl2, mockHandler1.ChangeBatches[0].Zones[zone2.ID])
	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
	require.Equal(t, 1, len(mockHandler2.ChangeBatches[0].Zones))
	require.Equal(t, lvl2, mockHandler2.ChangeBatches[0].Zones[zone2.ID])
}

func TestUnsubscribe(t *testing.T) {
	// System contains sensors and zones
	sys, zones, sensors := makeSystemWithZonesAndSensors(3, 3)
	zone1 := zones[0]
	zone2 := zones[1]
	zone3 := zones[2]
	sensor1 := sensors[0]
	sensor2 := sensors[1]
	sensor3 := sensors[2]

	evtBus := evtbus.NewBus(100, 100)
	m := gohome.NewMonitor(sys, evtBus, nil, nil)

	mockHandler1 := &MockChangeHandler{}
	mockHandler2 := &MockChangeHandler{}

	// Request to monitor certain items
	group1 := &gohome.MonitorGroup{
		Sensors: make(map[string]bool),
		Zones:   make(map[string]bool),
		Handler: mockHandler1,
		Timeout: time.Duration(300) * time.Second,
	}
	group2 := &gohome.MonitorGroup{
		Sensors: make(map[string]bool),
		Zones:   make(map[string]bool),
		Handler: mockHandler2,
		Timeout: time.Duration(300) * time.Second,
	}

	// Got two monitor groups, both contain sensor2 and zone2
	group1.Zones[zone1.ID] = true
	group1.Zones[zone2.ID] = true
	group1.Sensors[sensor1.ID] = true
	group1.Sensors[sensor2.ID] = true

	group2.Zones[zone2.ID] = true
	group2.Zones[zone3.ID] = true
	group2.Sensors[sensor2.ID] = true
	group2.Sensors[sensor3.ID] = true

	// Begin the subscription, should get back a monitor ID
	mID1, _ := m.Subscribe(group1, true)
	require.NotEqual(t, "", mID1)
	mID2, _ := m.Subscribe(group2, true)
	require.NotEqual(t, "", mID2)

	// Processing is async, small delay to let event bus process
	time.Sleep(time.Millisecond * 1000)

	// Clear out any previous change notifications
	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil

	// zone1 change should only affect handler1
	lvl1 := cmd.Level{Value: 10}
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone1.ID,
		Level:  lvl1,
	})
	time.Sleep(1000 * time.Millisecond)

	require.Equal(t, 1, len(mockHandler1.ChangeBatches))
	require.Equal(t, 0, len(mockHandler2.ChangeBatches))

	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil

	// zone3 change should only affect handler2
	lvl3 := cmd.Level{Value: 30}
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone3.ID,
		Level:  lvl3,
	})
	time.Sleep(1000 * time.Millisecond)

	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
	require.Equal(t, 0, len(mockHandler1.ChangeBatches))

	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil

	// zone2 change should affect handler1 and handler2
	lvl2 := cmd.Level{Value: 20}
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone2.ID,
		Level:  lvl2,
	})
	time.Sleep(1000 * time.Millisecond)

	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
	require.Equal(t, 1, len(mockHandler1.ChangeBatches))

	// Zone1 update should not come back to handler1 any more, also zone2 should not
	// come back to handler1
	m.Unsubscribe(mID1)
	m.InvalidateValues(mID1)
	m.InvalidateValues(mID2)
	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone1.ID,
		Level:  lvl1,
	})
	time.Sleep(1000 * time.Millisecond)
	require.Equal(t, 0, len(mockHandler1.ChangeBatches))
	require.Equal(t, 0, len(mockHandler2.ChangeBatches))

	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone2.ID,
		Level:  lvl2,
	})
	time.Sleep(1000 * time.Millisecond)
	require.Equal(t, 0, len(mockHandler1.ChangeBatches))
	require.Equal(t, 1, len(mockHandler2.ChangeBatches))

	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone3.ID,
		Level:  lvl3,
	})
	time.Sleep(1000 * time.Millisecond)
	require.Equal(t, 0, len(mockHandler1.ChangeBatches))
	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
}

func TestGroupExpires(t *testing.T) {
	sys, zones, sensors := makeSystemWithZonesAndSensors(1, 1)
	zone1 := zones[0]
	sensor1 := sensors[0]

	evtBus := evtbus.NewBus(100, 100)
	m := gohome.NewMonitor(sys, evtBus, nil, nil)

	mockHandler1 := &MockChangeHandler{}
	group1 := &gohome.MonitorGroup{
		Sensors: make(map[string]bool),
		Zones:   make(map[string]bool),
		Handler: mockHandler1,
		Timeout: time.Duration(3) * time.Second,
	}

	group1.Zones[zone1.ID] = true
	group1.Sensors[sensor1.ID] = true

	mID1, _ := m.Subscribe(group1, true)
	require.NotEqual(t, "", mID1)
	require.Equal(t, 0, len(mockHandler1.ExpiredIDs))

	// Group expires in 3 seconds, monitor checks every 5 so wait until after
	time.Sleep(time.Second * 6)

	require.Equal(t, mID1, mockHandler1.ExpiredIDs[0])

	//Expired group should not receive any updates
	mockHandler1.ChangeBatches = nil
	lvl1 := cmd.Level{Value: 10}
	evtBus.Enqueue(&gohome.ZoneLevelChangedEvt{
		ZoneID: zone1.ID,
		Level:  lvl1,
	})
	attr1 := sensor1.Attr
	attr1.Value = "10"
	evtBus.Enqueue(&gohome.SensorAttrChangedEvt{
		SensorID: sensor1.ID,
		Attr:     attr1,
	})

	time.Sleep(1000 * time.Millisecond)
	require.Equal(t, 0, len(mockHandler1.ChangeBatches))
}
