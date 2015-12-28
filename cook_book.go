package gohome

type CookBook struct {
	Identifiable
	LogoURL  string
	Triggers []Trigger
	Actions  []Action
}
