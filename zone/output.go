package zone

type Output uint32

const (
	OTContinuous Output = iota
	OTBinary
	OTRGB
	OTUnknown
)

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

//TODO: Gen this automatically
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
