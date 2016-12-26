package intg

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/extensions/belkin"
	"github.com/markdaws/gohome/extensions/connectedbytcp"
	"github.com/markdaws/gohome/extensions/fluxwifi"
	"github.com/markdaws/gohome/extensions/honeywell"
	"github.com/markdaws/gohome/extensions/lutron"
	"github.com/markdaws/gohome/extensions/testing"
	"github.com/markdaws/gohome/log"
)

// RegisterExtensions loads all of the know extensions into the specified system
func RegisterExtensions(sys *gohome.System) error {
	log.V("registering extensions")

	log.V("register extension - belkin")
	sys.Extensions.Register(belkin.NewExtension())

	log.V("register extension - connectedbytcp")
	sys.Extensions.Register(connectedbytcp.NewExtension())

	log.V("register extension - fluxwifi")
	sys.Extensions.Register(fluxwifi.NewExtension())

	log.V("register extension - honeywell")
	sys.Extensions.Register(honeywell.NewExtension())

	log.V("register extension - lutron")
	sys.Extensions.Register(lutron.NewExtension())

	/*
		// An example piece of hardware
		log.V("register extension - example")
		sys.Extensions.Register(example.NewExtension())
	*/

	//Uncomment for testing
	log.V("register extension - testing")
	sys.Extensions.Register(testing.NewExtension())

	return nil
}
