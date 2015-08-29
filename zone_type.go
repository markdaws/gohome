package gohome

type ZoneType uint32

const (
	ZTUnknown ZoneType = iota
	ZTLight
	ZTShade

	//Garage door
	//Sprinkler
	//Heat?

	//TODO: What are the most common smart devices, add support
	//TODO: Document how to add support for a new device
)
