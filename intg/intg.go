package intg

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/extensions/belkin"
	"github.com/markdaws/gohome/extensions/connectedbytcp"
	"github.com/markdaws/gohome/extensions/fluxwifi"
	"github.com/markdaws/gohome/extensions/lutron"
	"github.com/markdaws/gohome/log"
)

// Extension represents the interface any extension has to implement in order to
// be added to the system
type Extension interface {

	// RegisterCmdBuilders allows an extension to add cmd.Builder instances to the
	// map.  The key is the model number and the Builder then knows how to take abstract
	// commands like ZoneSetLevel and translate that to device specific commands for a
	// specific pieve of hardware
	RegisterCmdBuilders(*gohome.System, map[string]cmd.Builder)

	// RegisterDiscoverers allows extensions to register gohome.Discoverer instances for
	// model numbers.  Discoverers know how to scan the local network looking for a
	// particular kind of device
	RegisterDiscoverers(*gohome.System, map[string]gohome.Discoverer)

	// RegisterImports allows extensions to register importers for different model numbers.
	// An importer knows how to take device specific config files and convert that into
	// data types known by gohome
	RegisterImporters(*gohome.System, map[string]gohome.Importer)
}

// RegisterExtensions loads all of the know extensions into the specified system
func RegisterExtensions(sys *gohome.System) error {
	log.V("registering extensions")

	log.V("register extension - belkin")
	registerExtension(sys, belkin.NewExtension())

	log.V("register extension - connectedbytcp")
	registerExtension(sys, connectedbytcp.NewExtension())

	log.V("register extension - fluxwifi")
	registerExtension(sys, fluxwifi.NewExtension())

	log.V("register extension - lutron")
	registerExtension(sys, lutron.NewExtension())

	return nil
}

func registerExtension(sys *gohome.System, ext Extension) {
	ext.RegisterCmdBuilders(sys, sys.Extensions.CmdBuilders)
	ext.RegisterDiscoverers(sys, sys.Extensions.Discoverers)
	ext.RegisterImporters(sys, sys.Extensions.Importers)
}
