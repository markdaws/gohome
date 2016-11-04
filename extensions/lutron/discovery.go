package lutron

import (
	"errors"

	"github.com/markdaws/gohome"
)

type discovery struct {
	System *gohome.System
}

func (d *discovery) Discoverers() []gohome.DiscovererInfo {
	return []gohome.DiscovererInfo{gohome.DiscovererInfo{
		ID:          "lutron.l-bdgpro2-wh",
		Name:        "Lutron Smart Bridge Pro",
		Description: "Discover Lutron Smart Bridge Pro hubs",
		Type:        "FromString",
	}}
}

func (d *discovery) DiscovererFromID(ID string) gohome.Discoverer {
	switch ID {
	case "lutron.l-bdgpro2-wh":
		return &discoverer{System: d.System}
	default:
		return nil
	}
}

type discoverer struct {
	System *gohome.System
}

func (d *discoverer) ScanDevices(sys *gohome.System) (*gohome.DiscoveryResults, error) {
	return nil, errors.New("unsupported")
}
func (d *discoverer) FromString(body string) (*gohome.DiscoveryResults, error) {
	//TODO: Fix should not suck into a system ...

	importer := &importer{System: d.System}
	err := importer.FromString(d.System, body)
	return nil, err
}
