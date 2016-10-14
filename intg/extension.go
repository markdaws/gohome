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

	//TODO: Parse config files into system
	//TODO: Importing devices
	//TODO: Import UI, have webpack get necessary files
}

func RegisterExtensions(sys *gohome.System) map[string]cmd.Builder {
	log.V("registering extensions")

	builders := make(map[string]cmd.Builder)

	//TODO: look into maybe using "go generate" to somehow build this code dynamically

	log.V("register extension - belkin")
	registerExtension(sys, belkin.NewExtension(), builders)

	log.V("register extension - connectedbytcp")
	registerExtension(sys, connectedbytcp.NewExtension(), builders)

	log.V("register extension - fluxwifi")
	registerExtension(sys, fluxwifi.NewExtension(), builders)

	log.V("register extension - lutron")
	registerExtension(sys, lutron.NewExtension(), builders)

	return builders
}

func registerExtension(sys *gohome.System, ext Extension, builders map[string]cmd.Builder) {
	ext.RegisterCmdBuilders(sys, builders)
}
