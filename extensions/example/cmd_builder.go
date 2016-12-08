package example

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/cmd"
)

type cmdBuilder struct {
	System      *gohome.System
	Device      *gohome.Device
	ModelNumber string
}

func (b *cmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	// In this example our extension supports different types of hardware, so
	// we have different build methods for each piece of hardware
	switch b.ModelNumber {
	case "example.hardware.1":
		return b.buildHardwareOneCommands(c)
	case "example.hardware.2":
		return b.buildHardwareTwoCommands(c)
	default:
		return nil, errors.New("unsupported hardware found")
	}
}

func (b *cmdBuilder) buildHardwareOneCommands(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.FeatureSetAttrs:

		// Grab the feature that has changed
		f := b.System.FeatureByID(command.FeatureID)
		if f == nil {
			return nil, fmt.Errorf("invalid feature ID: %s", command.FeatureID)
		}

		// Get the device that owns the feature
		dev := b.System.DeviceByID(f.DeviceID)
		if dev == nil {
			return nil, fmt.Errorf("invalid device ID: %s", f.DeviceID)
		}
		_ = dev

		// For the abstract command we then return a cmd.Func instance that translates
		// the abstract command in to commands to send to the example device
		return &cmd.Func{
			Func: func() error {
				// Depending on the features you exported from your hardware, you will perform
				// different actions here.  For example we exported a light zone and a sensor in
				// this example. So we can loop through the updated attrs passed in the event
				// and if the attr is a Brightness attribute, we can set the new value.

				// NOTE: There may be many attributes updated in one event, it's up to you
				// to check for each one and perform the appropriate actions
				for _, attribute := range command.Attrs {
					if attribute.Type == attr.ATBrightness {
						// Pretend to set the light brightness

						// Make some example command string the example hardware needs
						// IMPORTANT: note how we cast the value to float32, you need to cast the value
						// to the type you want before using it, since it is stored as interface{}
						finalCmd := fmt.Sprintf("LIGHT SET %s, %f", f.ID, attribute.Value.(float32))

						// In a real extension you would perform a network action to send the
						// command here, for example look at extensions/belkin/cmd_builder.go
						// in this example we just print the command
						fmt.Println(finalCmd)
					}
				}

				return nil
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
	return nil, nil
}

func (b *cmdBuilder) buildHardwareTwoCommands(c cmd.Command) (*cmd.Func, error) {
	// Same details as the above function, but just send different commands for the
	// different type of hardware you support
	fmt.Println("Sending command to example2 hardware")
	return nil, nil
}
