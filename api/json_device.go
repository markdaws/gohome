package api

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
	AddressRequired bool          `json:"addressRequired"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	ModelNumber     string        `json:"modelNumber"`
	ModelName       string        `json:"modelName"`
	SoftwareVersion string        `json:"softwareVersion"`
	Zones           []jsonZone    `json:"zones"`
	ConnPool        *jsonConnPool `json:"connPool"`
	Type            string        `json:"type"`
	Sensors         []jsonSensor  `json:"sensors"`
	Auth            *jsonAuth     `json:"auth"`
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
