package example

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type cmdBuilder struct {
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
	//
	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		// For the abstract command we then return a cmd.Func instance that translates
		// the abstract command in to commands to send to the example device
		return &cmd.Func{
			Func: func() error {
				// Make some example command string the example hardware needs
				finalCmd := fmt.Sprintf("ZONESET %s, %f", command.ZoneID, command.Level.Value)

				// In a real extension you would perform a network action to send the
				// command here, for example look at extensions/belkin/cmd_builder.go
				// in this example we just print the command
				fmt.Println(finalCmd)

				return nil
			},
		}, nil
	case *cmd.ZoneTurnOn:
		return &cmd.Func{
			Func: func() error {
				fmt.Println("Sending TurnOn command to example1 hardware")
				return nil
			},
		}, nil
	case *cmd.ZoneTurnOff:
		return &cmd.Func{
			Func: func() error {
				fmt.Println("Sending TurnOff command to example1 hardware")
				return nil
			},
		}, nil
	case *cmd.ButtonPress:
		// Hardware doesn't support button presses, just return nil
		return nil, nil
	case *cmd.ButtonRelease:
		return nil, nil
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
