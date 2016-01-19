package gohome

type Scene struct {
	LocalID     string
	GlobalID    string
	Name        string
	Description string
	Commands    []Command
}
