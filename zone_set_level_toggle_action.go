package gohome

import "fmt"

type ZoneSetLevelToggleAction struct {
	ZoneID      string
	FirstLevel  float32
	SecondLevel float32

	second bool
}

func (a *ZoneSetLevelToggleAction) Type() string {
	return "gohome.ZoneSetLevelToggleAction"
}

func (a *ZoneSetLevelToggleAction) Name() string {
	return "Set Zone Level Toggle"
}

func (a *ZoneSetLevelToggleAction) Description() string {
	return "Toggles the specified zone level between two values"
}

func (a *ZoneSetLevelToggleAction) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			ID:          "FirstLevel",
			Name:        "First Intensity Level",
			Description: "The first target intensity for the zone",
			Type:        "float",
			Required:    true,
		},
		Ingredient{
			ID:          "SecondLevel",
			Name:        "Second Intensity Level",
			Description: "The second target intensity for the zone",
			Type:        "float",
			Required:    true,
		},
		Ingredient{
			ID:          "ZoneID",
			Name:        "Zone ID",
			Description: "The ID of the target zone",
			Type:        "string",
			Required:    true,
			Reference:   "zone",
		},
	}
}

func (a *ZoneSetLevelToggleAction) Execute(s *System) error {
	zone, ok := s.Zones[a.ZoneID]
	if !ok {
		return fmt.Errorf("Unknown ZoneID %s", a.ZoneID)
	}

	var level float32
	if a.second {
		a.second = false
		level = a.SecondLevel
	} else {
		a.second = true
		level = a.FirstLevel
	}

	return s.CmdProcessor.Enqueue(&ZoneSetLevelCommand{
		Zone:  zone,
		Level: level,
	})
}

func (a *ZoneSetLevelToggleAction) New() Action {
	return &ZoneSetLevelToggleAction{}
}
