package gohome

import "fmt"

type EventBroker interface {
	AddProducer(EventProducer)
	AddConsumer(EventConsumer)
	RemoveConsumer(EventConsumer)
	Init()
}

type EventProducer interface {
	StartProducingEvents() (<-chan Event, <-chan bool)
}

type EventConsumer interface {
	EventConsumerID() string
	StartConsumingEvents() chan<- Event
}

func NewEventBroker() EventBroker {
	return &broker{
		consumers: make(map[string]chan<- Event),
	}
}

type broker struct {
	consumers  map[string]chan<- Event
	eventQueue chan Event
}

func (b *broker) Init() {
	b.eventQueue = make(chan Event, 10000)

	// Want to process the events serially incase the order is important
	// for triggers vs. processing many events in parallel
	go func() {
		for {
			select {
			case e := <-b.eventQueue:
				for _, c := range b.consumers {
					c <- e
				}
			}
		}
	}()
}

func (b *broker) AddProducer(p EventProducer) {
	ec, dc := p.StartProducingEvents()
	go func() {
		for {
			select {
			case e := <-ec:
				b.eventQueue <- e
			case <-dc:
				//TODO:
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

	//fmt.Printf("Adding consumer: %s\n", c.EventConsumerID())
	b.consumers[c.EventConsumerID()] = ec
}

func (b *broker) RemoveConsumer(c EventConsumer) {
	id := c.EventConsumerID()
	eventChannel, ok := b.consumers[id]
	_ = eventChannel
	if !ok {
		return
	}

	delete(b.consumers, id)
}
