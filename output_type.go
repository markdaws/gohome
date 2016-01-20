package gohome

type OutputType uint32

const (
	OTContinuous OutputType = iota
	OTBinary
	OTUnknown

	//TODO: RGB ?
)

func OutputTypeFromString(ot string) OutputType {
	switch ot {
	case "continuous":
		return OTContinuous
	case "binary":
		return OTBinary
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
	case OTUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}
