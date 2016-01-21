package zone

type Zone struct {
	Address     string
	ID          string
	Name        string
	Description string
	DeviceID    string
	Type        Type
	Output      Output
	Controller  string

	//TODO: Describe max, min, step e.g. on/off vs dimmable
}
