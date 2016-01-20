package cmd

type Func struct {
	Func     func() error
	Friendly string
}

func (c *Func) Execute() error {
	return c.Func()
}
func (c *Func) FriendlyString() string {
	return c.Friendly
}
func (c *Func) String() string {
	return "cmd.Func"
}
