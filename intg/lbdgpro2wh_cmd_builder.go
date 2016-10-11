package intg

import (
	"fmt"
	"io"
	"strings"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/lutron"
)

type lbdgpro2whCmdBuilder struct {
	System *gohome.System
	device lutron.Device
}

func (b *lbdgpro2whCmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {

	if b.device == nil {
		lDev, err := lutron.DeviceFromModelNumber(b.ID())
		if err != nil {
			return nil, err
		}
		b.device = lDev
	}

	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		return &cmd.Func{
			Func: func() error {
				z := b.System.Zones[command.ZoneID]
				dev := b.System.Devices[z.DeviceID]

				return getWriterAndExec(dev, func(w io.Writer) error {
					return b.device.SetLevel(command.Level.Value, command.ZoneAddress, w)
				})
			},
		}, nil
	case *cmd.ZoneTurnOn:
		return &cmd.Func{
			Func: func() error {
				z := b.System.Zones[command.ZoneID]
				dev := b.System.Devices[z.DeviceID]
				return getWriterAndExec(dev, func(w io.Writer) error {
					return b.device.SetLevel(100.0, command.ZoneAddress, w)
				})
			},
		}, nil
	case *cmd.ZoneTurnOff:
		return &cmd.Func{
			Func: func() error {
				z := b.System.Zones[command.ZoneID]
				dev := b.System.Devices[z.DeviceID]
				return getWriterAndExec(dev, func(w io.Writer) error {
					return b.device.SetLevel(0.0, command.ZoneAddress, w)
				})
			},
		}, nil
	case *cmd.ButtonPress:
		return &cmd.Func{
			Func: func() error {
				dev := b.System.Devices[command.DeviceID]
				return getWriterAndExec(dev, func(w io.Writer) error {
					return b.device.ButtonPress(command.DeviceAddress, command.ButtonAddress, w)
				})
			},
		}, nil
	case *cmd.ButtonRelease:
		return &cmd.Func{
			Func: func() error {
				dev := b.System.Devices[command.DeviceID]
				return getWriterAndExec(dev, func(w io.Writer) error {
					return b.device.ButtonPress(command.DeviceAddress, command.ButtonAddress, w)
				})
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported command type")
	}
	return nil, nil
}

func (b *lbdgpro2whCmdBuilder) ID() string {
	return "l-bdgpro2-wh"
}

func getWriterAndExec(d gohome.Device, f func(io.Writer) error) error {
	conn := d.Connections.Get()
	if conn == nil {
		return fmt.Errorf("error connecting, pool returned nil")
	}

	defer func() {
		d.Connections.Release(conn)
	}()

	err := f(conn)
	if err != nil {
		return fmt.Errorf("Failed to send command %s\n", err)
	}
	return nil
}

func sendStringCommand(d gohome.Device, cmd *cmd.StringCommand) error {

	log.V(
		"sending command \"%s\"",
		strings.Replace(strings.Replace(cmd.String(), "\r", "\\r", -1), "\n", "\\n", -1),
	)

	conn := d.Connections.Get()
	if conn == nil {
		return fmt.Errorf("StringCommand - error connecting, pool returned nil")
	}

	defer func() {
		d.Connections.Release(conn)
	}()

	_, err := conn.Write([]byte(cmd.String()))
	if err != nil {
		return fmt.Errorf("Failed to send string_command %s\n", err)
	}
	return nil
}
