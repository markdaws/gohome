package connectedbytcp

import (
	"net"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type extension struct {
}

func (e *extension) RegisterCmdBuilders(sys *gohome.System, lookupTable map[string]cmd.Builder) {
	builder := &cmdBuilder{System: sys}
	lookupTable[builder.ID()] = builder
}

func (e *extension) RegisterNetwork(sys *gohome.System, lookupTable map[string]gohome.Network) {
	lookupTable["tcp600gwb"] = &network{System: sys}
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
