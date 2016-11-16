package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/validation"
)

// SensorDataType represents the type of data the sensor is returning
type SensorDataType string

const (
	// SDTInt - int data type
	SDTInt SensorDataType = "int"

	// SDTFloat - float data type
	SDTFloat SensorDataType = "float"

	// SDTBool - bool data type
	SDTBool SensorDataType = "bool"

	// SDTString - string data type
	SDTString SensorDataType = "string"
)

// SensorAttr is the attribute the sensor is monitoring. It may be temperature, an open/close
// state etc. It can be anything.  Each sensor can monitor one attribute. If you have a piece
// of hardware that monitors many values, when that would be modelled as the hardware having
// multiple sensors
type SensorAttr struct {
	// Name the name of the attribute e.g. temperature
	Name string

	// Value the current value of the attribute
	Value string

	// DataType the type of data the sensor is returning
	DataType SensorDataType

	UnitOfMeasure string

	// States is a map of values -> Name that can be displayed in the UI. For example,
	// the sensor may have Value==1 which means "Open" and 0 which means "Closed" you can
	// add those mappings to the States map so the UI can show a user friendly string instead
	// of the raw sensor values
	States map[string]string
}

// String returns a debug string containing information about the sensor
func (a SensorAttr) String() string {
	return fmt.Sprintf("SensorAttr - Name:%s, Value:%s, DataType:%s", a.Name, a.Value, a.DataType)
}

// Sensor represents a physical sensor that a piece of hardware may contains, such as a temperature
// sensor, an open/close sensor, wet/dry sensor etc.
type Sensor struct {
	ID          string
	Name        string
	Description string
	Address     string
	DeviceID    string
	Attr        SensorAttr
}

// Validate verifies if the sensor is in a valid state
func (s *Sensor) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if s.ID == "" {
		errors.Add("required field", "ID")
	}

	if s.Name == "" {
		errors.Add("required field", "Name")
	}

	if s.Address == "" {
		errors.Add("required field", "Address")
	}

	if s.DeviceID == "" {
		errors.Add("required field", "DeviceID")
	}

	if errors.Has() {
		return errors
	}
	return nil
}
