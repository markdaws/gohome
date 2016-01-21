package zone

type Type uint32

const (
	ZTUnknown Type = iota
	ZTLight
	ZTShade

	//Garage door
	//Sprinkler
	//Heat?

	//TODO: What are the most common smart devices, add support
	//TODO: Document how to add support for a new device
)

func TypeFromString(zt string) Type {
	switch zt {
	case "light":
		return ZTLight
	case "shade":
		return ZTShade
	case "unknown":
		return ZTUnknown
	default:
		return ZTUnknown
	}
}

//TODO: Gen this automatically
func (zt Type) ToString() string {
	switch zt {
	case ZTLight:
		return "light"
	case ZTShade:
		return "shade"
	default:
		return "unknown"
	}
}
