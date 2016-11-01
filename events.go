package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
)

//TODO: ButtonPressed
//TODO: ButtonReleased
//TODO: SceneSet
//TODO: Delete event/event.go

// SensorAttrChangedEvt represents an event when the attributes of a sensor have
// changed state
type SensorAttrChangedEvt struct {
	// SensorID is the ID of the sensor whos values have changed
	SensorID string

	// SensorName is the name of the sensor whos values have changed
	SensorName string

	// Information on the attribute that changed
	Attr SensorAttr
}

func (e *SensorAttrChangedEvt) String() string {
	return fmt.Sprintf("SensorAttrChangedEvt[Name:%s, ID:%s, %s]", e.SensorName, e.SensorID, e.Attr.String())
}

// ZoneLevelChangedEvt represents an event when a zones level changes
type ZoneLevelChangedEvt struct {
	// ZoneID is the ID of the zone whos value has changed
	ZoneID string

	// ZoneName is the name of the zone whos value changed
	ZoneName string

	// Level contains the current zone level information
	Level cmd.Level
}

func (e *ZoneLevelChangedEvt) String() string {
	return fmt.Sprintf("ZoneLevelChangedEvt[Name:%s, ID:%s %s]", e.ZoneName, e.ZoneID, e.Level.String())
}

// SensorsReportEvt is an event indicating that sensors included in the SensorIDs map
// should report their current attribute status on the event bus
type SensorsReportEvt struct {
	SensorIDs map[string]bool
}

func (sr *SensorsReportEvt) Add(sensorID string) {
	if sr.SensorIDs == nil {
		sr.SensorIDs = make(map[string]bool)
	}
	sr.SensorIDs[sensorID] = true
}
func (e *SensorsReportEvt) String() string {
	return fmt.Sprintf("SensorsReportEvt[#sensors: %d]", len(e.SensorIDs))
}

// SensorsReportingEvt is an event that is fired when sensors are reporting changes in their
// attribute values
type SensorsReportingEvt struct {
	Sensors map[string]SensorAttr
}

func (sr *SensorsReportingEvt) Add(sensorID string, attr SensorAttr) {
	if sr.Sensors == nil {
		sr.Sensors = make(map[string]SensorAttr)
	}
	sr.Sensors[sensorID] = attr
}
func (sr *SensorsReportingEvt) String() string {
	return fmt.Sprintf("SensorsReportingEvt[#sensors: %d]", len(sr.Sensors))
}

// ZonesReportEvt is an event indicating that the specified zones should report
// their current value
type ZonesReportEvt struct {
	ZoneIDs map[string]bool
}

// Add adds a new zone to the report
func (zr *ZonesReportEvt) Add(zoneID string) {
	if zr.ZoneIDs == nil {
		zr.ZoneIDs = make(map[string]bool)
	}
	zr.ZoneIDs[zoneID] = true
}

// Merge merges the zones in the provided parameter into the target report
func (zr *ZonesReportEvt) Merge(rpt *ZonesReportEvt) {
	for zoneID := range rpt.ZoneIDs {
		zr.ZoneIDs[zoneID] = true
	}
}
func (zr *ZonesReportEvt) String() string {
	return fmt.Sprintf("ZonesReportEvt[#zones: %d]", len(zr.ZoneIDs))
}

// ZonesReportingEvt is an event that fires when zones are reporting changes to
// their current level
type ZonesReportingEvt struct {
	Zones map[string]cmd.Level
}

func (zr *ZonesReportingEvt) Add(zoneID string, level cmd.Level) {
	if zr.Zones == nil {
		zr.Zones = make(map[string]cmd.Level)
	}
	zr.Zones[zoneID] = level
}
func (zr *ZonesReportingEvt) String() string {
	return fmt.Sprintf("ZonesReportingEvt[#zones: %d]", len(zr.Zones))
}

// DeviceProducingEvt is raised when a device starts producing events in the system
type DeviceProducingEvt struct {
	Device *Device
}

func (dp *DeviceProducingEvt) String() string {
	return fmt.Sprintf("DeviceProducingEvt[%s]", dp.Device)
}
