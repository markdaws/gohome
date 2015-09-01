package gohome

type CookBook struct {
	Identifiable
	Triggers []Trigger
	Actions  []Action
}
