package gohome

type Zone struct {
	Identifiable
	Type ZoneType
	//TODO: Describe discrete, continuous, max, min, step e.g. on/off vs dimmable

	//TODO: Why public, calling this from server.go?
	//setCommand Command
	setCommand   func(args ...interface{}) Command
	cmdProcessor CommandProcessor
}

func (z *Zone) SetLevel(value float32) error {
	cmd := z.setCommand(value)
	z.cmdProcessor.Enqueue(cmd)
	//z.SetCommand.Execute(value)
	//TODO: error
	return nil
}
