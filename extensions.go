package gohome

import (
	"net"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome/cmd"
)

//TODO: Ping mechanism
//TODO: Check connection is bad don't put back in the pool
//TODO: Set write, read timeouts for connections
//TODO: Store retry time in system config file
type Network interface {
	Devices(sys *System, modelNumber string) ([]*Device, error)
	NewConnection(sys *System, d *Device) (func(pool.Config) (net.Conn, error), error)
}

type Importer interface {
	FromString(sys *System, data string, modelNumber string) error
}

type Extensions struct {
	CmdBuilders map[string]cmd.Builder
	Network     map[string]Network
	Importers   map[string]Importer
}

//TODO: Store as linked list, iterate through, you pass in device, keep going until
//one responds that it handles this device

func NewExtensions() *Extensions {
	exts := &Extensions{}
	exts.CmdBuilders = make(map[string]cmd.Builder)
	exts.Network = make(map[string]Network)
	exts.Importers = make(map[string]Importer)

	return exts
}
