package gohome

import (
	"errors"
	"fmt"
)

type SceneSetAction struct {
	SceneID string
}

func (a *SceneSetAction) Type() string {
	return "gohome.SceneSetAction"
}

func (a *SceneSetAction) Name() string {
	return "Set Scene"
}

func (a *SceneSetAction) Description() string {
	return "Sets the specified scene"
}

func (a *SceneSetAction) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			Identifiable: Identifiable{
				ID:          "SceneID",
				Name:        "Scene ID",
				Description: "The ID of the Scene to set",
			},
			Type:     "string",
			Required: true,
		},
	}
}

func (a *SceneSetAction) Execute(s *System) error {
	scene, ok := s.Scenes[a.SceneID]
	if !ok {
		return errors.New(fmt.Sprintf("Unknown Scene ID %s", a.SceneID))
	}
	return scene.Execute()
}

func (a *SceneSetAction) New() Action {
	return &SceneSetAction{}
}
