package cmd

import "fmt"

type ButtonRelease struct {
	ButtonLocalID  string
	ButtonGlobalID string
	DeviceName     string
	DeviceLocalID  string
	DeviceGlobalID string
}

func (c *ButtonRelease) FriendlyString() string {
	return fmt.Sprintf("Device [%s] \"%s\" release %s [%s]",
		c.DeviceGlobalID, c.DeviceName, c.ButtonLocalID, c.ButtonGlobalID)
}
func (c *ButtonRelease) String() string {
	return "cmd.ButtonRelease"
}
