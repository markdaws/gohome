package gohome

import (
	"testing"

	"github.com/markdaws/gohome"
)

func TestXxx(t *testing.T) {
	s := &gohome.System{
		Identifiable: gohome.Identifiable{
			ID:          "Sys1",
			Name:        "Test System",
			Description: "Test System description",
		},
		Devices: map[string]*gohome.Device{},
		Zones: map[string]*gohome.Zone{
			"7": &gohome.Zone{
				Identifiable: gohome.Identifiable{
					ID:          "zid1",
					Name:        "Kitchen Lights",
					Description: "Main lights in the kitchen",
				},
			},
		},
	}
	d := &gohome.Device{
		Identifiable: gohome.Identifiable{
			ID:          "",
			Name:        "",
			Description: ""},
		System:     s,
		Connection: nil,
	}
	gohome.ParseCommandString(d, "~OUTPUT,7,1,25.00")
}
