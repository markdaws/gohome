package gohome

import "github.com/markdaws/gohome/cmd"

type Scene struct {
	LocalID     string
	GlobalID    string
	Name        string
	Description string
	Commands    []cmd.Command
}
