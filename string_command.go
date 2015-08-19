package gohome

type StringCommand struct {
	Value string
}

func (c *StringCommand) Data() []byte {
	return []byte(c.Value)
}
