package www

type jsonDevice struct {
	Address     string     `json:"address"`
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ModelNumber string     `json:"modelNumber"`
	Token       string     `json:"token"`
	ClientID    string     `json:"clientId,omitempty"`
	Zones       []jsonZone `json:"zones"`
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
