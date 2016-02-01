package gohome

// Action is an interface that describes a unit of execution. An action could be something
// like 'set zone level', or 'apply scene'. The Ingredients method returns a list of settings
// for the action
type Action interface {
	Execute(*System) error
	Type() string
	Name() string
	Description() string
	Ingredients() []Ingredient
	New() Action
}
