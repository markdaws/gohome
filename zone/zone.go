package zone

import "github.com/markdaws/gohome/validation"

type Zone struct {
	Address     string
	ID          string
	Name        string
	Description string
	DeviceID    string
	Type        Type
	Output      Output
	//TODO: use zone controller
	Controller string

	//TODO: Describe max, min, step e.g. on/off vs dimmable
	//TODO: Value presets?
}

func (z *Zone) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if z.Name == "" {
		errors.Add("required field", "Name")
	}

	if z.Address == "" {
		errors.Add("required field", "Address")
	}

	if z.DeviceID == "" {
		errors.Add("required field", "DeviceID")
	}

	if errors.Has() {
		return errors
	}
	return nil
	//TODO: Type/Output/Controller?
}
