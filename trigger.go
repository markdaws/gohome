package gohome

type Trigger interface {
	Type() string
	Ingredients() []Ingredient
	Name() string
	Description() string
	New() Trigger
	Init() (<-chan bool, bool)
	ProcessEvent(Event) bool
}
