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
			Identifiable: Identifiable{
				ID:          "FirstLevel",
				Name:        "First Intensity Level",
				Description: "The first target intensity for the zone",
			},
			Type:     "float",
			Required: true,
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "SecondLevel",
				Name:        "Second Intensity Level",
				Description: "The second target intensity for the zone",
			},
			Type:     "float",
			Required: true,
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "ZoneID",
				Name:        "Zone ID",
				Description: "The ID of the target zone",
			},
			Type:     "string",
			Required: true,
		},
	}
}

func (a *ZoneSetLevelToggleAction) Execute(s *System) error {
	zone, ok := s.Zones[a.ZoneID]
	if !ok {
		return fmt.Errorf("Unknown ZoneID %s", a.ZoneID)
	}

	if a.second {
		a.second = false
		return zone.Set(a.SecondLevel)
	} else {
		a.second = true
		return zone.Set(a.FirstLevel)
	}
}

func (a *ZoneSetLevelToggleAction) New() Action {
	return &ZoneSetLevelToggleAction{}
}
