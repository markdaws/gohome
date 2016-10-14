package lutron

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type extension struct {
}

func (e *extension) RegisterCmdBuilders(sys *gohome.System, lookupTable map[string]cmd.Builder) {
	builder := &cmdBuilder{System: sys}
	lookupTable[builder.ID()] = builder
}

func (e *extension) RegisterDiscoverers(sys *gohome.System, lookupTable map[string]gohome.Discoverer) {
	//TODO: Implement
}

func (e *extension) RegisterImporters(sys *gohome.System, lookupTable map[string]gohome.Importer) {
	importer := &importer{System: sys}

	// Register a mapping from a moelnumber to an importer.  We can
	// register as many importers as we want here for multiple model numbers
	lookupTable["l-bdgpro2-wh"] = importer
}

func NewExtension() *extension {
	return &extension{}
}
