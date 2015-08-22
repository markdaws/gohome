package gohome

type System struct {
	Identifiable
	Devices map[string]*Device
	Scenes  map[string]*Scene
	Zones   map[string]*Zone
}
