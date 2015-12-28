package gohome

type Scene struct {
	Identifiable
	Commands []Command
}

func (s *Scene) Execute() error {
	for _, c := range s.Commands {
		c.Execute()
	}

	//TODO: error
	return nil
}
