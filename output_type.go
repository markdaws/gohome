package gohome

type OutputType uint32

const (
	OTContinuous OutputType = iota
	OTBinary
)

//TODO: Gen this automatically
func (ot OutputType) ToString() string {
	switch ot {
	case OTContinuous:
		return "continuous"
	case OTBinary:
		return "binary"
	default:
		return "unknown"
	}
}
