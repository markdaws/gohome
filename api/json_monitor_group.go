package api

type jsonMonitorGroup struct {
	TimeoutInSeconds int      `json:"timeoutInSeconds"`
	SensorIDs        []string `json:"sensorIds"`
	ZoneIDs          []string `json:"zoneIds"`
}

//TODO: Move
type jsonZoneLevel struct {
	Value float32 `json:"value"`
	R     byte    `json:"r"`
	G     byte    `json:"g"`
	B     byte    `json:"b"`
}
type jsonMonitorGroupResponse struct {
	Sensors map[string]jsonSensorAttr
	Zones   map[string]jsonZoneLevel
}
