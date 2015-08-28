package gohome

import (
	"fmt"
	"time"
)

type EventProducer interface {
	GetEventProducerChans() (<-chan Event, <-chan bool)
}

type Event struct {
	Time        time.Time
	StringValue string
	Device      *Device
}

func (e *Event) String() string {
	return e.StringValue
}

type EventBroker interface {
	AddProducer(EventProducer)
}

func NewEventBroker() EventBroker {
	return &broker{}
}

type broker struct {
}

func (e *broker) AddProducer(p EventProducer) {
	ec, dc := p.GetEventProducerChans()
	go func() {
		for {
			select {
			case e := <-ec:
				//Got a new event, process
				fmt.Println("got an event:", e.String())
			case <-dc:
				fmt.Println("Producer has stopped")
				return
			}
		}
	}()
}
