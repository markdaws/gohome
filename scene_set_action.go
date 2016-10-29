package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
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
			ID:          "SceneID",
			Name:        "Scene ID",
			Description: "The ID of the Scene to set",
			Type:        "string",
			Required:    true,
			Reference:   "scene",
		},
	}
}

func (a *SceneSetAction) Execute(s *System) error {
	scene, ok := s.Scenes[a.SceneID]
	if !ok {
		return fmt.Errorf("Unknown Scene ID %s", a.SceneID)
	}

	desc := fmt.Sprintf("Set Scene: %s", scene.Name)
	return s.Services.CmdProcessor.Enqueue(NewCommandGroup(desc, &cmd.SceneSet{
		SceneID:   scene.ID,
		SceneName: scene.Name,
	}))
}

func (a *SceneSetAction) New() Action {
	return &SceneSetAction{}
}
