package api

import "strings"

type jsonButton struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"fullName"`
	Description string `json:"description"`
	Address     string `json:"address"`
}
type buttons []jsonButton

func (slice buttons) Len() int {
	return len(slice)
}
func (slice buttons) Less(i, j int) bool {
	return slice[i].FullName < slice[j].FullName
}
func (slice buttons) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type jsonCommand struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}

type jsonConnPool struct {
	Name     string `json:"name"`
	PoolSize int32  `json:"poolSize"`
}

type jsonAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type jsonDevice struct {
	ID              string        `json:"id"`
	Address         string        `json:"address"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	ModelNumber     string        `json:"modelNumber"`
	ModelName       string        `json:"modelName"`
	SoftwareVersion string        `json:"softwareVersion"`
	Zones           []jsonZone    `json:"zones"`
	Sensors         []jsonSensor  `json:"sensors"`
	Buttons         []jsonButton  `json:"buttons"`
	ConnPool        *jsonConnPool `json:"connPool"`
	Type            string        `json:"type"`
	Auth            *jsonAuth     `json:"auth"`
	IsDupe          bool          `json:"isDupe"`
}
type devices []jsonDevice

func (slice devices) Len() int {
	return len(slice)
}
func (slice devices) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}
func (slice devices) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type jsonUIField struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Default     string `json:"default"`
	Required    bool   `json:"required"`
}

type jsonDiscovererInfo struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	PreScanInfo string        `json:"preScanInfo"`
	UIFields    []jsonUIField `json:"uiFields"`
}

type jsonMonitorGroup struct {
	TimeoutInSeconds int      `json:"timeoutInSeconds"`
	SensorIDs        []string `json:"sensorIds"`
	ZoneIDs          []string `json:"zoneIds"`
}

type jsonZoneLevel struct {
	Value float32 `json:"value"`
	R     byte    `json:"r"`
	G     byte    `json:"g"`
	B     byte    `json:"b"`
}
type jsonMonitorGroupResponse struct {
	Sensors map[string]jsonSensorAttr `json:"sensors"`
	Zones   map[string]jsonZoneLevel  `json:"zones"`
}

type jsonRecipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
type jsonRecipes []jsonRecipe

func (slice jsonRecipes) Len() int {
	return len(slice)
}
func (slice jsonRecipes) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}
func (slice jsonRecipes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type jsonScene struct {
	Address     string        `json:"address"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Managed     bool          `json:"managed"`
	Commands    []jsonCommand `json:"commands"`
}
type scenes []jsonScene

func (slice scenes) Len() int {
	return len(slice)
}
func (slice scenes) Less(i, j int) bool {
	return strings.ToLower(slice[i].Name) < strings.ToLower(slice[j].Name)
}
func (slice scenes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type jsonSensorAttr struct {
	Name     string            `json:"name"`
	Value    string            `json:"value"`
	DataType string            `json:"dataType"`
	States   map[string]string `json:"states"`
}

type jsonSensor struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Address     string         `json:"address"`
	DeviceID    string         `json:"deviceId"`
	Attr        jsonSensorAttr `json:"attr"`
	IsDupe      bool           `json:"isDupe"`
}
type sensors []jsonSensor

func (slice sensors) Len() int {
	return len(slice)
}
func (slice sensors) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}
func (slice sensors) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type jsonZone struct {
	Address     string `json:"address"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DeviceID    string `json:"deviceId"`
	Type        string `json:"type"`
	Output      string `json:"output"`
	IsDupe      bool   `json:"isDupe"`
}
type zones []jsonZone

func (slice zones) Len() int {
	return len(slice)
}
func (slice zones) Less(i, j int) bool {
	return strings.ToLower(slice[i].Name) < strings.ToLower(slice[j].Name)
}
func (slice zones) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
