package zone

import (
	"fmt"

	"github.com/markdaws/gohome/validation"
)

type Zone struct {
	Address     string
	ID          string
	Name        string
	Description string
	DeviceID    string
	Type        Type
	Output      Output

	//TODO: Describe max, min, step e.g. on/off vs dimmable
	//TODO: Value presets?
}

func (z *Zone) String() string {
	return fmt.Sprintf("Zone[ID: %s, Name: %s, Address: %s]", z.ID, z.Name, z.Address)
}

func (z *Zone) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if z.ID == "" {
		errors.Add("required field", "ID")
	}

	if z.Name == "" {
		errors.Add("required field", "Name")
	}

	if z.DeviceID == "" {
		errors.Add("required field", "DeviceID")
	}

	// Zones do not require an address, if they are connected to a device that only
	// controls one zone it is not required
	if errors.Has() {
		return errors
	}
	return nil
}
