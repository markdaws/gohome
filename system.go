package gohome

type System struct {
	ID          string
	Name        string
	Description string
	Devices     map[string]*Device
	Scenes      map[string]*Scene
	Zones       map[string]*Zone
}
