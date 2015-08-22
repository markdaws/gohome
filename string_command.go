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
	c.Device.Connection.Send([]byte(str))
}
