package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/log"
)

type CommandProcessor interface {
	Process()
	Enqueue(Command) error
}

func NewCommandProcessor() CommandProcessor {
	return &commandProcessor{
		commands: make(chan *FCommand, 10000),
	}
}

type commandProcessor struct {
	commands chan *FCommand
}

func (cp *commandProcessor) Process() {
	//TODO: Have multiple workers?
	for c := range cp.commands {
		err := c.Func()
		if err != nil {
			log.W("cmpProcesor:execute error:%s", err)
		} else {
			log.V("cmdProcessor:executed:%s", c)
		}
	}
}

func (cp *commandProcessor) Enqueue(c Command) error {
	log.V("cmdProcessor:enqueue:%s", c)

	switch cmd := c.(type) {
	case *ZoneSetLevelCommand:
		zCmd, err := cmd.Zone.Device.BuildCommand(cmd)
		if err != nil {
			return err
		}
		cp.commands <- zCmd

	case *SceneSetCommand:
		for _, sceneCmd := range cmd.Scene.Commands {
			err := cp.Enqueue(sceneCmd)
			if err != nil {
				return err
			}
		}

	case *ButtonPressCommand:
		bCmd, err := cmd.Button.Device.BuildCommand(cmd)
		if err != nil {
			return err
		}
		cp.commands <- bCmd

	case *ButtonReleaseCommand:
		bCmd, err := cmd.Button.Device.BuildCommand(cmd)
		if err != nil {
			return err
		}
		cp.commands <- bCmd

	default:
		return fmt.Errorf("unknown command, cannot process")
	}
	return nil
}
