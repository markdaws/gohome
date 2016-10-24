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

//TODO: ZoneLevelChanged
//TODO: ButtonPressed
//TODO: ButtonReleased
//TODO: SceneSet
