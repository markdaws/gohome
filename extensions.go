package gohome

import (
	"github.com/markdaws/gohome/cmd"
)

type Discoverer interface {
	Devices(sys *System, modelNumber string) ([]*Device, error)
}

type Importer interface {
	FromString(sys *System, data string, modelNumber string) error
}

type Extensions struct {
	CmdBuilders map[string]cmd.Builder
	Discoverers map[string]Discoverer
	Importers   map[string]Importer
}

func NewExtensions() *Extensions {
	exts := &Extensions{}
	exts.CmdBuilders = make(map[string]cmd.Builder)
	exts.Discoverers = make(map[string]Discoverer)
	exts.Importers = make(map[string]Importer)

	return exts
}
