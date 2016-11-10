package cmd

import "fmt"

type ZoneTurnOn struct {
	ID          string
	ZoneAddress string
	ZoneID      string
	ZoneName    string
}

func (c *ZoneTurnOn) FriendlyString() string {
	return fmt.Sprintf("Zone [%s] \"%s\" turn on", c.ZoneID, c.ZoneName)
}
func (c *ZoneTurnOn) String() string {
	return "cmd.ZoneTurnOn"
}
