package www

type jsonZone struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
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
