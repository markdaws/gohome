package gohome

type Zone struct {
	Identifiable
	SetCommand Command
}

func (z *Zone) Set(value float32) {
	//somehow need to insert value
	z.SetCommand.Execute(value)
}

//getlevel
