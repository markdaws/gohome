package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/attr"
)

//TODO: SceneSet

// AutomationTriggeredEvt is fired when a piece of automation is triggered
type AutomationTriggeredEvt struct {
	Name string
}

// String returns a debug string
func (e *AutomationTriggeredEvt) String() string {
	return fmt.Sprintf("AutomationTriggeredEvt[Name: %s]", e.Name)
}

// SunriseEvt is fired when it is sunrise
type SunriseEvt struct{}

// String returns a debug string
func (e *SunriseEvt) String() string {
	return "SunriseEvt"
}

// Sunset is fired when it is sunset
type SunsetEvt struct{}

// String returns a debug string
func (e *SunsetEvt) String() string {
	return "SunsetEvt"
}

// FeaturesReportEvt is fired when the system wants certain features to
// report their current values
type FeaturesReportEvt struct {
	FeatureIDs map[string]bool
}

// Add adds a feature to the request list
func (e *FeaturesReportEvt) Add(featureID string) {
	if e.FeatureIDs == nil {
		e.FeatureIDs = make(map[string]bool)
	}
	e.FeatureIDs[featureID] = true
}

// String returns a debug string
func (e *FeaturesReportEvt) String() string {
	return fmt.Sprintf("FeaturesReportEvt[#features: %d]", len(e.FeatureIDs))
}

// FeatureReportingEvt is fired when a devices feature(s) responds to a FeaturesReportEvt
// for the current state of the device features
type FeatureReportingEvt struct {
	FeatureID string
	Attrs     map[string]*attr.Attribute
}

// String returns a debug string
func (e *FeatureReportingEvt) String() string {
	return fmt.Sprintf("FeatureReportingEvt[ID: %s]", e.FeatureID)
}

// FeatureAttrsChangedEvt is fired when the system determines an attribute of a feature has changed
// value.  Note this is not always reflective of the actual hardware in that when the system first
// starts and has no initial value, this will fire even if the value didn't change on the hardware intially
type FeatureAttrsChangedEvt struct {
	FeatureID string
	Attrs     map[string]*attr.Attribute
}

// String returns a debug string
func (e *FeatureAttrsChangedEvt) String() string {
	return fmt.Sprintf("FeatureAttrsChangedEvt[ID:%s, Attrs:%s]", e.FeatureID, e.Attrs)
}

// DeviceProducingEvt is raised when a device starts producing events in the system
type DeviceProducingEvt struct {
	Device *Device
}

// String returns a debug string
func (dp *DeviceProducingEvt) String() string {
	return fmt.Sprintf("DeviceProducingEvt[%s]", dp.Device)
}

// DeviceLostEvt indicates that connection to a device has been lost
type DeviceLostEvt struct {
	DeviceName string
	DeviceID   string
}

// String returns a debug string
func (dl *DeviceLostEvt) String() string {
	return fmt.Sprintf("DeviceLostEvt[ID: %s, Name: %s]", dl.DeviceName, dl.DeviceID)
}

// ClientConnectedEvt is raised when a client registers to get updates for zone and sensor values
type ClientConnectedEvt struct {
	ConnectionID string `json:"connectionId"`
	MonitorID    string `json:"-"`
	Origin       string `json:"origin"`
}

// String returns a debug string
func (cc *ClientConnectedEvt) String() string {
	return fmt.Sprintf("ClientConnectedEvt[ConnID: %s, MonitorID: %s, Origin:%s]",
		cc.ConnectionID, cc.MonitorID, cc.Origin)
}

// ClientDisconnectedEvt is raised when a client connection is closed
type ClientDisconnectedEvt struct {
	ConnectionID string `json:"connectionId"`
}

// String returns a debug string
func (cc *ClientDisconnectedEvt) String() string {
	return fmt.Sprintf("ClientDisconnectedEvt[ConnID: %s", cc.ConnectionID)
}

// UserLoginEvt is fired when a user logs in to the system, or there is an unsuccessful
// login attempt
type UserLoginEvt struct {
	Login   string `json:"login"`
	Success bool   `json:"success"`
}

// String returns a debug string
func (ul *UserLoginEvt) String() string {
	return fmt.Sprintf("UserLoginEvt[Login: %s, Success: %t]", ul.Login, ul.Success)
}

// UserLogoutEvt is fired when a user logs out of the system explicitly
type UserLogoutEvt struct {
	Login string `json:"login"`
}

// String returns a debug string
func (ul *UserLogoutEvt) String() string {
	return fmt.Sprintf("UserLogoutEvt[Login: %s]", ul.Login)
}

// ServerStartEvt fires when the server is started
type ServerStartedEvt struct{}

// String returns a debug string
func (e *ServerStartedEvt) String() string {
	return "ServerStartEvt"
}

//TODO: Finish device lost plumbing
//TODO: DeviceConnectedEvt
