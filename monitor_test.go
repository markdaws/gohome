package gohome_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/stretchr/testify/require"
)

type MockChangeHandler struct {
	ChangeBatches []*gohome.ChangeBatch
}

func (h *MockChangeHandler) Update(cb *gohome.ChangeBatch) {
	h.ChangeBatches = append(h.ChangeBatches, cb)

	fmt.Printf("got update callback %+v\n", cb)
}

func makeSystem() (*gohome.System, *gohome.Sensor, *gohome.Sensor, *gohome.Sensor, *gohome.Sensor) {
	sys := gohome.NewSystem("", "", nil, 1)
	sensor1 := &gohome.Sensor{
		Name:     "test sensor 1",
		Address:  "1",
		DeviceID: "1234",
		Attr: gohome.SensorAttr{
			Name:     "sensor1",
			DataType: "int",
			Value:    "-1",
		},
	}
	sensor2 := &gohome.Sensor{
		Name:     "test sensor 2",
		Address:  "2",
		DeviceID: "12345",
		Attr: gohome.SensorAttr{
			Name:     "sensor2",
			DataType: "int",
			Value:    "-1",
		},
	}
	sensor3 := &gohome.Sensor{
		Name:     "test sensor 3",
		Address:  "3",
		DeviceID: "123456",
		Attr: gohome.SensorAttr{
			Name:     "sensor3",
			DataType: "int",
			Value:    "-1",
		},
	}
	sensor4 := &gohome.Sensor{
		Name:     "test sensor 4",
		Address:  "4",
		DeviceID: "1234567",
		Attr: gohome.SensorAttr{
			Name:     "sensor4",
			DataType: "int",
			Value:    "-1",
		},
	}

	sys.AddSensor(sensor1)
	sys.AddSensor(sensor2)
	sys.AddSensor(sensor3)
	sys.AddSensor(sensor4)

	return sys, sensor1, sensor2, sensor3, sensor4
}

type EventConsumer struct {
	SensorsReport *gohome.SensorsReport
}

func (ec *EventConsumer) ConsumerName() string {
	return "EventConsumer"
}
func (ec *EventConsumer) StartConsuming(ch chan evtbus.Event) {
	go func() {
		for e := range ch {
			switch evt := e.(type) {
			case *gohome.SensorsReport:
				ec.SensorsReport = evt
			}
		}
	}()
}
func (ec *EventConsumer) StopConsuming() {
}

// Test the Subscribe function.  Should make sure that the monitor returns and
// values it already knows about and requests values for ones it doesn't
func TestSubscribe(t *testing.T) {

	// System contains sensors and zones
	sys, sensor1, sensor2, sensor3, sensor4 := makeSystem()

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

	m := gohome.NewMonitor(sys, evtBus, initialSensorValues)

	mockHandler := &MockChangeHandler{}

	// Request to monitor certain items
	group := &gohome.MonitorGroup{
		Sensors:          make(map[string]bool),
		Handler:          mockHandler,
		TimeoutInSeconds: 300,
	}

	// Add a sensor to the group, so we monitor it
	group.Sensors[sensor1.ID] = true
	group.Sensors[sensor2.ID] = true
	group.Sensors[sensor3.ID] = true
	group.Sensors[sensor4.ID] = true

	// Begin the subscription, should get back a monitor ID
	mID := m.Subscribe(group, true)
	require.NotEqual(t, "", mID)

	// Processing is async, small delay to let event bus process
	time.Sleep(time.Millisecond * 100)

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
	reporting := &gohome.SensorsReporting{}
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

	time.Sleep(time.Millisecond * 100)

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
	sys, sensor1, sensor2, sensor3, sensor4 := makeSystem()

	evtBus := evtbus.NewBus(100, 100)
	evtConsumer := &EventConsumer{}
	evtBus.AddConsumer(evtConsumer)

	m := gohome.NewMonitor(sys, evtBus, nil)

	mockHandler1 := &MockChangeHandler{}
	mockHandler2 := &MockChangeHandler{}

	group1 := &gohome.MonitorGroup{
		Sensors:          make(map[string]bool),
		Handler:          mockHandler1,
		TimeoutInSeconds: 300,
	}
	group1.Sensors[sensor1.ID] = true
	group1.Sensors[sensor2.ID] = true

	group2 := &gohome.MonitorGroup{
		Sensors:          make(map[string]bool),
		Handler:          mockHandler2,
		TimeoutInSeconds: 300,
	}
	group2.Sensors[sensor2.ID] = true
	group2.Sensors[sensor3.ID] = true
	group2.Sensors[sensor4.ID] = true

	mID1 := m.Subscribe(group1, false)
	require.NotEqual(t, "", mID1)

	mID2 := m.Subscribe(group2, false)
	require.NotEqual(t, "", mID2)

	// Sensor1 update should only update handler1
	attr1 := sensor1.Attr
	attr1.Value = "10"
	evtBus.Enqueue(&gohome.SensorAttrChanged{
		SensorID: sensor1.ID,
		Attr:     attr1,
	})

	time.Sleep(time.Millisecond * 100)
	require.Equal(t, 1, len(mockHandler1.ChangeBatches))
	require.Equal(t, 1, len(mockHandler1.ChangeBatches[0].Sensors))
	require.Equal(t, attr1, mockHandler1.ChangeBatches[0].Sensors[sensor1.ID])
	require.Equal(t, 0, len(mockHandler2.ChangeBatches))

	// Sensor3 update should only update handler2
	mockHandler1.ChangeBatches = nil
	attr3 := sensor3.Attr
	attr3.Value = "30"
	evtBus.Enqueue(&gohome.SensorAttrChanged{
		SensorID: sensor3.ID,
		Attr:     attr3,
	})

	time.Sleep(time.Millisecond * 100)
	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
	require.Equal(t, 1, len(mockHandler2.ChangeBatches[0].Sensors))
	require.Equal(t, attr3, mockHandler2.ChangeBatches[0].Sensors[sensor3.ID])
	require.Equal(t, 0, len(mockHandler1.ChangeBatches))

	// Sensor2 update should update handler1 and handler2 since they both subscribe to it
	mockHandler1.ChangeBatches = nil
	mockHandler2.ChangeBatches = nil
	attr2 := sensor2.Attr
	attr2.Value = "20"
	evtBus.Enqueue(&gohome.SensorAttrChanged{
		SensorID: sensor2.ID,
		Attr:     attr2,
	})
	time.Sleep(time.Millisecond * 100)
	require.Equal(t, 1, len(mockHandler1.ChangeBatches))
	require.Equal(t, 1, len(mockHandler1.ChangeBatches[0].Sensors))
	require.Equal(t, attr2, mockHandler1.ChangeBatches[0].Sensors[sensor2.ID])
	require.Equal(t, 1, len(mockHandler2.ChangeBatches))
	require.Equal(t, 1, len(mockHandler2.ChangeBatches[0].Sensors))
	require.Equal(t, attr2, mockHandler2.ChangeBatches[0].Sensors[sensor2.ID])
}

//TODO: TestExpire
