package cmd

import "fmt"

type ZoneSetLevel struct {
	ID          string
	ZoneAddress string
	ZoneID      string
	ZoneName    string
	Level       Level
}

func (c *ZoneSetLevel) FriendlyString() string {
	return fmt.Sprintf("Zone [%s] \"%s\" set to %.2f%%", c.ZoneID, c.ZoneName, c.Level.Value)
}
func (c *ZoneSetLevel) String() string {
	return "cmd.ZoneSetLevel"
}
