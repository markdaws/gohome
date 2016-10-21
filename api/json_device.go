package api

type jsonCmdBuilder struct {
	ID string `json:"id"`
}

type jsonConnPool struct {
	Name     string `json:"name"`
	PoolSize int32  `json:"poolSize"`
}

type jsonDevice struct {
	Address         string          `json:"address"`
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	ModelNumber     string          `json:"modelNumber"`
	ModelName       string          `json:"modelName"`
	SoftwareVersion string          `json:"softwareVersion"`
	Token           string          `json:"token"`
	ClientID        string          `json:"clientId,omitempty"`
	Zones           []jsonZone      `json:"zones"`
	CmdBuilder      *jsonCmdBuilder `json:"cmdBuilder"`
	ConnPool        *jsonConnPool   `json:"connPool"`
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
