package gohome

import "fmt"

type StringCommand struct {
	Value    string
	Friendly string
	Device   *Device
}

//TODO: return error
func (c *StringCommand) Execute(args ...interface{}) {
	str := fmt.Sprintf(c.Value, args...)
	fmt.Println("Setting command:", str)

	//TODO: Should use connection pool, don't assume connection is
	//just open to send on
	c.Device.Connection.Send([]byte(str))
}

func (c *StringCommand) String() string {
	return c.Value
}

func (c *StringCommand) FriendlyString() string {
	return c.Friendly
}
