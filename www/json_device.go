package www

type jsonDevice struct {
	LocalID     string `json:"localId"`
	GlobalID    string `json:"globalId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ModelNumber string `json:"modelNumber"`
	//TODO: Stream
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