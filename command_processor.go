package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

type CommandProcessor interface {
	Process()
	Enqueue(cmd.Command) error
	SetSystem(s *System)
}

type CommandBuilder interface {
	Build(cmd.Command) (*cmd.Func, error)
}

func NewCommandProcessor() CommandProcessor {
	return &commandProcessor{
		commands: make(chan *cmd.Func, 10000),
	}
}

type commandProcessor struct {
	commands chan *cmd.Func
	system   *System
}

func (cp *commandProcessor) SetSystem(s *System) {
	cp.system = s
}

func (cp *commandProcessor) Process() {
	//TODO: Have multiple workers?
	for c := range cp.commands {
		err := c.Func()
		if err != nil {
			log.E("cmdProcessor:execute error:%s", err)
		} else {
			log.V("cmdProcessor:executed:%s", c)
		}
	}
}

func (cp *commandProcessor) Enqueue(c cmd.Command) error {
	log.V("cmdProcessor:enqueue:%s", c)

	//TODO: use devicer (defined in this namespace), remove reference to system, move into cmd package
	switch command := c.(type) {
	case *cmd.ZoneSetLevel:
		z, ok := cp.system.Zones[command.ZoneID]
		if !ok {
			return fmt.Errorf("unknown zone ID %s", command.ZoneID)
		}
		zCmd, err := cp.system.Devices[z.DeviceID].BuildCommand(command)
		if err != nil {
			return err
		}
		cp.commands <- zCmd

	case *cmd.ZoneTurnOn:
		z, ok := cp.system.Zones[command.ZoneID]
		if !ok {
			return fmt.Errorf("unknown zone ID %s", command.ZoneID)
		}
		zCmd, err := cp.system.Devices[z.DeviceID].BuildCommand(command)
		if err != nil {
			return err
		}
		cp.commands <- zCmd

	case *cmd.ZoneTurnOff:
		z, ok := cp.system.Zones[command.ZoneID]
		if !ok {
			return fmt.Errorf("unknown zone ID %s", command.ZoneID)
		}
		zCmd, err := cp.system.Devices[z.DeviceID].BuildCommand(command)
		if err != nil {
			return err
		}
		cp.commands <- zCmd

	case *cmd.SceneSet:
		s, ok := cp.system.Scenes[command.SceneID]
		if !ok {
			return fmt.Errorf("unknown scene ID %s", command.SceneID)
		}
		for _, sceneCmd := range s.Commands {
			err := cp.Enqueue(sceneCmd)
			if err != nil {
				return err
			}
		}

	case *cmd.ButtonPress:
		b, ok := cp.system.Buttons[command.ButtonID]
		if !ok {
			return fmt.Errorf("unknown button ID %s", command.ButtonID)
		}

		// The hub is the device that is used to talk to the target device. If the device
		// doesn't have a hub it is assumed to be a hub
		hub := b.Device.Hub()
		if hub == nil {
			hub = b.Device
		}
		bCmd, err := hub.BuildCommand(command)
		if err != nil {
			return err
		}
		cp.commands <- bCmd

	case *cmd.ButtonRelease:
		b, ok := cp.system.Buttons[command.ButtonID]
		if !ok {
			return fmt.Errorf("unknown button ID %s", command.ButtonID)
		}
		hub := b.Device.Hub()
		if hub == nil {
			hub = b.Device
		}
		bCmd, err := hub.BuildCommand(command)
		if err != nil {
			return err
		}
		cp.commands <- bCmd

	default:
		return fmt.Errorf("unknown command, cannot process")
	}
	return nil
}
