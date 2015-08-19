package gohome

type StringCommand struct {
	Value  string
	Device Device
}

func (c *StringCommand) Execute() {
	c.Device.Connection.Send([]byte(c.Value))
}
