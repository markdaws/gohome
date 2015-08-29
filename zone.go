package gohome

//TODO: Should have some type e.g. Light/Shade/Other ...

type Zone struct {
	Identifiable
	Type       ZoneType
	SetCommand Command
}

func (z *Zone) Set(value float32) {
	//somehow need to insert value
	z.SetCommand.Execute(value)
}

//getlevel
