package cmd

import "fmt"

type StringCommand struct {
	Value string
	Args  []interface{}
}

func (c *StringCommand) String() string {
	return fmt.Sprintf(c.Value, c.Args...)
}
