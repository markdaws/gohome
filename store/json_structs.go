package store

import "github.com/markdaws/gohome/feature"

type systemJSON struct {
	Version     string       `json:"version"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Scenes      []sceneJSON  `json:"scenes"`
	Devices     []deviceJSON `json:"devices"`
	Users       []userJSON   `json:"users"`
}

type areaJSON struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ParentID    string   `json:"parentId"`
	FeatureIDs  []string `json:"featureIds"`
	AreaIDs     []string `json:"areaIds"`
}

type userJSON struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	HashedPwd string `json:"hashedPwd"`
	Salt      string `json:"salt"`
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
	ID              string             `json:"id"`
	Address         string             `json:"address"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	ModelNumber     string             `json:"modelNumber"`
	ModelName       string             `json:"modelName"`
	SoftwareVersion string             `json:"softwareVersion"`
	HubID           string             `json:"hubId"`
	DeviceIDs       []string           `json:"deviceIds"`
	Auth            *authJSON          `json:"auth"`
	ConnPool        *connPoolJSON      `json:"connPool"`
	Features        []*feature.Feature `json:"features"`
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
