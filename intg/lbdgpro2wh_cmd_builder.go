package intg

import (
	"fmt"
	"strings"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

type lbdgpro2whCmdBuilder struct {
	System *gohome.System
}

func (b *lbdgpro2whCmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		return &cmd.Func{
			Func: func() error {
				newCmd := &cmd.StringCommand{
					Value: "#OUTPUT," + command.ZoneAddress + ",1,%.2f\r\n",
					Args:  []interface{}{command.Level.Value},
				}

				z := b.System.Zones[command.ZoneID]
				dev := b.System.Devices[z.DeviceID]
				return sendStringCommand(dev, newCmd)
			},
		}, nil
	case *cmd.ZoneTurnOn:
		return &cmd.Func{
			Func: func() error {
				newCmd := &cmd.StringCommand{
					Value: "#OUTPUT," + command.ZoneAddress + ",1,%.2f\r\n",
					Args:  []interface{}{100.0},
				}
				z := b.System.Zones[command.ZoneID]
				dev := b.System.Devices[z.DeviceID]
				return sendStringCommand(dev, newCmd)
			},
		}, nil
	case *cmd.ZoneTurnOff:
		return &cmd.Func{
			Func: func() error {
				newCmd := &cmd.StringCommand{
					Value: "#OUTPUT," + command.ZoneAddress + ",1,%.2f\r\n",
					Args:  []interface{}{0.0},
				}
				z := b.System.Zones[command.ZoneID]
				dev := b.System.Devices[z.DeviceID]
				return sendStringCommand(dev, newCmd)
			},
		}, nil
	case *cmd.ButtonPress:
		return &cmd.Func{
			Func: func() error {
				newCmd := &cmd.StringCommand{
					//TODO: This device also has a local id of 1 which we have to send as well as
					//an actual IP address... fix
					Value: "#DEVICE," + command.DeviceAddress + "," + command.ButtonAddress + ",3\r\n",
				}
				dev := b.System.Devices[command.DeviceID]
				return sendStringCommand(dev, newCmd)
			},
		}, nil
	case *cmd.ButtonRelease:
		return &cmd.Func{
			Func: func() error {
				newCmd := &cmd.StringCommand{
					Value: "#DEVICE," + command.DeviceAddress + "," + command.ButtonAddress + ",4\r\n",
				}
				dev := b.System.Devices[command.DeviceID]
				return sendStringCommand(dev, newCmd)
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
		return fmt.Errorf("Failed to string string_command %s\n", err)
	}
	return nil
}
