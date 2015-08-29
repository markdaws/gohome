package gohome

//TODO: turn enum to string representation
type CommandType uint32

const (
	CTUnknown CommandType = iota
	CTZoneSetLevel
	CTDevicePressButton
	CTDeviceReleaseButton
	CTSystemSetScene
)
