package gohome

import "fmt"

type Command interface {
	Execute(args ...interface{}) error
	String() string
	FriendlyString() string
	GetType() CommandType
}

type CommandBuilderParams struct {
	CommandType  CommandType
	Zone         *Zone
	Intensity    float64
	Device       *Device
	SourceDevice *Device
	//TODO: Needed?
	ComponentID string
}

//TODO: The command builder will depend on the device that it will be sent to
//makes more sense to be a function inside a device - then would have a factory
//that returns the right kind of builder for a certain device, keyed of model number?
func BuildCommand(p CommandBuilderParams) Command {
	switch p.CommandType {
	case CTZoneSetLevel:
		return &StringCommand{
			Device:   p.Device,
			Friendly: fmt.Sprintf("Zone \"%s\" set to %.2f%%", p.Zone.Name, p.Intensity),
			Value:    fmt.Sprintf("#OUTPUT,%s,1,%.2f\r\n", p.Zone.ID, p.Intensity),
			Type:     p.CommandType,
		}

	case CTDevicePressButton:
		return &StringCommand{
			Device:   p.Device,
			Friendly: fmt.Sprintf("Device \"%s\" press button %s", p.SourceDevice.Name, p.ComponentID),
			Value:    fmt.Sprintf("#DEVICE,%s,%s,3\r\n", p.SourceDevice.Name, p.ComponentID),
			Type:     p.CommandType,
		}

	case CTDeviceReleaseButton:
		return &StringCommand{
			Device:   p.Device,
			Friendly: fmt.Sprintf("Device \"%s\" release button %s", p.SourceDevice.Name, p.ComponentID),
			Value:    fmt.Sprintf("#DEVICE,%s,%s,4\r\n", p.SourceDevice.Name, p.ComponentID),
			Type:     p.CommandType,
		}

	default:
		return nil
	}
}
