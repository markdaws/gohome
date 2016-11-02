package store

import "github.com/markdaws/gohome"

type systemJSON struct {
	Version     string              `json:"version"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Scenes      []sceneJSON         `json:"scenes"`
	Devices     []deviceJSON        `json:"devices"`
	Recipes     []gohome.RecipeJSON `json:"recipes"`
	Users       []userJSON          `json:"users"`
}

type userJSON struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	HashedPwd string `json:"hashedPwd"`
	Salt      string `json:"salt"`
}

type buttonJSON struct {
	Address     string `json:"address"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type zoneJSON struct {
	Address     string `json:"address"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DeviceID    string `json:"deviceId"`
	Type        string `json:"type"`
	Output      string `json:"output"`
}

type sceneJSON struct {
	Address     string        `json:"address"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Commands    []commandJSON `json:"commands"`
}

type connPoolJSON struct {
	Name     string `json:"name"`
	PoolSize int32  `json:"poolSize"`
}

type deviceJSON struct {
	ID              string        `json:"id"`
	Address         string        `json:"address"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	ModelNumber     string        `json:"modelNumber"`
	ModelName       string        `json:"modelName"`
	SoftwareVersion string        `json:"softwareVersion"`
	HubID           string        `json:"hubId"`
	Buttons         []buttonJSON  `json:"buttons"`
	Zones           []zoneJSON    `json:"zones"`
	Sensors         []sensorJSON  `json:"sensors"`
	DeviceIDs       []string      `json:"deviceIds"`
	Auth            *authJSON     `json:"auth"`
	ConnPool        *connPoolJSON `json:"connPool"`
}

type sensorAttrJSON struct {
	Name     string            `json:"name"`
	Value    string            `json:"value"`
	DataType string            `json:"dataType"`
	States   map[string]string `json:"states"`
}

type sensorJSON struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Address     string         `json:"address"`
	DeviceID    string         `json:"deviceId"`
	Attr        sensorAttrJSON `json:"attr"`
}

type authJSON struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type commandJSON struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
}
