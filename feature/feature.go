package feature

import (
	"fmt"

	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/validation"
)

const (
	// FTButton button, can be phycial or virtual
	FTButton string = "Button"

	// FTCoolZone cooling zone
	FTCoolZone string = "CoolZone"

	// FTHeatZone heating zone
	FTHeatZone string = "HeatZone"

	// FTLightZone lighting zone
	FTLightZone string = "LightZone"

	// FTOutlet outlet
	FTOutlet string = "Outlet"

	// FTSensor sensor
	FTSensor string = "Sensor"

	// FTSwitch switch
	FTSwitch string = "Switch"

	// FTWindowTreatment window treatment such as a shade or curtain
	FTWindowTreatment string = "WindowTreatment"
)

// Attrs a map of attributes, keyed by the attributes local ID
type Attrs map[string]*attr.Attribute

// NewAttrs returns a map of attributes keyed by the LocalID field
func NewAttrs(attrs ...*attr.Attribute) Attrs {
	out := Attrs{}
	for _, attribute := range attrs {
		if attribute == nil {
			continue
		}
		out[attribute.LocalID] = attribute
	}
	return out
}

// Feature represents a unit of functionality that a device may provide, such as a
// button, sensor, switch etc.  Devices can potentially export many features.
type Feature struct {
	// ID a globally unique ID for the feature
	ID string `json:"id"`

	// Type represents a concrete type for the feature e.g. LightZone, Outlet etc.
	Type string `json:"type"`

	// Address an optional address to identify the feature locally
	Address string `json:"address"`

	// Name a user friendly name that may be displayed in the UI
	Name string `json:"name"`

	// Description more details about the feature which may be shown in the UI
	Description string `json:"description"`

	// DeviceID is the ID of the device that owns the feature
	DeviceID string `json:"deviceId"`

	// Attrs is a map of attributes, keyed by the attributes LocalID field
	Attrs Attrs `json:"attrs"`

	//TODO: Remove - don't use, only for API serialization purposes
	IsDupe bool `json:"isDupe"`
}

// String returns a debug string for the feature
func (f *Feature) String() string {
	return fmt.Sprintf("Feature[ID:%s, Type:%s, Address:%s, Name:%s, DeviceID:%s]",
		f.ID, f.Type, f.Address, f.Name, f.DeviceID)
}

// Validate returns nil if the feature is in a valid state, otherwise a validation
// error is returned detailing the issues
func (f *Feature) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if f.ID == "" {
		errors.Add("required field", "ID")
	}

	if f.Name == "" {
		errors.Add("required field", "Name")
	}

	if f.DeviceID == "" {
		errors.Add("required field", "DeviceID")
	}

	if f.Type == "" {
		errors.Add("require field", "Type")
	}

	if errors.Has() {
		return errors
	}
	return nil
}

// NewFromType returns a new feature instance based on the speicifed feature type
// passed in to the function
func NewFromType(ID, fType string) *Feature {
	switch fType {
	case FTHeatZone:
		return NewHeatZone(ID)
	case FTLightZone:
		return NewLightZone(ID, LightZoneModeHSL)
	case FTOutlet:
		return NewOutlet(ID)
	case FTSwitch:
		return NewSwitch(ID)
	case FTWindowTreatment:
		return NewWindowTreatment(ID)
	case FTButton:
		//TODO:
		return nil
	case FTCoolZone:
		//TODO:
		return nil
	case FTSensor:
		// Don't support sensor since a sensor needs an attribute
		// as well to be initialized
		return nil
	default:
		return nil
	}
}

// NewSensor returns a feature instance initialized as a Sensor
func NewSensor(ID string, attribute *attr.Attribute) *Feature {
	s := &Feature{
		ID:    ID,
		Type:  FTSensor,
		Attrs: Attrs{attribute.LocalID: attribute},
	}
	return s
}

const (
	// HeatZoneCurrentTempLocalID is the localID value for the current temp attribute
	HeatZoneCurrentTempLocalID string = "currenttemp"

	// HeatZoneTargetTempLocalID is the localID value for the target temp attribute
	HeatZoneTargetTempLocalID string = "targettemp"
)

// NewHeatZone returns a featuer instance initialized as a heat zone. Heat Zones represents
// zones on a thermostat that can provide heat
func NewHeatZone(ID string) *Feature {
	s := &Feature{
		ID:   ID,
		Type: FTHeatZone,
	}

	// The current temp is read only since you set the target temp not the current
	current := attr.NewTemp("currenttemp", nil)
	current.Perms = attr.PermsReadOnly
	current.Name = "Current Temperature"

	target := attr.NewTemp("targettemp", nil)
	target.Name = "Target Temperature"
	s.Attrs = Attrs{
		current.LocalID: current,
		target.LocalID:  target,
	}
	return s
}

// HeatZoneCloneAttrs clone the common attributes for a heat zone so they can be updated
func HeatZoneCloneAttrs(f *Feature) (current, target *attr.Attribute) {
	var ok bool
	if current, ok = f.Attrs[HeatZoneCurrentTempLocalID]; ok {
		current = current.Clone()
	}

	if target, ok = f.Attrs[HeatZoneTargetTempLocalID]; ok {
		target = target.Clone()
	}

	return
}

const (
	// LightZoneOnOffLocalID is the local ID for the onoff attribute
	LightZoneOnOffLocalID string = "onoff"

	// LightZoneBrightnessLocalID is the local ID for the brightness attribute
	LightZoneBrightnessLocalID string = "brightness"

	// LightZoneHSLLocalID is the local ID for the HSL attribute
	LightZoneHSLLocalID string = "hsl"

	// LightZoneModeBinary indicates the light can only be in either an on or off
	// state, nothing inbetween i.e. not dimmable
	LightZoneModeBinary = "binary"

	// LightZoneModeContinuous indicates the light can be dimmed
	LightZoneModeContinuous = "continuous"

	// LightZoneModeHSL indicates the light supports different colours mapped to
	// values in the HSL color space
	LightZoneModeHSL = "hsl"
)

// NewLightZone returns a feature initialized as a LightZone.  A Light zone represents
// a single or multiple bulbs
func NewLightZone(ID, mode string) *Feature {
	s := &Feature{
		ID:    ID,
		Type:  FTLightZone,
		Attrs: Attrs{},
	}

	// All light zones can be turned on and off
	onOff := attr.NewOnOff("onoff", nil)
	onOff.Name = "On/Off"
	s.Attrs[onOff.LocalID] = onOff

	switch mode {
	case LightZoneModeBinary:
		// Nothing else to do, only support on/off

	case LightZoneModeContinuous:
		// Light can be dimmed
		brightness := attr.NewBrightness("brightness", nil)
		brightness.Name = "Brightness"
		s.Attrs[brightness.LocalID] = brightness

	case LightZoneModeHSL:
		hsl := attr.NewHSL("hsl", nil)
		hsl.Name = "HSL"
		s.Attrs[hsl.LocalID] = hsl
	}

	return s
}

// LightZoneCloneAttrs clone the common attributes for a light zone so they can be updated
func LightZoneCloneAttrs(f *Feature) (onoff, brightness, hsl *attr.Attribute) {
	var ok bool
	if brightness, ok = f.Attrs[LightZoneBrightnessLocalID]; ok {
		brightness = brightness.Clone()
	}

	if onoff, ok = f.Attrs[LightZoneOnOffLocalID]; ok {
		onoff = onoff.Clone()
	}

	if hsl, ok = f.Attrs[LightZoneHSLLocalID]; ok {
		hsl = hsl.Clone()
	}
	return
}

// LightZoneGetBrightness returns the brightness the light should be set to. It takes in to account
// if you have an onoff and brightness attribute or just a brightness attribute
func LightZoneGetBrightness(attrs map[string]*attr.Attribute) (float32, error) {
	onoff := attrs[LightZoneOnOffLocalID]
	brightness := attrs[LightZoneBrightnessLocalID]

	var level float32 = -1
	if onoff != nil {
		if onoff.Value.(int32) == attr.OnOffOff {
			level = 0
		} else {
			if brightness != nil {
				level = brightness.Value.(float32)
			} else {
				level = 100
			}
		}
	} else {
		if brightness != nil {
			level = brightness.Value.(float32)
		}
	}

	if level == -1 {
		return 0, fmt.Errorf("both onoff and brightness attributes are missing, require at least one")
	}
	return level, nil
}

const (
	// WindowTreatmentOffsetLocalID is the local ID of the offset attribute
	WindowTreatmentOffsetLocalID string = "offset"

	// WindowTreatmentOpenCloseLocalID is the local ID of the openclose attribute
	WindowTreatmentOpenCloseLocalID string = "openclose"
)

// NewWindowTreatment returns a new feature initialized as a WindowTreatment.  A window treatment can
// be a shade or curtain or anything covering a window that can be monitored and controlled and has
// an offset position.  An offset of 100% means the window treamtent is fully open, 0% means
// fully closed
func NewWindowTreatment(ID string) *Feature {
	s := &Feature{
		ID:   ID,
		Type: FTWindowTreatment,
	}
	offset := attr.NewOffset("offset", nil)
	offset.Name = "Offset"
	openClose := attr.NewOpenClose("openclose", nil)
	openClose.Name = "Open/Close"
	s.Attrs = Attrs{
		offset.LocalID:    offset,
		openClose.LocalID: openClose,
	}
	return s
}

// WindowTreatmentGetOffset returns the desired offset, taking into account the open/close value
// if it is present in the attributes. So if openclose is closed even if there is a offset
// value it is ignored.
func WindowTreatmentGetOffset(attrs map[string]*attr.Attribute) (float32, error) {
	openclose := attrs[WindowTreatmentOpenCloseLocalID]
	offset := attrs[WindowTreatmentOffsetLocalID]

	var level float32 = -1
	if openclose != nil {
		if openclose.Value.(int32) == attr.OpenCloseClosed {
			level = 0
		} else {
			if offset != nil {
				level = offset.Value.(float32)
			} else {
				level = 100
			}
		}
	} else {
		if offset != nil {
			level = offset.Value.(float32)
		}
	}

	if level == -1 {
		return 0, fmt.Errorf("both openclose and offset attributes are missing, require at least one")
	}
	return level, nil
}

// WindowTreatmentCloneAttrs clone the common attributes for a window treatment so they can be updated
func WindowTreatmentCloneAttrs(f *Feature) (openClose, offset *attr.Attribute) {
	var ok bool
	if openClose, ok = f.Attrs[WindowTreatmentOpenCloseLocalID]; ok {
		openClose = openClose.Clone()
	}

	if offset, ok = f.Attrs[WindowTreatmentOffsetLocalID]; ok {
		offset = offset.Clone()
	}
	return
}

// NewButton returns a new button instance
// TODO: Attributes?
func NewButton(ID string) *Feature {
	b := &Feature{
		ID:    ID,
		Type:  FTButton,
		Attrs: make(map[string]*attr.Attribute),
	}
	return b
}

const (
	// SwitchOnOffLocalID is the local ID of the onoff attribute
	SwitchOnOffLocalID string = "onoff"
)

// NewSwitch returns a new feature initialized as a Switch feature. Switches can be turned on and off
func NewSwitch(ID string) *Feature {
	s := &Feature{
		ID:   ID,
		Type: FTSwitch,
	}
	onOff := attr.NewOnOff("onoff", nil)
	onOff.Name = "On/Off"
	s.Attrs = Attrs{onOff.LocalID: onOff}
	return s
}

// SwitchCloneAttrs clone the common attributes for a switch so they can be updated
func SwitchCloneAttrs(f *Feature) (onOff *attr.Attribute) {
	var ok bool
	if onOff, ok = f.Attrs[SwitchOnOffLocalID]; ok {
		onOff = onOff.Clone()
	}
	return
}

const (
	// OutletOnOffLocalID is the local ID of the onoff attribute
	OutletOnOffLocalID string = "onoff"
)

// NewOutlet returns a new feature initialized as an outlet.  An outlet is a plug socket that
// can be turned on or off
func NewOutlet(ID string) *Feature {
	s := &Feature{
		ID:   ID,
		Type: FTOutlet,
	}
	onOff := attr.NewOnOff("onoff", nil)
	onOff.Name = "On/Off"
	s.Attrs = Attrs{onOff.LocalID: onOff}
	return s
}

// OutletCloneAttrs clone the common attributes for an Outlet so they can be updated
func OutletCloneAttrs(f *Feature) (onOff *attr.Attribute) {
	var ok bool
	if onOff, ok = f.Attrs[OutletOnOffLocalID]; ok {
		onOff = onOff.Clone()
	}
	return
}
