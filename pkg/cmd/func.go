package cmd

// Func represents a function along with a friendly name for that function
type Func struct {
	Func     func() error
	Friendly string
}

// Execute calls the func associated with this Func instance
func (c *Func) Execute() error {
	return c.Func()
}

// FriendlyString returns a human readable friendly name for this function e.g. "Zone Set to 10%"
func (c *Func) FriendlyString() string {
	return c.Friendly
}

// String returns a human readable friendly name for this function e.g. "Zone Set to 10%"
func (c *Func) String() string {
	if c.Friendly == "" {
		return "cmd.Func"
	}
	return c.Friendly
}
