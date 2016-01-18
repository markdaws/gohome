package gohome

type Zone struct {
	LocalID     string
	GlobalID    string
	Name        string
	Description string
	Device      Device
	Type        ZoneType
	Output      OutputType
	//TODO: Describe discrete, continuous, max, min, step e.g. on/off vs dimmable
}
