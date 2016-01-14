package gohome

import "fmt"

type StringCommandAction struct {
	DeviceID string
	Command  string
	Friendly string
}

func (a *StringCommandAction) Type() string {
	return "gohome.StringCommandAction"
}

func (a *StringCommandAction) Name() string {
	return "String Command Action"
}

func (a *StringCommandAction) Description() string {
	return "Sends the specified string command to the specified device"
}

func (a *StringCommandAction) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			ID:          "DeviceID",
			Name:        "Device ID",
			Description: "The ID of the device to send the command to",
			Type:        "string",
			Required:    true,
		},
		Ingredient{
			ID:          "Command",
			Name:        "Command String",
			Description: "The string to send to the specified device",
			Type:        "string",
			Required:    true,
		},
		Ingredient{
			ID:          "Friendly",
			Name:        "Description",
			Description: "A human friendly description of the command that will appear in the event log",
			Type:        "string",
			Required:    true,
		},
	}
}

func (a *StringCommandAction) Execute(s *System) error {
	device, ok := s.Devices[a.DeviceID]
	if !ok {
		return fmt.Errorf("Unknown Device ID %s", a.DeviceID)
	}

	c := StringCommand{
		Value:    a.Command,
		Device:   device,
		Type:     CTDeviceSendCommand,
		Friendly: a.Friendly,
	}
	return c.Execute()
}

func (a *StringCommandAction) New() Action {
	return &StringCommandAction{}
}
