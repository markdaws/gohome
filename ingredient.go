package gohome

type Ingredient struct {
	ID          string
	Name        string
	Description string
	Type        string //string,integer,float,duration,datetime,boolean
	Required    bool
}

type Ingredientor interface {
	Ingredients() []Ingredient
}
