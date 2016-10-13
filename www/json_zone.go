package www

type jsonZone struct {
	Address     string `json:"address"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DeviceID    string `json:"deviceId"`
	Type        string `json:"type"`
	Output      string `json:"output"`

	// ClientID is an ID assigned on the client to the zone if the zone
	// hasn't been created yet but still needs to be referenced uniquely
	// by the client.
	ClientID string `json:"clientId,omitempty"`
}
type zones []jsonZone

func (slice zones) Len() int {
	return len(slice)
}
func (slice zones) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}
func (slice zones) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
