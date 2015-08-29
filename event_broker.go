package gohome

import "fmt"

type EventBroker interface {
	AddProducer(EventProducer)
}

type EventProducer interface {
	GetEventProducerChans() (<-chan Event, <-chan bool)
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
