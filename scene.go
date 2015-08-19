package gohome

type Scene struct {
	Identifiable
	Device   *Device
	Commands []Command
}

func (s *Scene) Execute() {
	for _, c := range s.Commands {
		s.Device.ExecuteCommand(c)
	}
}
