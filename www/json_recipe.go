package www

type jsonRecipe struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
type jsonRecipes []jsonRecipe

func (slice jsonRecipes) Len() int {
	return len(slice)
}
func (slice jsonRecipes) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}
func (slice jsonRecipes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
