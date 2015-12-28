package gohome

type Trigger interface {
	Start() (<-chan bool, <-chan bool)
	Stop()
	Type() string
	Ingredients() []Ingredient
	Name() string
	Description() string
	Enabled() bool
	SetEnabled(bool)
	New() Trigger
}
