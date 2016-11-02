package gohome

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/markdaws/gohome/cmd"
	"github.com/markdaws/gohome/log"
)

// CommandGroup contains a collection of commands that need to be run sequentially
type CommandGroup struct {
	Desc string
	Cmds []cmd.Command
}

// NewCommandGroup returns a CommandGroup instance with the Desc and Cmds field set
func NewCommandGroup(desc string, cmds ...cmd.Command) CommandGroup {
	return CommandGroup{Desc: desc, Cmds: cmds}
}

// CommandProcessor represents an interface to a type that knows how to process commands
type CommandProcessor interface {
	Start()
	Stop()
	Enqueue(CommandGroup) error
}

// CommandBuilder know how to take an abstract command like ZoneSetLevel and turn it
// in to a device specific set of instructions, for a specific piee of hardware
type CommandBuilder interface {
	Build(cmd.Command) (*cmd.Func, error)
}

// NewCommandProcessor returns an initialized type that implements the CommandProcessor interface
func NewCommandProcessor(system *System, maxWorkers, queueSize int) CommandProcessor {
	return &commandProcessor{
		system:     system,
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

func (cp *commandProcessor) Enqueue(cg CommandGroup) error {
	select {
	case cp.requests <- cg:
		log.V("CommandProcessor - enqueued: %s", cg.Desc)
		return nil
	default:
		err := errors.New("CommandProcessor - CommandGroup enqueue failed, CommandProcessor queue is full")
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

func (cp *commandProcessor) startWorker(index int) (errRet error) {

	log.V("CommandProcessor - starting worker %d", index)

	// If there is a panic for any reason trying to execute the commands
	// recover and log the error
	defer func() {
		if r := recover(); r != nil {
			errRet = fmt.Errorf("%s, %s", r, debug.Stack())
		}
	}()

	for cg := range cp.requests {
		log.V("CommandProcessor - execute group: %s", cg.Desc)

		cmds, err := cp.buildCommands(cg)
		if err != nil {
			log.E("CommandProcessor - unable to generate commands: %s, %s", cg.Desc, err)
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
			log.V("CommandProcessor - executed command: %s", c)
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
	var cmds []*cmd.Func

	for _, c := range cg.Cmds {
		finalCmd, err := cp.buildCommand(c)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, finalCmd...)
	}
	return cmds, nil
}

func (cp *commandProcessor) buildCommand(c cmd.Command) ([]*cmd.Func, error) {

	var cmds []*cmd.Func
	var finalCmd *cmd.Func
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

		hub := d.Hub
		if hub == nil {
			hub = d
		}

		var zCmd *cmd.Func
		var err error
		if hub.CmdBuilder != nil {
			zCmd, err = hub.CmdBuilder.Build(command)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("no command builder for device id:%s", z.DeviceID)
		}
		finalCmd = zCmd

	case *cmd.ZoneTurnOff:
		z, ok := cp.system.Zones[command.ZoneID]
		if !ok {
			return nil, fmt.Errorf("unknown zone ID %s", command.ZoneID)
		}
		d, ok := cp.system.Devices[z.DeviceID]
		if !ok {
			return nil, fmt.Errorf("unknown device ID %s", z.DeviceID)
		}

		hub := d.Hub
		if hub == nil {
			hub = d
		}

		var zCmd *cmd.Func
		var err error
		if hub.CmdBuilder != nil {
			zCmd, err = hub.CmdBuilder.Build(command)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("no command builder for device id:%s", z.DeviceID)
		}
		finalCmd = zCmd

	case *cmd.ZoneSetLevel:
		z, ok := cp.system.Zones[command.ZoneID]
		if !ok {
			return nil, fmt.Errorf("unknown zone ID %s", command.ZoneID)
		}
		d, ok := cp.system.Devices[z.DeviceID]
		if !ok {
			return nil, fmt.Errorf("unknown device ID %s", z.DeviceID)
		}

		hub := d.Hub
		if hub == nil {
			hub = d
		}

		var zCmd *cmd.Func
		var err error
		if hub.CmdBuilder != nil {
			zCmd, err = hub.CmdBuilder.Build(command)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("no command builder for device id:%s", z.DeviceID)
		}
		finalCmd = zCmd

	case *cmd.SceneSet:
		s, ok := cp.system.Scenes[command.SceneID]
		if !ok {
			return nil, fmt.Errorf("unknown scene ID %s", command.SceneID)
		}
		for _, sceneCmd := range s.Commands {
			// Scenes are a list of commands, so we may get multiple commands
			// that we need to execute, also scenes can execute other scenes
			sceneCmds, err := cp.buildCommand(sceneCmd)
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, sceneCmds...)
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
			hub = b.Device
		}

		var err error
		var bCmd *cmd.Func
		if hub.CmdBuilder != nil {
			bCmd, err = hub.CmdBuilder.Build(command)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("no command builder for device id:%s", hub.ID)
		}
		finalCmd = bCmd

	case *cmd.ButtonRelease:
		b, ok := cp.system.Buttons[command.ButtonID]
		if !ok {
			return nil, fmt.Errorf("unknown button ID %s", command.ButtonID)
		}
		hub := b.Device.Hub
		if hub == nil {
			hub = b.Device
		}

		var err error
		var bCmd *cmd.Func
		if hub.CmdBuilder != nil {
			bCmd, err = hub.CmdBuilder.Build(command)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("no command builder for device id:%s", hub.ID)
		}
		finalCmd = bCmd

	default:
		return nil, fmt.Errorf("unknown command, cannot process")
	}

	if finalCmd != nil {
		if finalCmd.Friendly == "" {
			finalCmd.Friendly = c.FriendlyString()
		}

		cmds = append(cmds, finalCmd)
	}

	return cmds, nil
}
