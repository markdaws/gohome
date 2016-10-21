package belkin

import (
	"net"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type extension struct {
}

func (e *extension) RegisterCmdBuilders(sys *gohome.System, lookupTable map[string]cmd.Builder) {
	//Belkin WeMo Insight
	lookupTable["f7c029v2"] = &cmdBuilder{System: sys, id: "f7c029v2"}
	//Belkin WeMo Maker
	lookupTable["f7c043fc"] = &cmdBuilder{System: sys, id: "f7c043fc"}
}

func (e *extension) RegisterNetwork(sys *gohome.System, lookupTable map[string]gohome.Network) {
	//Belkin WeMo Insight
	lookupTable["f7c029v2"] = &network{System: sys}
	//Belkin WeMo Maker
	lookupTable["f7c043fc"] = &network{System: sys}
}

func (e *extension) RegisterImporters(sys *gohome.System, lookupTable map[string]gohome.Importer) {
}

func (e *extension) RegisterConnFactories(
	sys *gohome.System,
	lookupTable map[string]func(*gohome.Device, pool.Config) (net.Conn, error)) {
}

func NewExtension() *extension {
	return &extension{}
}
