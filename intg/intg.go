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

type Extension interface {
	RegisterCmdBuilders(*gohome.System, map[string]cmd.Builder)
	RegisterDiscoverers(*gohome.System, map[string]gohome.Discoverer)

	//TODO: Parse config files into system
	//TODO: Importing devices
	//TODO: Import UI, have webpack get necessary files
}

func RegisterExtensions(sys *gohome.System) error {
	log.V("registering extensions")

	//TODO: look into maybe using "go generate" to somehow build this code dynamically

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
}
