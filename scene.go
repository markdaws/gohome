package gohome

type Scene struct {
	Identifiable
	Commands     []Command
	cmdProcessor CommandProcessor
}

func (s *Scene) Execute() error {
	for _, c := range s.Commands {
		s.cmdProcessor.Enqueue(c)
	}
	return nil
}
