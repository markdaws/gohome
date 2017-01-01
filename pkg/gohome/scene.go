package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/pkg/cmd"
	"github.com/markdaws/gohome/pkg/validation"
)

// Scene represents a target state for the system. A scene is just a list of commands
// that should be executed when the scene is chosen
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

// DeleteCommand deletes a command from the scene
func (s *Scene) DeleteCommand(cmdID string) error {
	index := -1
	for i, c := range s.Commands {
		if c.GetID() == cmdID {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("invalid command ID")
	}

	s.Commands, s.Commands[len(s.Commands)-1] = append(s.Commands[:index], s.Commands[index+1:]...), nil
	return nil
}

// AddCommand adds a new command to the scene
func (s *Scene) AddCommand(c cmd.Command) error {
	s.Commands = append(s.Commands, c)
	return nil
}

// Validate verfies the scene is in a good state
func (s *Scene) Validate() *validation.Errors {
	//TODO: Verify that there isn't an infinite loop where scene A -> B -> C -> A otherwise
	//app will crash

	errors := &validation.Errors{}

	if s.Name == "" {
		errors.Add("required field", "Name")
	}

	if errors.Has() {
		return errors
	}
	return nil
}
