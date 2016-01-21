package gohome

import "github.com/markdaws/gohome/cmd"

type Scene struct {
	Address     string
	ID          string
	Name        string
	Description string
	Commands    []cmd.Command
}
