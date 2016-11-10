package cmd

import "fmt"

type ZoneTurnOff struct {
	ID          string
	ZoneAddress string
	ZoneID      string
	ZoneName    string
}

func (c *ZoneTurnOff) FriendlyString() string {
	return fmt.Sprintf("Zone [%s] \"%s\" turn on", c.ZoneID, c.ZoneName)
}
func (c *ZoneTurnOff) String() string {
	return "cmd.ZoneTurnOff"
}
