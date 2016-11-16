package lutron

import "fmt"

// Event represents an event fired by a device
type Event interface {
	fmt.Stringer
}

// UnknownEvt is an event which we don't parse into a concrete type, can inspect the Msg field
type UnknownEvt struct {
	Msg string
}

func (e UnknownEvt) String() string {
	return e.Msg
}

// ZoneLevelEvt is when the system reports a zones current level
type ZoneLevelEvt struct {
	Address string
	Level   float32
}

func (e ZoneLevelEvt) String() string {
	return fmt.Sprintf("address: %s, level:%f", e.Address, e.Level)
}

// BtnPressEvt is fired when a button is pressed, not this is only the press action not also
// the release, this is not a click
type BtnPressEvt struct {
	Address       string
	DeviceAddress string
}

func (e BtnPressEvt) String() string {
	return fmt.Sprintf("button address: %s, device address: %s", e.Address, e.DeviceAddress)
}

// BtnReleaseEvt is fired when a button is released
type BtnReleaseEvt struct {
	Address       string
	DeviceAddress string
}

func (e BtnReleaseEvt) String() string {
	return fmt.Sprintf("button address: %s, device address: %s", e.Address, e.DeviceAddress)
}
