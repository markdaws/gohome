package api

type jsonSensor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	DeviceID    string `json:"deviceId"`
	ClientID    string `json:"clientId"`
	//TODO: attributes
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
