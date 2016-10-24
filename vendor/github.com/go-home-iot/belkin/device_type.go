package belkin

// DeviceType represents an identifier for the type of Belkin device you want to
// scan the network for
type DeviceType string

const (
	// DTBridge - belkin bridge
	DTBridge DeviceType = "urn:Belkin:device:bridge:1"

	// DTSwitch - belkin switch
	DTSwitch = "urn:Belkin:device:controllee:1"

	// DTMotion - belkin motion sensor
	DTMotion = "urn:Belkin:device:sensor:1"

	// DTMaker - belkin maker
	DTMaker = "urn:Belkin:device:Maker:1"

	// DTInsight - belkin insight
	DTInsight = "urn:Belkin:device:insight:1"

	// DTLightSwitch - belkin light switch
	DTLightSwitch = "urn:Belkin:device:lightswitch:1"
)
