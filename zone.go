package gohome

type Zone struct {
	Identifiable
	Type ZoneType
	//TODO: Describe discrete, continuous, max, min, step e.g. on/off vs dimmable
	SetCommand Command
}

func (z *Zone) Set(value float32) error {
	z.SetCommand.Execute(value)
	//TODO: error
	return nil
}
