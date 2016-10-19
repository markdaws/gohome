package belkin

// DeviceType represents an identifier for the type of Belkin device you want to
// scan the network for
type DeviceType string

const (
	DTBridge      DeviceType = "urn:Belkin:device:bridge:1"
	DTSwitch                 = "urn:Belkin:device:controllee:1"
	DTMotion                 = "urn:Belkin:device:sensor:1"
	DTMaker                  = "urn:Belkin:device:Maker:1"
	DTInsight                = "urn:Belkin:device:insight:1"
	DTLightSwitch            = "urn:Belkin:device:lightswitch:1"
)
