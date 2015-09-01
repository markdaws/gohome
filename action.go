package gohome

type Action interface {
	//TODO: Return error
	Execute()
	GetName() string
	GetDescription() string
	GetIngredients() []Ingredient
}

/*
type PrintAction struct {
}

func (a *PrintAction) Execute() {
	fmt.Println("I am a print action")
}
*/

type FuncAction struct {
	Name        string
	Description string
	Func        func()
}

func (a *FuncAction) GetName() string {
	return a.Name
}

func (a *FuncAction) GetDescription() string {
	return a.Description
}

func (a *FuncAction) GetIngredients() []Ingredient {
	//TODO: Where does func come from?
	return nil
}

func (a *FuncAction) Execute() {
	a.Func()
}

type SetSceneAction struct {
	Scene *Scene
}

func (a *SetSceneAction) GetName() string {
	return "Set Scene"
}

func (a *SetSceneAction) GetDescription() string {
	return "Sets the specified scene"
}

func (a *SetSceneAction) GetIngredients() []Ingredient {
	//TODO: how to enumerate scene
	return nil
}

func (a *SetSceneAction) Execute() {
	a.Scene.Execute()
}

//TODO: Allow multiple zone/level values?
type ZoneSetLevelAction struct {
	Zone  *Zone
	Level float32
}

func (a *ZoneSetLevelAction) GetName() string {
	return "Set Zone Level"
}

func (a *ZoneSetLevelAction) GetDescription() string {
	return "Sets the zone level to the specified value"
}

func (a *ZoneSetLevelAction) GetIngredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			Identifiable: Identifiable{
				ID:          "Level",
				Name:        "Intensity Level",
				Description: "The target intensity for the zone",
			},
			Type: "number",
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "Zone",
				Name:        "Zone",
				Description: "The target zone",
			},
			Type: "string", //TODO: Zone?
		},
	}
}

func (a *ZoneSetLevelAction) Execute() {
	a.Zone.Set(a.Level)
}
