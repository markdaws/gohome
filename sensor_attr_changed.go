package gohome

// SensorAttrChanged represents an event when the attributes of a sensor have
// changed state
type SensorAttrChanged struct {
	// SensorID is the ID of the sensor whos values have changed
	SensorID string

	// SensorName is the name of the sensor whos values have changed
	SensorName string

	// Attrs if a slice of attributes that have changed
	Attrs []SensorAttr
}

//TODO: ZoneSetLevel
//TODO: ButtonPress
//TODO: ButtonRelease
//TODO: SceneSet
