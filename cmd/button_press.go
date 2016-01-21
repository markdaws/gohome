package cmd

import "fmt"

type ButtonPress struct {
	ButtonAddress string
	ButtonID      string
	DeviceName    string
	DeviceAddress string
	DeviceID      string
}

func (c *ButtonPress) FriendlyString() string {
	return fmt.Sprintf("Device [%s] \"%s\" press %s [%s]",
		c.DeviceID, c.DeviceName, c.ButtonAddress, c.ButtonID)
}
func (c *ButtonPress) String() string {
	return "cmd.ButtonPress"
}
