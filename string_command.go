package gohome

import "fmt"

type StringCommand struct {
	Value  string
	Device *Device
}

//TODO: return error
func (c *StringCommand) Execute(args ...interface{}) {
	str := fmt.Sprintf(c.Value, args...)
	fmt.Println("Setting command:", str)

	//TODO: Should use connection pool, don't assume connection is
	//just open to send on
	c.Device.Connection.Send([]byte(str))
}
