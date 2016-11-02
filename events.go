package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
)

//TODO: SceneSet

// SensorAttrChangedEvt represents an event when the attributes of a sensor have
// changed state
type SensorAttrChangedEvt struct {
	// SensorID is the ID of the sensor whos values have changed
	SensorID string `json:"id"`

	// SensorName is the name of the sensor whos values have changed
	SensorName string `json:"name"`

	// Information on the attribute that changed
	Attr SensorAttr `json:"attr"`
}

func (e *SensorAttrChangedEvt) String() string {
	return fmt.Sprintf("SensorAttrChangedEvt[Name:%s, ID:%s, %s]", e.SensorName, e.SensorID, e.Attr.String())
}

// SensorAttrReportingEvt represents an event when the attributes of a sensor are
// reported back to the system. This doesn't mean they have changed, just that
// the system asked for the latest value and this is it
type SensorAttrReportingEvt struct {
	// SensorID is the ID of the sensor whos values have changed
	SensorID string `json:"id"`

	// SensorName is the name of the sensor whos values have changed
	SensorName string `json:"name"`

	// Information on the attribute that changed
	Attr SensorAttr `json:"attr"`
}

func (e *SensorAttrReportingEvt) String() string {
	return fmt.Sprintf("SensorAttrReportingEvt[Name:%s, ID:%s, %s]", e.SensorName, e.SensorID, e.Attr.String())
}

// ZoneLevelChangedEvt represents an event when a zones level changes
type ZoneLevelChangedEvt struct {
	// ZoneID is the ID of the zone whos value has changed
	ZoneID string `json:"id"`

	// ZoneName is the name of the zone whos value changed
	ZoneName string `json:"name"`

	// Level contains the current zone level information
	Level cmd.Level `json:"level"`
}

func (e *ZoneLevelChangedEvt) String() string {
	return fmt.Sprintf("ZoneLevelChangedEvt[Name:%s, ID:%s %s]", e.ZoneName, e.ZoneID, e.Level.String())
}

// ZoneLevelReportingEvt represents an event when a zones level was queried and the value is being reported
type ZoneLevelReportingEvt struct {
	// ZoneID is the ID of the zone whos value has changed
	ZoneID string `json:"id"`

	// ZoneName is the name of the zone whos value changed
	ZoneName string `json:"name"`

	// Level contains the current zone level information
	Level cmd.Level `json:"level"`
}

func (e *ZoneLevelReportingEvt) String() string {
	return fmt.Sprintf("ZoneLevelReportingEvt[Name:%s, ID:%s %s]", e.ZoneName, e.ZoneID, e.Level.String())
}

// ButtonPressEvt is raised when a button is pressed in the system
type ButtonPressEvt struct {
	// BtnAddress is the address of the button
	BtnAddress string `json:"address"`

	// BtnID is the global ID of the button
	BtnID string `json:"id"`

	// BtnName is the name of the button
	BtnName string `json:"name"`

	// DeviceID is the ID of the device to which the button belongs
	DeviceID string `json:"deviceId"`

	// DeviceName is the name of the device
	DeviceName string `json:"deviceName"`

	// DeviceAddress is the address of the device
	DeviceAddress string `json:"deviceAddress"`
}

func (e *ButtonPressEvt) String() string {
	return fmt.Sprintf("ButtonPressEvt[Addr:%s, ID: %s, Name:%s, DevAddr:%s, DevID:%s, DevName:%s",
		e.BtnAddress, e.BtnID, e.BtnName, e.DeviceAddress, e.DeviceID, e.DeviceName)
}

// ButtonPressEvt is raised when a button is released in the system
type ButtonReleaseEvt struct {
	// BtnAddress is the address of the button
	BtnAddress string `json:"address"`

	// BtnID is the global ID of the button
	BtnID string `json:"id"`

	// BtnName is the name of the button
	BtnName string `json:"name"`

	// DeviceID is the ID of the device to which the button belongs
	DeviceID string `json:"deviceId"`

	// DeviceName is the name of the device
	DeviceName string `json:"deviceName"`

	// DeviceAddress is the address of the device
	DeviceAddress string `json:"deviceAddress"`
}

func (e *ButtonReleaseEvt) String() string {
	return fmt.Sprintf("ButtonReleaseEvt[Addr:%s, ID: %s, Name:%s, DevAddr:%s, DevID:%s, DevName:%s",
		e.BtnAddress, e.BtnID, e.BtnName, e.DeviceAddress, e.DeviceID, e.DeviceName)
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

// DeviceLostEvt indicates that connection to a device has been lost
type DeviceLostEvt struct {
	DeviceName string
	DeviceID   string
}

func (dl *DeviceLostEvt) String() string {
	return fmt.Sprintf("DeviceLostEvt[ID: %s, Name: %s]", dl.DeviceName, dl.DeviceID)
}

// ClientConnectedEvt is raised when a client registers to get updates for zone and sensor values
type ClientConnectedEvt struct {
	ConnectionID string `json:"connectionId"`
	MonitorID    string `json:"-"`
	Origin       string `json:"origin"`
}

func (cc *ClientConnectedEvt) String() string {
	return fmt.Sprintf("ClientConnectedEvt[ConnID: %s, MonitorID: %s, Origin:%s]",
		cc.ConnectionID, cc.MonitorID, cc.Origin)
}

// ClientDisconnectedEvt is raised when a client connection is closed
type ClientDisconnectedEvt struct {
	ConnectionID string `json:"connectionId"`
}

func (cc *ClientDisconnectedEvt) String() string {
	return fmt.Sprintf("ClientDisconnectedEvt[ConnID: %s", cc.ConnectionID)
}

// UserLoginEvt is fired when a user logs in to the system, or there is an unsuccessful
// login attempt
type UserLoginEvt struct {
	Login   string `json:"login"`
	Success bool   `json:"success"`
}

func (ul *UserLoginEvt) String() string {
	return fmt.Sprintf("UserLoginEvt[Login: %s, Success: %t]", ul.Login, ul.Success)
}

// UserLogoutEvt is fired when a user logs out of the system explicitly
type UserLogoutEvt struct {
	Login string `json:"login"`
}

func (ul *UserLogoutEvt) String() string {
	return fmt.Sprintf("UserLogoutEvt[Login: %s]", ul.Login)
}

//TODO: Finish device lost plumbing
//TODO: DeviceConnectedEvt
