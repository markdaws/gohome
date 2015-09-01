package gohome

//TODO: Allow arrays? specify max size optional
type Ingredient struct {
	Identifiable
	Type     string //string, number (float64), bool
	Required bool
}
