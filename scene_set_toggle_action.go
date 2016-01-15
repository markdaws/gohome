package gohome

import "fmt"

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

	if a.second {
		a.second = false
		return second.Execute()
	} else {
		a.second = true
		return first.Execute()
	}
}

func (a *SceneSetToggleAction) New() Action {
	return &SceneSetToggleAction{}
}
