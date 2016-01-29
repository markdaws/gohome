package www

type jsonScene struct {
	Address     string        `json:"address"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Managed     bool          `json:"managed"`
	Commands    []jsonCommand `json:"commands"`
}
type scenes []jsonScene

func (slice scenes) Len() int {
	return len(slice)
}
func (slice scenes) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}
func (slice scenes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
