package cmd

import "fmt"

type ButtonRelease struct {
	ButtonAddress string
	ButtonID      string
	DeviceName    string
	DeviceAddress string
	DeviceID      string
}

func (c *ButtonRelease) FriendlyString() string {
	return fmt.Sprintf("Device [%s] \"%s\" release %s [%s]",
		c.DeviceID, c.DeviceName, c.ButtonAddress, c.ButtonID)
}
func (c *ButtonRelease) String() string {
	return "cmd.ButtonRelease"
}
