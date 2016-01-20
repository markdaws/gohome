package gohome

type Zone struct {
	Address     string
	ID          string
	Name        string
	Description string
	Device      Device
	Type        ZoneType
	Output      OutputType
	Controller  string

	//TODO: Describe max, min, step e.g. on/off vs dimmable
}
