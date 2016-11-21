package api

import (
	"strings"

	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/feature"
)

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
	ID              string             `json:"id"`
	Address         string             `json:"address"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	ModelNumber     string             `json:"modelNumber"`
	ModelName       string             `json:"modelName"`
	SoftwareVersion string             `json:"softwareVersion"`
	ConnPool        *jsonConnPool      `json:"connPool"`
	Type            string             `json:"type"`
	Auth            *jsonAuth          `json:"auth"`
	HubID           string             `json:"hubId"`
	DeviceIDs       []string           `json:"deviceIds"`
	IsDupe          bool               `json:"isDupe"`
	Features        []*feature.Feature `json:"features"`
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
	FeatureIDs       []string `json:"featureIds"`
}

type jsonMonitorGroupResponse struct {
	Features map[string]map[string]*attr.Attribute `json:"features"`
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
