package lutron

import (
	"bufio"
	"errors"
	"fmt"
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
	lookupTable["l-bdgpro2-wh"] = &network{sys}
}

func (e *extension) RegisterImporters(sys *gohome.System, lookupTable map[string]gohome.Importer) {
	importer := &importer{System: sys}

	// Register a mapping from a moelnumber to an importer.  We can
	// register as many importers as we want here for multiple model numbers
	lookupTable["l-bdgpro2-wh"] = importer
}

func NewExtension() *extension {
	return &extension{}
}

type network struct {
	System *gohome.System
}

func (d *network) Devices(sys *gohome.System, modelNumber string) ([]*gohome.Device, error) {
	return nil, errors.New("unsupported method - Devices")
}

func (d *network) NewConnection(sys *gohome.System, dev *gohome.Device) (func(pool.Config) (net.Conn, error), error) {
	return func(cfg pool.Config) (net.Conn, error) {
		conn, err := net.Dial("tcp", dev.Address)
		if err != nil {
			return nil, err
		}

		r := bufio.NewReader(conn)
		_, err = r.ReadString(':')
		if err != nil {
			return nil, fmt.Errorf("authenticate login failed: %s", err)
		}

		_, err = conn.Write([]byte(dev.Auth.Login + "\r\n"))
		if err != nil {
			return nil, fmt.Errorf("authenticate write login failed: %s", err)
		}

		_, err = r.ReadString(':')
		if err != nil {
			return nil, fmt.Errorf("authenticate password failed: %s", err)
		}

		_, err = conn.Write([]byte(dev.Auth.Password + "\r\n"))
		if err != nil {
			return nil, fmt.Errorf("authenticate password failed: %s", err)
		}

		return conn, nil
	}, nil
}
