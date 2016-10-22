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
/*func (c *TelnetConnection) Read(p []byte) (n int, err error) {
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	n, err = c.conn.Read(p)
	if err != nil {
		c.status = CSClosed
	}
	return
}
func (c *TelnetConnection) Write(p []byte) (n int, err error) {
	c.conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	n, err = c.conn.Write(p)
	if err != nil {
		c.status = CSClosed
	}
	return
}*/

// Network is an interface that must be exported by an extension that provides
// network related functionality pertaining to the extensions hardware
type Network interface {
	Devices(sys *System, modelNumber string) ([]*Device, error)
	NewConnection(sys *System, d *Device) (func(pool.Config) (net.Conn, error), error)
}

// Importer is an interface that the extension provides to generate gohome models from
// 3rd party config files
type Importer interface {
	FromString(sys *System, data string, modelNumber string) error
}

// Extension represents the interface any extension has to implement in order to
// be added to the system
type Extension interface {

	// BuilderForDevice should return a cmd.Builder if the extension exports a builder
	// for the device that was passed in to the function, nil otherwise
	BuilderForDevice(*System, *Device) cmd.Builder

	// NetworkForDevice should return a gohome.Network if the extension exports a Network interface
	// for the device that was passed in to the function, nil otherwise
	NetworkForDevice(*System, *Device) Network

	// ImporterForDevice should return a gohome.Importer if the extension exports an Importer
	// for the device that was passed in to the function, nil otherwise
	ImporterForDevice(*System, *Device) Importer

	// Name returns a friendly name for the extension
	Name() string
}

// Extensions contains references to all of the loaded extensions in a system
type Extensions struct {
	extensions []Extension
}

// Register adds a new extension to the Extensions instance
func (e *Extensions) Register(ext Extension) {
	e.extensions = append(e.extensions, ext)
}

// FindCmdBuilder returns a cmd.Builder instance if there is any extension that
// exports one for the device passed in to the function
func (e *Extensions) FindCmdBuilder(sys *System, d *Device) cmd.Builder {
	for _, ext := range e.extensions {
		builder := ext.BuilderForDevice(sys, d)
		if builder != nil {
			return builder
		}
	}
	return nil
}

// FindNetwork returns a gohome.Network instance if there is any extension that
// exports one for the device passed in to the function
func (e *Extensions) FindNetwork(sys *System, d *Device) Network {
	for _, ext := range e.extensions {
		network := ext.NetworkForDevice(sys, d)
		if network != nil {
			return network
		}
	}
	return nil
}

// FindImporter returns a gohome.Importer instance if there is any extension that
// exports one for the device passed in to the function
func (e *Extensions) FindImporter(sys *System, d *Device) Importer {
	for _, ext := range e.extensions {
		importer := ext.ImporterForDevice(sys, d)
		if importer != nil {
			return importer
		}
	}
	return nil
}

// NewExtensions inits and returns a new Extensions instance
func NewExtensions() *Extensions {
	exts := &Extensions{}
	return exts
}
