package gohome

type Trigger interface {
	Start() (<-chan bool, <-chan bool)
	Stop()
	GetIngredients() []Ingredient
	GetName() string
	GetDescription() string
}

//TODO: TrueTrigger? Just executes
