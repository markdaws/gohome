package gohome

import "github.com/markdaws/gohome/log"

type CommandProcessor interface {
	Process()
	Enqueue(Command) error
}

func NewCommandProcessor() CommandProcessor {
	return &commandProcessor{
		commands: make(chan Command, 10000),
	}
}

type commandProcessor struct {
	commands chan Command
}

func (cp *commandProcessor) Process() {
	//TODO: Have multiple workers?
	for c := range cp.commands {
		err := c.Execute()
		if err != nil {
			log.W("cmpProcesor:execute error:%s", err)
		} else {
			log.V("cmdProcessor:executed:%s", c)
		}
	}
}

func (cp *commandProcessor) Enqueue(c Command) error {
	log.V("cmdProcessor:enqueue:%s", c)
	cp.commands <- c
	return nil
}
