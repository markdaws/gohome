package gohome

import "fmt"

type EventBroker interface {
	AddProducer(EventProducer)
	AddConsumer(EventConsumer)
}

type EventProducer interface {
	StartProducingEvents() (<-chan Event, <-chan bool)
}

type EventConsumer interface {
	StartConsumingEvents() chan<- Event
}

func NewEventBroker() EventBroker {
	return &broker{}
}

type broker struct {
	consumers []chan<- Event
}

func (b *broker) AddProducer(p EventProducer) {
	ec, dc := p.StartProducingEvents()
	go func() {
		for {
			select {
			case e := <-ec:
				//Got a new event, process
				//fmt.Println("got an event:", e.String())

				//TODO: Async?
				for _, c := range b.consumers {
					c <- e
				}
			case <-dc:
				fmt.Println("Producer has stopped")
				return
			}
		}
	}()
}

func (b *broker) AddConsumer(c EventConsumer) {
	ec := c.StartConsumingEvents()
	if ec == nil {
		return
	}
	b.consumers = append(b.consumers, ec)
}
