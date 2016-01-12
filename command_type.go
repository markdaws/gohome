package gohome

type CommandType uint32

const (
	CTUnknown CommandType = iota
	CTZoneSetLevel
	CTDevicePressButton
	CTDeviceReleaseButton
	CTDeviceSendCommand
	CTSystemSetScene
)
