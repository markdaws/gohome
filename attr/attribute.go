package attr

import (
	"fmt"

	"github.com/markdaws/gohome/log"
)

const (
	// DTString string
	DTString string = "string"

	// DTBool boolean
	DTBool string = "bool"

	// DTFloat32 float32
	DTFloat32 string = "float32"

	// DTInt32 int32
	DTInt32 string = "int32"
)

const (
	// UTPercentage percentage
	UTPercentage string = "percentage"

	// UTCelcius celcius
	UTCelcius string = "celcius"

	// UTFarenheit farenheit
	UTFarenheit string = "farenheit"

	// UTMillisecond millisecond
	UTMilliSecond string = "millisecond"
)

const (
	// ATOpenClose represents an OpenClose attribute
	ATOpenClose string = "OpenClose"

	// ATBrightness represents a Brightness attribute
	ATBrightness string = "Brightness"

	// ATOnOff represents an OnOff attribute
	ATOnOff string = "OnOff"

	// ATHue represents a Hue attribute
	ATHue string = "Hue"

	// ATOffset represents an Offset attribute
	ATOffset string = "Offset"

	// ATTemperature represents a temperature attribute
	ATTemperature string = "Temperature"
)

const (
	// PermsReadOnly the attribute is a read-only value and cannot be set
	PermsReadOnly string = "r"

	// PermsReadWrite the attribute can be read and written to
	PermsReadWrite string = "rw"
)

// Attribute represents a value with a specified data type and range.
type Attribute struct {
	// LocalID is an identifier that should be mutually exclusive between all attributes
	// assigned to a feature, it doesn't need to be globally unique. It can be used to
	// distinguish between the various atrtibutes of a feature
	LocalID string `json:"localId"`

	// Type is the concrete type of attribute such as OpenClose, Brightness, Hue etc
	Type string `json:"type"`

	// DataType is the underlying data type used in the attribute, such as int32, bool etc
	DataType string `json:"dataType"`

	// Unit is the data units used by the attribute e.g. percentage
	Unit string `json:"unit"`

	// Name is a user friendly string which can be shown in the UI
	Name string `json:"name"`

	// Description provides more details about the attribute that can be shown in the UI
	Description string `json:"description"`

	// Value is the value of the attribute
	Value interface{} `json:"value"`

	// Min is the minimum allowed value
	Min interface{} `json:"min"`

	// Max is the max allowed value
	Max interface{} `json:"max"`

	// Step is the step size
	Step interface{} `json:"step"`

	// Perms specifies if the user has read or readwrite permissions, either 'r' or 'rw'
	Perms string `json:"perms"`
}

func (a Attribute) String() string {
	return fmt.Sprintf("Attribute[LocalID: %s, Type: %s, Value: %+v, Perms: %s]",
		a.LocalID, a.Type, a.Value, a.Perms,
	)
}

// FixJSON massages the values back from float64 which is the type given to the values
// when being unmarshalled, to their correct data type
func FixJSON(attrs map[string]*Attribute) {
	for _, attribute := range attrs {
		// When these are deserialized from JSON the interface values get the wrong type, need to
		// massage them back to the expected type
		if attribute.Value != nil {
			switch attribute.DataType {
			case DTFloat32:
				attribute.Value = float32(attribute.Value.(float64))
			case DTInt32:
				attribute.Value = int32(attribute.Value.(float64))
			}
		}
		if attribute.Min != nil {
			switch attribute.DataType {
			case DTFloat32:
				attribute.Min = float32(attribute.Min.(float64))
			case DTInt32:
				attribute.Min = int32(attribute.Min.(float64))
			}
		}
		if attribute.Max != nil {
			switch attribute.DataType {
			case DTFloat32:
				attribute.Max = float32(attribute.Max.(float64))
			case DTInt32:
				attribute.Max = int32(attribute.Max.(float64))
			}
		}
		if attribute.Step != nil {
			switch attribute.DataType {
			case DTFloat32:
				attribute.Step = float32(attribute.Step.(float64))
			case DTInt32:
				attribute.Step = int32(attribute.Step.(float64))
			}
		}
	}
}

// Clone returns a cloned copy of the attribute
func (a *Attribute) Clone() *Attribute {
	b := *a
	b.Value = a.cloneInterface(b.Value)
	b.Min = a.cloneInterface(b.Min)
	b.Max = a.cloneInterface(b.Max)
	b.Step = a.cloneInterface(b.Step)
	return &b
}

func (a *Attribute) cloneInterface(val interface{}) interface{} {
	if val == nil {
		return nil
	}

	switch val.(type) {
	case string:
		s := val.(string)
		return s
	case bool:
		b := val.(bool)
		return b
	case int32:
		i := val.(int32)
		return i
	case float32:
		f := val.(float32)
		return f
	case float64:
		panic("got a float64")
	default:
		log.E("unknown data type in attribute: %s", a.DataType)
		return nil
	}
}

// NewBool returns a new attribute initialized as a boolean data type
func NewBool(localID, concrete string, val *bool) *Attribute {
	attr := NewAttribute(localID, concrete, DTBool, val)
	return attr
}

// NewInt32 returns a new attribute initialized as an int32 data type
func NewInt32(localID, concrete string, val *int32) *Attribute {
	attr := NewAttribute(localID, concrete, DTInt32, val)
	return attr
}

// NewFloat32 returns a new attribute initialized as a float32 data type
func NewFloat32(localID, concrete string, val *float32) *Attribute {
	attr := NewAttribute(localID, concrete, DTFloat32, val)
	return attr
}

// BoolP returns a pointer to a bool containing the passed in value
func BoolP(val bool) *bool {
	return &val
}

// Int32P returns a pointer to an int32 containing the passed in value
func Int32P(val int32) *int32 {
	return &val
}

// Float32P returns a pointer to a float32 containing the passed in value
func Float32P(val float32) *float32 {
	return &val
}

// Clone copies the map and the contained attributes, returning a new copy
// that can be modified without chaning the original map or attributes
func CloneAttrs(attrs map[string]*Attribute) map[string]*Attribute {
	newAttrs := make(map[string]*Attribute)
	for k, v := range attrs {
		newAttrs[k] = v.Clone()
	}
	return newAttrs
}

// Only can be used if you only have one item in a map and want to access it
func Only(attrs map[string]*Attribute) *Attribute {
	for _, v := range attrs {
		return v
	}
	return nil
}

// NewAttribute creates and returns a new Attribute instance
func NewAttribute(localID, t, dataType string, val interface{}) *Attribute {
	a := &Attribute{
		LocalID:  localID,
		Type:     t,
		Perms:    PermsReadWrite,
		DataType: dataType,
	}

	switch dataType {
	case DTString:
		if v := val.(*string); v != nil {
			a.Value = *v
		}
	case DTBool:
		if v := val.(*bool); v != nil {
			a.Value = *v
		}
	case DTInt32:
		if v := val.(*int32); v != nil {
			a.Value = *v
		}
	case DTFloat32:
		if v := val.(*float32); v != nil {
			a.Value = *v
		}
	default:
		log.E("unknown data type in attribute: %s:%s:%t", dataType, DTFloat32, dataType == DTFloat32)
	}
	return a
}

const (
	// OpenCloseClosed indicates the OpenClose attribute is in a closed state
	OpenCloseClosed int32 = 1

	// OpenCloseOpen indicates the OpenClose attribute is in an open state
	OpenCloseOpen int32 = 2
)

// NweOpenClose returns a new attribute instance, initialized as an OpenClose type
func NewOpenClose(localID string, val *int32) *Attribute {
	attr := NewInt32(localID, ATOpenClose, val)
	return attr
}

const (
	// OnOffOff indicates the OnOff attribute is in the off state
	OnOffOff int32 = 1

	// OnOffOn indicates the OnOff attribute is in the on state
	OnOffOn int32 = 2
)

// NewOnOff returns a new Attribute instance initialized as an OnOff type
func NewOnOff(localID string, val *int32) *Attribute {
	attr := NewInt32(localID, ATOnOff, val)
	return attr
}

// NewBrightness returns a new Attribute initialized as a Brightness type
func NewBrightness(localID string, val *float32) *Attribute {
	attr := NewFloat32(localID, ATBrightness, val)
	attr.Unit = UTPercentage
	attr.Min = float32(0)
	attr.Max = float32(100)
	attr.Step = float32(1)
	return attr
}

// NewOffset returns a new Attribute instance initialized as an Offset type
func NewOffset(localID string, val *float32) *Attribute {
	attr := NewFloat32(localID, ATOffset, val)
	attr.Unit = UTPercentage
	attr.Min = float32(0)
	attr.Max = float32(100)
	attr.Step = float32(1)
	return attr
}

// NewHue returns a new Attribute instance initialized as a Hue type
func NewHue(localID string, val *int32) *Attribute {
	attr := NewInt32(localID, ATHue, val)
	attr.Min = int32(0)
	attr.Max = int32(359)
	attr.Step = int32(1)
	return attr
}

// NewTemp returns a new Attribute instance initialized as a Temp type
func NewTemp(localID string, val *int32) *Attribute {
	attr := NewInt32(localID, ATTemperature, val)
	attr.Unit = UTFarenheit
	attr.Min = int32(40)
	attr.Max = int32(80)
	attr.Step = int32(1)
	return attr
}

/*
const (
	HeatingCoolingModeOff  int = 0
	HeatingCoolingModeOn   int = 1
	HeatingCoolingModeAuto     = 2
)

type HeatingCoolingMode struct {
	*Int32
}

func NewHeatingCoolingMode() *HeatingCoolingMode {
	attr := NewInt32("HeatingCoolingMode", 2)
	return &HeatingCoolingMode{attr}
}
*/
