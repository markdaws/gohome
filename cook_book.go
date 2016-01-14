package gohome

type CookBook struct {
	ID          string
	Name        string
	Description string
	LogoURL     string
	Triggers    []Trigger
	Actions     []Action
}
