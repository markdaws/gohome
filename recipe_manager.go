/*
TODO:
- validate input on client and server
- specify ranges for ingredients
- return useful error messages for invalid ingredients
*/
package gohome

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"
	"unicode/utf8"
)

//TODO: make internal
type RecipeManager struct {
	CookBooks []*CookBook
	System    *System
	Recipes   []*Recipe

	dataPath       string
	eventBroker    EventBroker
	triggerFactory map[string]func() Trigger
	actionFactory  map[string]func() Action
}

func (rm *RecipeManager) Init(eb EventBroker, dataPath string) error {
	rm.eventBroker = eb
	rm.dataPath = dataPath
	rm.CookBooks = loadCookBooks(dataPath)
	rm.triggerFactory = buildTriggerFactory(rm.CookBooks)
	rm.actionFactory = buildActionFactory(rm.CookBooks)

	recipes, err := rm.loadRecipes(dataPath)
	if err != nil {
		return err
	}

	rm.Recipes = recipes
	for _, recipe := range recipes {
		rm.RegisterAndStart(recipe)
	}
	return nil
}

func (rm *RecipeManager) RegisterAndStart(r *Recipe) {
	rm.eventBroker.AddConsumer(r.Trigger.(EventConsumer))
	r.Start()
}

func (rm *RecipeManager) UnregisterAndStop(r *Recipe) {
	rm.eventBroker.RemoveConsumer(r.Trigger.(EventConsumer))
	r.Stop()
}

func (rm *RecipeManager) RecipeByID(id string) *Recipe {
	for _, recipe := range rm.Recipes {
		if recipe.ID == id {
			return recipe
		}
	}
	return nil
}

func (rm *RecipeManager) EnableRecipe(r *Recipe, enabled bool) error {
	oldEnabled := r.Trigger.Enabled()
	if oldEnabled == enabled {
		return nil
	}

	r.Trigger.SetEnabled(enabled)
	return rm.SaveRecipe(r, false)
}

type recipeJSON struct {
	ID          string
	Name        string
	Description string
	Enabled     bool `json:"enabled"`
	Trigger     triggerWrapper
	Action      actionWrapper
}

type triggerWrapper struct {
	Type    string                 `json:"type"`
	Trigger map[string]interface{} `json:"fields"`
}

type actionWrapper struct {
	Type   string                 `json:"type"`
	Action map[string]interface{} `json:"fields"`
}

type ErrUnmarshalRecipe struct {
	ParamID     string
	ErrorType   string
	Description string
}

func (e ErrUnmarshalRecipe) Error() string {
	return e.ParamID + " - " + e.Description
}

func (rm *RecipeManager) UnmarshalNewRecipe(data map[string]interface{}) (*Recipe, error) {
	if _, ok := data["name"]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "name",
			ErrorType:   "missingParam",
			Description: "required field",
		}
	}
	name, ok := data["name"].(string)
	if !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "name",
			ErrorType:   "invalidData",
			Description: "must be a string",
		}
	}
	if utf8.RuneCountInString(name) == 0 {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "name",
			ErrorType:   "invalidData",
			Description: "must have a value",
		}
	}

	if _, ok = data["description"]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "description",
			ErrorType:   "missingParam",
			Description: "required field",
		}
	}
	desc, ok := data["description"].(string)
	if !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "description",
			ErrorType:   "invalidData",
			Description: "must be a string",
		}
	}
	if utf8.RuneCountInString(desc) == 0 {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "description",
			ErrorType:   "invalidData",
			Description: "must have a value",
		}
	}

	if _, ok = data["trigger"]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "trigger",
			ErrorType:   "missingParam",
			Description: "required field",
		}
	}
	triggerData, ok := data["trigger"].(map[string]interface{})
	if !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "trigger",
			ErrorType:   "invalidData",
			Description: "must be an object",
		}
	}

	if _, ok = triggerData["id"]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "trigger.id",
			ErrorType:   "missingParam",
			Description: "required field",
		}
	}
	triggerID, ok := triggerData["id"].(string)
	if !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "trigger.id",
			ErrorType:   "invalidData",
			Description: "must be a string",
		}
	}

	if _, ok = triggerData["ingredients"]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "trigger.ingredients",
			ErrorType:   "missingParam",
			Description: "required field",
		}
	}
	triggerIngredients, ok := triggerData["ingredients"].(map[string]interface{})
	if !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "trigger.ingredients",
			ErrorType:   "invalidData",
			Description: "must be an object",
		}
	}

	if _, ok = rm.triggerFactory[triggerID]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "trigger.id",
			ErrorType:   "invalidData",
			Description: fmt.Sprintf("unknown trigger ID: %s", triggerID),
		}
	}

	if _, ok = data["action"]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "action",
			ErrorType:   "missingParam",
			Description: "required field",
		}
	}
	actionData, ok := data["action"].(map[string]interface{})
	if !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "action",
			ErrorType:   "invalidData",
			Description: "must be an object",
		}
	}
	if _, ok = actionData["id"]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "action.id",
			ErrorType:   "missingParam",
			Description: "required field",
		}
	}
	actionID, ok := actionData["id"].(string)
	if !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "action.id",
			ErrorType:   "invalidData",
			Description: "must be a string",
		}
	}

	if _, ok = actionData["ingredients"]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "action.ingredients",
			ErrorType:   "missingParam",
			Description: "required field",
		}
	}
	actionIngredients, ok := actionData["ingredients"].(map[string]interface{})
	if !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "action.ingredients",
			ErrorType:   "invalidData",
			Description: "must be an object",
		}
	}

	if _, ok = rm.actionFactory[actionID]; !ok {
		return nil, &ErrUnmarshalRecipe{
			ParamID:     "action.id",
			ErrorType:   "invalidData",
			Description: fmt.Sprintf("unknown trigger ID: %s", actionID),
		}
	}

	trigger := rm.triggerFactory[triggerID]()
	action := rm.actionFactory[actionID]()

	err := setIngredients(trigger, triggerIngredients, reflect.ValueOf(trigger).Elem())
	if err != nil {
		err.ParamID = "trigger." + err.ParamID
		return nil, err
	}
	err = setIngredients(action, actionIngredients, reflect.ValueOf(action).Elem())
	if err != nil {
		err.ParamID = "action." + err.ParamID
		return nil, err
	}

	enabled := true
	recipe, rErr := NewRecipe(name, desc, enabled, trigger, action, rm.System)
	return recipe, rErr
}

func (rm *RecipeManager) SaveRecipe(r *Recipe, appendTo bool) error {
	// Since Trigger and Action are interfaces, we need to also save the underlying
	// concrete type to the JSON file so we can unmarshal to the correct type later

	out := recipeJSON{}
	out.ID = r.ID
	out.Name = r.Name
	out.Description = r.Description
	out.Enabled = r.Trigger.Enabled()

	out.Trigger = triggerWrapper{Type: r.Trigger.Type(), Trigger: getIngredientValueMap(r.Trigger, reflect.ValueOf(r.Trigger).Elem())}
	out.Action = actionWrapper{Type: r.Action.Type(), Action: getIngredientValueMap(r.Action, reflect.ValueOf(r.Action).Elem())}

	b, err := json.Marshal(out)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(rm.recipePath(r), b, 0644)
	if err != nil {
		return err
	}

	if appendTo {
		rm.Recipes = append(rm.Recipes, r)
	}
	return nil
}

func (rm *RecipeManager) recipePath(r *Recipe) string {
	return filepath.Join(rm.dataPath, r.ID+".json")
}

func (rm *RecipeManager) DeleteRecipe(r *Recipe) error {
	err := os.Remove(rm.recipePath(r))
	if err != nil {
		return err
	}

	for i, recipe := range rm.Recipes {
		if recipe.ID == r.ID {
			rm.Recipes, rm.Recipes[len(rm.Recipes)-1] = append(rm.Recipes[:i], rm.Recipes[i+1:]...), nil
			break
		}
	}

	rm.UnregisterAndStop(r)
	return nil
}

func getIngredientValueMap(i Ingredientor, v reflect.Value) map[string]interface{} {
	values := make(map[string]interface{})
	for _, ingredient := range i.Ingredients() {
		// Want to store duration as ms, so need to massage
		var value interface{}
		if ingredient.Type == "duration" {
			value = int64(time.Duration(v.FieldByName(ingredient.ID).Int()) / time.Millisecond)
		} else {
			value = v.FieldByName(ingredient.ID).Interface()
		}
		values[ingredient.ID] = value
	}
	return values
}

func (rm *RecipeManager) loadRecipes(path string) ([]*Recipe, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	recipes := make([]*Recipe, 0)
	for _, fileInfo := range files {
		filepath := filepath.Join(path, fileInfo.Name())
		recipe, err := rm.loadRecipe(filepath)
		if err != nil {
			//TODO: log error
			fmt.Println(err)
			continue
		}

		//fmt.Printf("appending %+v", recipe)
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (rm *RecipeManager) loadRecipe(path string) (*Recipe, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var recipeWrapper recipeJSON
	err = json.Unmarshal(b, &recipeWrapper)
	if err != nil {
		return nil, err
	}

	recipe := &Recipe{}
	recipe.system = rm.System
	recipe.ID = recipeWrapper.ID
	recipe.Name = recipeWrapper.Name
	recipe.Description = recipeWrapper.Description

	trigger, err := rm.makeTrigger(recipeWrapper.Trigger.Type, recipeWrapper.Trigger.Trigger)
	if err != nil {
		return nil, err
	}
	trigger.SetEnabled(recipeWrapper.Enabled)

	action, err := rm.makeAction(recipeWrapper.Action.Type, recipeWrapper.Action.Action)
	if err != nil {
		return nil, err
	}

	recipe.Trigger = trigger
	recipe.Action = action
	return recipe, nil
}

func (rm *RecipeManager) makeTrigger(triggerID string, triggerIngredients map[string]interface{}) (Trigger, error) {
	trigger := rm.triggerFactory[triggerID]()

	err := setIngredients(trigger, triggerIngredients, reflect.ValueOf(trigger).Elem())
	if err != nil {
		return nil, err
	}
	return trigger, nil
}

func (rm *RecipeManager) makeAction(actionID string, actionIngredients map[string]interface{}) (Action, error) {
	action := rm.actionFactory[actionID]()

	err := setIngredients(action, actionIngredients, reflect.ValueOf(action).Elem())
	if err != nil {
		return nil, err
	}
	return action, nil
}

func setIngredients(i Ingredientor, ingredientValues map[string]interface{}, s reflect.Value) *ErrUnmarshalRecipe {
	for _, ingredient := range i.Ingredients() {
		_, ok := ingredientValues[ingredient.ID]

		if !ok && ingredient.Required {
			return &ErrUnmarshalRecipe{
				ParamID:     ingredient.ID,
				ErrorType:   "invalidData",
				Description: "must have a value",
			}
		}

		if ok {
			field := s.FieldByName(ingredient.ID)
			switch ingredient.Type {
			case "string":
				value, ok := ingredientValues[ingredient.ID].(string)
				if !ok {
					return &ErrUnmarshalRecipe{
						ParamID:     ingredient.ID,
						ErrorType:   "invalidData",
						Description: "must be a string",
					}
				}

				/*
					//TODO: Have a reference field that we can use to validate
					//ids are valid objects in the system
					if ingredient.Reference != "" {
						//TODO: needed?
						parts := strings.Split(ingredient.Reference)
						for i, v := range parts {
							if i == len(parts) - 1 {
								//The last part tells us the kind of
							}
						}
					}*/
				field.SetString(value)

			case "boolean":
				value, ok := ingredientValues[ingredient.ID].(bool)
				if !ok {
					return &ErrUnmarshalRecipe{
						ParamID:     ingredient.ID,
						ErrorType:   "invalidData",
						Description: "must be a boolean (true or false)",
					}
				}
				field.SetBool(value)

			case "integer":
				value, ok := ingredientValues[ingredient.ID].(float64)
				if !ok {
					return &ErrUnmarshalRecipe{
						ParamID:     ingredient.ID,
						ErrorType:   "invalidData",
						Description: "must be an integer",
					}
				}
				field.SetInt(int64(value))

			case "float":
				value, ok := ingredientValues[ingredient.ID].(float64)
				if !ok {
					return &ErrUnmarshalRecipe{
						ParamID:     ingredient.ID,
						ErrorType:   "invalidData",
						Description: "must be a floating point number",
					}
				}
				field.SetFloat(value)

			case "duration":
				value, ok := ingredientValues[ingredient.ID].(float64)
				if !ok {
					return &ErrUnmarshalRecipe{
						ParamID:     ingredient.ID,
						ErrorType:   "invalidData",
						Description: "must be an integer representing milliseconds",
					}
				}
				field.Set(reflect.ValueOf(time.Duration(int64(value)) * time.Millisecond))

			case "datetime":
				//TODO: implement

			default:
				return &ErrUnmarshalRecipe{
					ParamID:     ingredient.ID,
					ErrorType:   "invalidData",
					Description: "unknown ingredient ID",
				}
			}
		}
	}
	return nil
}

func loadCookBooks(dataPath string) []*CookBook {
	// For every cook book we support, add to this list, at some point these can
	// be defined in a config file or in a DB
	cookBooks := []*CookBook{
		{
			ID:          "1",
			Name:        "Lutron Smart Bridge Pro",
			Description: "Cook up some goodness for the Lutron Smart Bridge Pro",
			LogoURL:     "lutron_400x400.png",
			Triggers: []Trigger{
				// New triggers need to be added to this slice
				&ButtonClickTrigger{},
				&TimeTrigger{},
			},
			Actions: []Action{
				// New actions need to be added to this slice
				&ZoneSetLevelAction{},
				&ZoneSetLevelToggleAction{},
				&SceneSetAction{},
				&SceneSetToggleAction{},
				&StringCommandAction{},
			},
		},
	}
	return cookBooks
}

func buildTriggerFactory(cookBooks []*CookBook) map[string]func() Trigger {
	factory := make(map[string]func() Trigger)
	for _, cookBook := range cookBooks {
		for _, trigger := range cookBook.Triggers {
			factory[trigger.Type()] = trigger.New
		}
	}
	return factory
}

func buildActionFactory(cookBooks []*CookBook) map[string]func() Action {
	factory := make(map[string]func() Action)
	for _, cookBook := range cookBooks {
		for _, action := range cookBook.Actions {
			factory[action.Type()] = action.New
		}
	}
	return factory
}
