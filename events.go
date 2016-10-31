package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
)

//TODO: ButtonPressed
//TODO: ButtonReleased
//TODO: SceneSet
//TODO: Delete event/event.go

// SensorAttrChanged represents an event when the attributes of a sensor have
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
	return fmt.Sprintf("SensorAttrChanged, Name:%s, ID:%s, %s", e.SensorName, e.SensorID, e.Attr.String())
}

// ZoneLevelChanged represents an event when a zones level changes
type ZoneLevelChangedEvt struct {
	// ZoneID is the ID of the zone whos value has changed
	ZoneID string

	// ZoneName is the name of the zone whos value changed
	ZoneName string

	// Level contains the current zone level information
	Level cmd.Level
}

func (e *ZoneLevelChangedEvt) String() string {
	return fmt.Sprintf("ZoneLevelChanged, Name:%s, ID:%s %s", e.ZoneName, e.ZoneID, e.Level.String())
}

// SensorsReport is an event indicating that sensors included in the SensorIDs map
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
	return fmt.Sprintf("SensorsReport, contains %d sensors", len(e.SensorIDs))
}

// SensorsReporting is an event that is fired when sensors are reporting changes in their
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
	return fmt.Sprintf("SensorsReporting, contains %d sensors", len(sr.Sensors))
}

// ZonesReport is an event indicating that the specified zones should report
// their current value
type ZonesReportEvt struct {
	ZoneIDs map[string]bool
}

func (zr *ZonesReportEvt) Add(zoneID string) {
	if zr.ZoneIDs == nil {
		zr.ZoneIDs = make(map[string]bool)
	}
	zr.ZoneIDs[zoneID] = true
}
func (zr *ZonesReportEvt) String() string {
	return fmt.Sprintf("ZonesReport, contains %d zones", len(zr.ZoneIDs))
}

// ZonesReporting is an event that fires when zones are reporting changes to
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
	return fmt.Sprintf("ZonesReporting, contains %d zones", len(zr.Zones))
}
