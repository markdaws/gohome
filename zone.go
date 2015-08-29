package gohome

type Zone struct {
	Identifiable
	Type       ZoneType
	SetCommand Command
}

func (z *Zone) Set(value float32) {
	z.SetCommand.Execute(value)
}
