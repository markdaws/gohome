package gohome

type Zone struct {
	Identifiable
	Type       ZoneType
	SetCommand Command
}

func (z *Zone) Set(value float32) {
	z.SetCommand.Execute(value)
}

//TODO: Support multiple channels e.g. r/g/b vs. just intensity
