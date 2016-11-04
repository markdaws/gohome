package connectedbytcp

import (
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type extension struct {
	gohome.NullExtension
}

func (e *extension) BuilderForDevice(sys *gohome.System, d *gohome.Device) cmd.Builder {
	switch d.ModelNumber {
	case "tcp600gwb":
		return &cmdBuilder{System: sys}
	default:
		return nil
	}
}

func (e *extension) NetworkForDevice(sys *gohome.System, d *gohome.Device) gohome.Network {
	switch d.ModelNumber {
	case "tcp600gwb":
		return &network{}
	default:
		return nil
	}
}

//TODO: Events

func (e *extension) Discovery(sys *gohome.System) gohome.Discovery {
	return &discovery{System: sys}
}

func (e *extension) Name() string {
	return "connectedbytcp"
}

func NewExtension() *extension {
	return &extension{}
}
