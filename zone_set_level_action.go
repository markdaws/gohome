package gohome

import "fmt"

type ZoneSetLevelAction struct {
	ZoneID string
	Level  float32
}

func (a *ZoneSetLevelAction) Type() string {
	return "gohome.ZoneSetLevelAction"
}

func (a *ZoneSetLevelAction) Name() string {
	return "Set Zone Level"
}

func (a *ZoneSetLevelAction) Description() string {
	return "Sets the zone level to the specified value"
}

func (a *ZoneSetLevelAction) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			ID:          "Level",
			Name:        "Intensity Level",
			Description: "The target intensity for the zone",
			Type:        "float",
			Required:    true,
		},
		Ingredient{
			ID:          "ZoneID",
			Name:        "Zone ID",
			Description: "The ID of the target zone",
			Type:        "string",
			Required:    true,
		},
	}
}

func (a *ZoneSetLevelAction) Execute(s *System) error {
	zone, ok := s.Zones[a.ZoneID]
	if !ok {
		return fmt.Errorf("Unknown ZoneID %s", a.ZoneID)
	}
	return zone.SetLevel(a.Level)
}

func (a *ZoneSetLevelAction) New() Action {
	return &ZoneSetLevelAction{}
}
