package zone

// Type represents the type of outputs the zone controls
type Type uint32

const (
	// ZTUnknown indicates we don't know the kind of output this zone is controlling
	ZTUnknown Type = iota

	// ZTLight represents a light
	ZTLight

	// ZTShade represents a shade
	ZTShade

	// ZTSwitch represents a switch
	ZTSwitch
)

// TypeFromString returns a Type based in the value passed in to the function
func TypeFromString(zt string) Type {
	switch zt {
	case "light":
		return ZTLight
	case "switch":
		return ZTSwitch
	case "shade":
		return ZTShade
	case "unknown":
		return ZTUnknown
	default:
		return ZTUnknown
	}
}

// ToString returns a string representation of the specified Type parameter. NOTE: this
// value is for internal use only, don't use it anywhere in your app, it can change at any time
func (zt Type) ToString() string {
	switch zt {
	case ZTLight:
		return "light"
	case ZTSwitch:
		return "switch"
	case ZTShade:
		return "shade"
	default:
		return "unknown"
	}
}
