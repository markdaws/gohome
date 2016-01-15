package gohome

type Scene struct {
	LocalID      string
	GlobalID     string
	Name         string
	Description  string
	Commands     []Command
	cmdProcessor CommandProcessor
}

func (s *Scene) Execute() error {
	for _, c := range s.Commands {
		s.cmdProcessor.Enqueue(c)
	}
	return nil
}
