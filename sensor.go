package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/validation"
)

//TODO: Monitoring - sometimes only want to monitor when the UI is activly viewing
//the sensor, but other times always want to monitor e.g. fridge temperature, flags?

type SensorDataType string

const (
	SDTInt    SensorDataType = "int"
	SDTFloat  SensorDataType = "float"
	SDTBool   SensorDataType = "bool"
	SDTString SensorDataType = "string"
)

type SensorAttr struct {
	Name          string
	Value         string
	DataType      SensorDataType
	UnitOfMeasure string
	States        map[string]string
}

func (a SensorAttr) String() string {
	return fmt.Sprintf("SensorAttr - Name:%s, Value:%s, DataType:%s", a.Name, a.Value, a.DataType)
}

type Sensor struct {
	ID          string
	Name        string
	Description string
	Address     string
	DeviceID    string
	Attr        SensorAttr
}

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
