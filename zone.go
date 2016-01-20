package gohome

type Zone struct {
	LocalID     string
	GlobalID    string
	Name        string
	Description string
	Device      Device
	Type        ZoneType
	Output      OutputType

	//TODO: Serialize/deserialize
	Controller string

	//TODO: Describe discrete, continuous, max, min, step e.g. on/off vs dimmable

	//TODO: Bulbs, Shades, etc
	//TODO: RLevel, GLevel, BLevel when OTRGB
}

/*

*/
