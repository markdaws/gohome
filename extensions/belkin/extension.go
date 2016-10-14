package belkin

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
	//Belkin WeMo Insight
	lookupTable["f7c029v2"] = &discoverer{System: sys}
}

func (e *extension) RegisterImporters(sys *gohome.System, lookupTable map[string]gohome.Importer) {
}

func NewExtension() *extension {
	return &extension{}
}
