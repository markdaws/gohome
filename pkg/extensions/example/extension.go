/*
package example shows an example implementation of an extension.  Extensions are used to interface
with different types of hardware. To allow goHOME to interface with new hardware e.g. Philips Hue bulbs
we would create a new extension that knows how to monitor and talk to Hue bulbs

IMPORTANT: Once you have implemented this, you need to update the code in gohome/intg/intg.go to
include this new extension
*/
package example

import (
	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/gohome"
)

type extension struct {
	// You must composite the NullExtension in your type, this gives default
	// methods incase you do not need to implement them.
	gohome.NullExtension
}

func (e *extension) Name() string {
	// Return a unique string for your extension, just returning the package
	// name will suffice
	return "example"
}

func (e *extension) BuilderForDevice(sys *gohome.System, d *gohome.Device) cmd.Builder {
	// There are several abstract commands that gohome may send to your extension, you need
	// to look at the commands and return a cmd.Builder instance that knows how to turn the
	// abstract command to a device specific set of commands.

	// This function will be called for all devices in the system, including ones this extension
	// doesn't control.  You need to look at the device and see if you own it, if not then just
	// return nil.

	// It may be that you support several different pieces of hardware, you can return different
	// builders for each device, or even check fields like d.SoftwareVersion to see if you need
	// to return a different builder for different versions of software installed on the device

	// In the example below we check ModelNumber to see if this is a device we support, we support
	// two different pieces of hardware
	switch d.ModelNumber {
	case "example.hardware.1":
		return &cmdBuilder{ModelNumber: d.ModelNumber, Device: d, System: sys}
	case "example.hardware.2":
		return &cmdBuilder{ModelNumber: d.ModelNumber, Device: d, System: sys}
	default:
		// This device is not one that we know how to control, return nil
		return nil
	}
}

func (e *extension) NetworkForDevice(sys *gohome.System, d *gohome.Device) gohome.Network {
	// This method will be called for all devices in the system, including ones this extension doesn't
	// know how to control, check the Device and if you don't know how to control it, return nil

	// The only reason you need to implement this method is if your hardware uses a TCP connection pool,
	// if you just send HTTP requests, simply return nil from this function

	// If your hardware supports TCP connection pooling then this function is called when the
	// system needs to know how to open a connection to your device.

	switch d.ModelNumber {
	case "example.hardware.1":
		return &network{Device: d}
	default:
		return nil
	}
}

func (e *extension) EventsForDevice(sys *gohome.System, d *gohome.Device) *gohome.ExtEvents {
	// If your device needs to react to events in the system, or can produce events, such as monitoring
	// a changing value that you want to report back to the system, you need to implement this function

	// Again, this function is called for every device in the system, you need to look at the device and
	// determine if your extension owns the hardware, if not, return nil

	switch d.ModelNumber {
	case "example.hardware.1":
		evts := &gohome.ExtEvents{}
		// If you don't produce any events, just leave this nil
		evts.Producer = &producer{
			Device: d,
			System: sys,
		}

		// If you don't consume any events, just leave this nil
		evts.Consumer = &consumer{
			Device: d,
			System: sys,
		}
		return evts
	default:
		return nil
	}
}

func (e *extension) Discovery(sys *gohome.System) gohome.Discovery {
	// The discovery function allows an extension to document what type of hardware it supports. The information
	// you return here will appear in the import page on the app.
	return &discovery{}
}

// You need a NewExtension() function to return your extension instance
func NewExtension() *extension {
	return &extension{}
}
