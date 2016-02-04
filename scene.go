package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/validation"
)

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

func (s *Scene) DeleteCommand(i int) error {
	if i < 0 || i >= len(s.Commands) {
		return fmt.Errorf("invalid command index")
	}

	s.Commands, s.Commands[len(s.Commands)-1] = append(s.Commands[:i], s.Commands[i+1:]...), nil
	return nil
}

func (s *Scene) AddCommand(c cmd.Command) error {
	s.Commands = append(s.Commands, c)
	return nil
}

func (s *Scene) Validate() *validation.Errors {
	errors := &validation.Errors{}

	if s.Name == "" {
		errors.Add("required field", "Name")
	}

	if errors.Has() {
		return errors
	}
	return nil
}
