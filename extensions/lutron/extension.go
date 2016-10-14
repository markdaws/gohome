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

func NewExtension() *extension {
	return &extension{}
}
