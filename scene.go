package gohome

type Scene struct {
	Identifiable
	Commands []Command `json:"-"`
}

func (s *Scene) Execute() {
	for _, c := range s.Commands {
		c.Execute()
	}
}
