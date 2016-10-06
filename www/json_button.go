package www

type jsonButton struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`

	// ClientID is an ID assigned on the client to the zone if the zone
	// hasn't been created yet but still needs to be referenced uniquely
	// by the client.
	ClientID string `json:"clientId,omitempty"`
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
