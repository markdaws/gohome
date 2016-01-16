package gohome

type Action interface {
	Execute(*System) error
	Type() string
	Name() string
	Description() string
	Ingredients() []Ingredient
	New() Action
}
