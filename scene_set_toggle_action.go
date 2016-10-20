package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
)

type SceneSetToggleAction struct {
	FirstSceneID  string
	SecondSceneID string

	second bool
}

func (a *SceneSetToggleAction) Type() string {
	return "gohome.SceneSetToggleAction"
}

func (a *SceneSetToggleAction) Name() string {
	return "Set Scene Toggle"
}

func (a *SceneSetToggleAction) Description() string {
	return "Toggles between setting the two specified scenes"
}

func (a *SceneSetToggleAction) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			ID:          "FirstSceneID",
			Name:        "First Scene ID",
			Description: "The ID of the first Scene to set",
			Type:        "string",
			Required:    true,
			Reference:   "scene",
		},
		Ingredient{
			ID:          "SecondSceneID",
			Name:        "Second Scene ID",
			Description: "The ID of the second Scene to set",
			Type:        "string",
			Required:    true,
			Reference:   "scene",
		},
	}
}

func (a *SceneSetToggleAction) Execute(s *System) error {
	first, ok := s.Scenes[a.FirstSceneID]
	if !ok {
		return fmt.Errorf("Unknown First Scene ID %s", a.FirstSceneID)
	}

	second, ok := s.Scenes[a.SecondSceneID]
	if !ok {
		return fmt.Errorf("Unknown Second Scene ID %s", a.SecondSceneID)
	}

	var scene *Scene
	if a.second {
		a.second = false
		scene = second
	} else {
		a.second = true
		scene = first
	}

	desc := fmt.Sprintf("Toggle Scene: %s", scene.Name)
	return s.CmdProcessor.Enqueue(NewCommandGroup(desc, &cmd.SceneSet{
		SceneID:   scene.ID,
		SceneName: scene.Name,
	}))
}

func (a *SceneSetToggleAction) New() Action {
	return &SceneSetToggleAction{}
}
