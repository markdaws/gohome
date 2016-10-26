package gohome

import (
	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/log"
)

// LogConsumer consumes events from the event bus and prints them to the log
// TODO: Log as json objects so we can replay everything that happens in the system
type LogConsumer struct{}

func (c *LogConsumer) ConsumerName() string {
	return "LogConsumer"
}

func (c *LogConsumer) StartConsuming(ch chan evtbus.Event) {
	log.V("LogConsumer - start consuming events")

	go func() {
		for e := range ch {
			log.V("event: %s", e.String())
		}
		log.V("LogConsumer - event channel has closed")
	}()
}

func (c *LogConsumer) StopConsuming() {
	log.V("LogConsumer - stop consuming events")
	//TODO:
}
