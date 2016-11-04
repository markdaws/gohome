package fluxwifi

import (
	"net"
	"time"

	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
)

type network struct{}

func (d *network) NewConnection(sys *gohome.System, dev *gohome.Device) (func(pool.Config) (net.Conn, error), error) {
	return func(cfg pool.Config) (net.Conn, error) {
		return net.DialTimeout("tcp", dev.Address, time.Second*10)
	}, nil
}
