package gohome

type FuncCommand struct {
	Friendly    string
	CommandType CommandType
	Func        func() error
}

func (c *FuncCommand) Execute() error {
	return c.Func()
}
func (c *FuncCommand) FriendlyString() string {
	return c.Friendly
}
func (c *FuncCommand) String() string {
	return c.FriendlyString()
}
func (c *FuncCommand) CMDType() CommandType {
	return c.CommandType
}
