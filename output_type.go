package gohome

type OutputType uint32

const (
	OTContinuous OutputType = iota
	OTBinary
	OTRGB
	OTUnknown
)

func OutputTypeFromString(ot string) OutputType {
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

//TODO: Gen this automatically
func (ot OutputType) ToString() string {
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
