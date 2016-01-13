package gohome

type Zone struct {
	Identifiable
	Type   ZoneType
	Output OutputType
	//TODO: Describe discrete, continuous, max, min, step e.g. on/off vs dimmable
	setCommand   func(args ...interface{}) Command
	cmdProcessor CommandProcessor
}

func (z *Zone) SetLevel(value float32) error {
	cmd := z.setCommand(value)
	z.cmdProcessor.Enqueue(cmd)
	return nil
}
