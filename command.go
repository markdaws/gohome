package gohome

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Command interface {
	Execute(args ...interface{})
	String() string
	FriendlyString() string
}

func ParseCommandString(d *Device, cmd string) Command {
	switch {
	case strings.HasPrefix(cmd, "~OUTPUT"),
		strings.HasPrefix(cmd, "#OUTPUT"):
		return parseZoneCommand(d, cmd)

	case strings.HasPrefix(cmd, "~DEVICE"),
		strings.HasPrefix(cmd, "#DEVICE"):
		return parseDeviceCommand(d, cmd)

	default:
		//TODO: Error?
		//fmt.Println("unknown: " + cmd)
		return nil
	}
}

func parseDeviceCommand(d *Device, cmd string) Command {
	matches := regexp.MustCompile("[~|#]DEVICE,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(cmd)
	if matches == nil || len(matches) != 4 {
		fmt.Println("no matches")
		return nil
	}

	deviceID := matches[1]
	componentID := matches[2]
	cmdID := matches[3]
	sourceDevice := d.System.Devices[deviceID]
	if sourceDevice == nil {
		//TODO: Error?
		return nil
	}

	var ct CommandType
	switch cmdID {
	case "3":
		ct = CTDevicePressButton
	case "4":
		ct = CTDeviceReleaseButton
	default:
		ct = CTUnknown
	}

	return BuildCommand(CommandBuilderParams{
		Device:       d,
		CommandType:  ct,
		SourceDevice: sourceDevice,
		ComponentID:  componentID,
	})
}

func parseZoneCommand(d *Device, cmd string) Command {
	matches := regexp.MustCompile("[~|?]OUTPUT,([^,]+),([^,]+),(.+)\r\n").FindStringSubmatch(cmd)
	if matches == nil || len(matches) != 4 {
		fmt.Println("no matches")
		return nil
	}

	zoneID := matches[1]
	cmdID := matches[2]
	intensity, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		//TODO: Error
		return nil
	}

	z := d.System.Zones[zoneID]
	if z == nil {
		//TODO: Error log
		return nil
	}

	var ct CommandType
	switch cmdID {
	case "1":
		ct = CTZoneSetLevel
	default:
		ct = CTUnknown
	}

	return BuildCommand(CommandBuilderParams{
		Device:      d,
		CommandType: ct,
		Intensity:   intensity,
		Zone:        z,
	})
}

type CommandBuilderParams struct {
	CommandType  CommandType
	Zone         *Zone
	Intensity    float64
	Device       *Device
	SourceDevice *Device
	ComponentID  string
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
		}

	case CTDevicePressButton:
		return &StringCommand{
			Device:   p.Device,
			Friendly: fmt.Sprintf("Device \"%s\" press button %s", p.SourceDevice.Name, p.ComponentID),
			Value:    fmt.Sprintf("#DEVICE,%s,%s,3\r\n", p.SourceDevice.Name, p.ComponentID),
		}

	case CTDeviceReleaseButton:
		return &StringCommand{
			Device:   p.Device,
			Friendly: fmt.Sprintf("Device \"%s\" release button %s", p.SourceDevice.Name, p.ComponentID),
			Value:    fmt.Sprintf("#DEVICE,%s,%s,4\r\n", p.SourceDevice.Name, p.ComponentID),
		}

	default:
		return nil
	}
}
