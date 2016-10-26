package gohome

import "fmt"

//TODO: Move to event package

// SensorAttrChanged represents an event when the attributes of a sensor have
// changed state
type SensorAttrChanged struct {
	// SensorID is the ID of the sensor whos values have changed
	SensorID string

	// SensorName is the name of the sensor whos values have changed
	SensorName string

	// Information on the attribute that changed
	Attr SensorAttr
}

func (e *SensorAttrChanged) String() string {
	return fmt.Sprintf("SensorAttrChanged, Sensor Name:%s, ID:%s, %s", e.SensorName, e.SensorID, e.Attr.String())
}

//TODO: Different file
// SensorsReport is an event indicating that sensors included in the SensorIDs map
// should report their current attribute status on the event bus
type SensorsReport struct {
	SensorIDs map[string]bool
}

func (sr *SensorsReport) Add(sensorID string) {
	if sr.SensorIDs == nil {
		sr.SensorIDs = make(map[string]bool)
	}
	sr.SensorIDs[sensorID] = true
}

func (e *SensorsReport) String() string {
	return fmt.Sprintf("SensorsReport, contains %d sensors", len(e.SensorIDs))
}

type SensorsReporting struct {
	Sensors map[string]SensorAttr
}

func (sr *SensorsReporting) Add(sensorID string, attr SensorAttr) {
	if sr.Sensors == nil {
		sr.Sensors = make(map[string]SensorAttr)
	}
	sr.Sensors[sensorID] = attr
}

func (sr *SensorsReporting) String() string {
	return fmt.Sprintf("SensorsReporting, contains %d sensors", len(sr.Sensors))
}

//TODO: ZoneLevelChanged
//TODO: ButtonPressed
//TODO: ButtonReleased
//TODO: SceneSet
//TODO: Delete event/event.go
