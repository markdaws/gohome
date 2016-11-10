package cmd

import "fmt"

type StringCommand struct {
	ID    string
	Value string
	Args  []interface{}
}

func (c *StringCommand) String() string {
	return fmt.Sprintf(c.Value, c.Args...)
}
