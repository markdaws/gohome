package gohome

type Ingredient struct {
	Identifiable
	Type     string //string,integer,float,duration,datetime,boolean
	Required bool
}

type Ingredientor interface {
	Ingredients() []Ingredient
}
