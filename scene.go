package gohome

import "github.com/markdaws/gohome/cmd"

type Scene struct {
	Address     string
	ID          string
	Name        string
	Description string

	// Managed is true if gohome knows about the items the scene is controlling,
	// false otherwise. For example when you import a Lutron scene, we don't know
	// what zones it is controlling, it is managed by the Lutron device, but if you
	// create a GoHome scene, we do know what zones it will affect
	Managed  bool
	Commands []cmd.Command
}
