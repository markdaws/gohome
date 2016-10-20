package gohome

import (
	"errors"
	"fmt"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

// CommandGroup contains a collection of commands that need to be run sequentially
type CommandGroup struct {
	Desc string
	Cmds []cmd.Command
}

func NewCommandGroup(desc string, cmds ...cmd.Command) CommandGroup {
	return CommandGroup{Desc: desc, Cmds: cmds}
}

type CommandProcessor interface {
	Start()
	Stop()
	Enqueue(CommandGroup) error
	SetSystem(s *System)
}

// CommandBuilder know how to take an abstract command like ZoneSetLevel and turn it
// in to a device specific set of instructions, for a specific piee of hardware
type CommandBuilder interface {
	Build(cmd.Command) (*cmd.Func, error)
}

func NewCommandProcessor(maxWorkers, queueSize int) CommandProcessor {
	return &commandProcessor{
		queueSize:  queueSize,
		maxWorkers: maxWorkers,
	}
}

type commandProcessor struct {
	maxWorkers int
	queueSize  int
	requests   chan CommandGroup
	system     *System
}

func (cp *commandProcessor) SetSystem(s *System) {
	cp.system = s
}

func (cp *commandProcessor) Enqueue(cg CommandGroup) error {
	select {
	case cp.requests <- cg:
		log.V("CommandGroup enqueued: %s", cg.Desc)
		return nil
	default:
		err := errors.New("CommandGroup enqueue failed, CommandProcessor queue is full")
		log.E(err.Error())
		return err
	}
}

func (cp *commandProcessor) Start() {
	log.V("CommandProcessor - starting")

	cp.requests = make(chan CommandGroup, cp.queueSize)

	for i := 0; i < cp.maxWorkers; i++ {
		i := i
		go func() {
			for {
				err := cp.startWorker(i)
				if err != nil {
					log.E("CommandProcessor - worker[%d] panic: %s", i, err)

					// restart the failed worker
				} else {
					//worker quit gracefully, allow shutdown, break out of the loop
					break
				}
			}
		}()
	}
}

func (cp *commandProcessor) startWorker(index int) (errRet interface{}) {

	log.V("CommandProcessor - starting worker %d", index)

	// If there is a panic for any reason trying to execute the commands
	// recover and log the error
	defer func() {
		if r := recover(); r != nil {
			errRet = r
		}
	}()

	for cg := range cp.requests {
		log.V("CommandProcessor - execute group: %s", cg.Desc)

		cmds, err := cp.buildCommands(cg)
		if err != nil {
			log.E("CommandProcessor - unable to generate commands: %s", cg.Desc)
			continue
		}

		for _, c := range cmds {
			log.V("CommandProcessor - executing command: %s", c)

			err := c.Func()
			if err != nil {
				log.E("CommandProcessor - execute error: %s", err)

				// Don't continue with any other commands in the group
				break
			}
		}
	}

	errRet = nil
	return
}

func (cp *commandProcessor) Stop() {
	log.V("CommandProcessor - stopping")
	close(cp.requests)

	//TODO: Wait until all workers return?
}

func (cp *commandProcessor) buildCommands(cg CommandGroup) ([]*cmd.Func, error) {

	cmds := make([]*cmd.Func, len(cg.Cmds))

	// TODO: Shouldn't have this code here, if each command contains the CmdBuild information
	// we can pull it directly vs having to get it back here.
	for i, c := range cg.Cmds {
		switch command := c.(type) {
		case *cmd.ZoneTurnOn:
			z, ok := cp.system.Zones[command.ZoneID]
			if !ok {
				return nil, fmt.Errorf("unknown zone ID %s", command.ZoneID)
			}
			d, ok := cp.system.Devices[z.DeviceID]
			if !ok {
				return nil, fmt.Errorf("unknown device ID %s", z.DeviceID)
			}

			var zCmd *cmd.Func
			var err error
			if d.CmdBuilder != nil {
				zCmd, err = d.CmdBuilder.Build(command)
			}
			if err != nil {
				return nil, err
			}
			cmds[i] = zCmd

		case *cmd.ZoneTurnOff:
			z, ok := cp.system.Zones[command.ZoneID]
			if !ok {
				return nil, fmt.Errorf("unknown zone ID %s", command.ZoneID)
			}
			d, ok := cp.system.Devices[z.DeviceID]
			if !ok {
				return nil, fmt.Errorf("unknown device ID %s", z.DeviceID)
			}

			var zCmd *cmd.Func
			var err error
			if d.CmdBuilder != nil {
				zCmd, err = d.CmdBuilder.Build(command)
			}
			if err != nil {
				return nil, err
			}
			cmds[i] = zCmd

		case *cmd.ZoneSetLevel:
			z, ok := cp.system.Zones[command.ZoneID]
			if !ok {
				return nil, fmt.Errorf("unknown zone ID %s", command.ZoneID)
			}
			d, ok := cp.system.Devices[z.DeviceID]
			if !ok {
				return nil, fmt.Errorf("unknown device ID %s", z.DeviceID)
			}

			var zCmd *cmd.Func
			var err error
			if d.CmdBuilder != nil {
				zCmd, err = d.CmdBuilder.Build(command)
			}
			if err != nil {
				return nil, err
			}
			cmds[i] = zCmd

		case *cmd.SceneSet:
			s, ok := cp.system.Scenes[command.SceneID]
			if !ok {
				return nil, fmt.Errorf("unknown scene ID %s", command.SceneID)
			}
			for _, sceneCmd := range s.Commands {
				err := cp.Enqueue(NewCommandGroup(cg.Desc, sceneCmd))
				if err != nil {
					return nil, err
				}
			}

		case *cmd.ButtonPress:
			b, ok := cp.system.Buttons[command.ButtonID]
			if !ok {
				return nil, fmt.Errorf("unknown button ID %s", command.ButtonID)
			}

			// The hub is the device that is used to talk to the target device. If the device
			// doesn't have a hub it is assumed to be a hub
			hub := b.Device.Hub
			if hub == nil {
				hub = &b.Device
			}

			//TODO: Remove hub or use it here
			var err error
			var bCmd *cmd.Func
			if b.Device.CmdBuilder != nil {
				bCmd, err = b.Device.CmdBuilder.Build(command)
			}
			if err != nil {
				return nil, err
			}
			cmds[i] = bCmd

		case *cmd.ButtonRelease:
			b, ok := cp.system.Buttons[command.ButtonID]
			if !ok {
				return nil, fmt.Errorf("unknown button ID %s", command.ButtonID)
			}
			hub := b.Device.Hub
			if hub == nil {
				hub = &b.Device
			}

			//TODO: remove hub or use it here
			var err error
			var bCmd *cmd.Func
			if b.Device.CmdBuilder != nil {
				bCmd, err = b.Device.CmdBuilder.Build(command)
			}
			if err != nil {
				return nil, err
			}
			cmds[i] = bCmd

		default:
			return nil, fmt.Errorf("unknown command, cannot process")
		}

		if cmds[i].Friendly == "" {
			cmds[i].Friendly = c.FriendlyString()
		}
	}
	return cmds, nil
}
