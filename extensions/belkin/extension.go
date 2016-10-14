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
	//TODO: Implement

	/*
		case "F7C029V2":
		responses, err := belkin.Scan(belkin.DTInsight, 5)
		_ = responses
		_ = err
		return nil, fmt.Errorf("//TODO:not implemented")
	*/
}

func (e *extension) RegisterImporters(sys *gohome.System, lookupTable map[string]gohome.Importer) {
}

func NewExtension() *extension {
	return &extension{}
}
