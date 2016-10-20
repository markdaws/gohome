package zone

// Output represents the different types of output a zone can control
type Output uint32

const (
	// OTContinuous represents a zone that can have continuous values from min to max set on it e.g. 31.2
	OTContinuous Output = iota

	// OTBinary represents a zone that can have only two states e.g. on/of or open/close
	OTBinary

	// OTRGB represents a zone that controls RGB outputs
	OTRGB

	// OTUnknown represents a zone whos output type we do not know
	OTUnknown
)

// OutputFromString returns a Output instance corresponding to the string passed in to the function
func OutputFromString(ot string) Output {
	switch ot {
	case "continuous":
		return OTContinuous
	case "binary":
		return OTBinary
	case "rgb":
		return OTRGB
	case "unknown":
		return OTUnknown
	default:
		return OTUnknown
	}
}

// ToString returns a string represneting the Output parameter passed in to the function. NOTE: this
// value is for internal use only, don't use it anywhere in your app, it can change at any time
func (ot Output) ToString() string {
	switch ot {
	case OTContinuous:
		return "continuous"
	case OTBinary:
		return "binary"
	case OTRGB:
		return "rgb"
	case OTUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}
