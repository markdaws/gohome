package event

import (
	"fmt"

	"github.com/markdaws/gohome/log"
)

// Broker is an interface for a type that implements a producer/consumer
// pattern. You can register producers who will produce events and then
// register consumers who will process the events. The job of the broker is
// to receive the events and then forward them to all of the consumers
type Broker interface {
	AddProducer(Producer)
	AddConsumer(Consumer)
	RemoveConsumer(Consumer)
	Init()
}

// NewBroker returns a type that implements the Broker interface.
func NewBroker() Broker {
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

func (b *broker) AddProducer(p Producer) {
	if !p.ProducesEvents() {
		return
	}

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

func (b *broker) AddConsumer(c Consumer) {
	ec := c.StartConsumingEvents()
	if ec == nil {
		return
	}

	log.V("%s adding consumer %s", b, c.EventConsumerID())
	b.consumers[c.EventConsumerID()] = ec
}

func (b *broker) RemoveConsumer(c Consumer) {
	id := c.EventConsumerID()
	eventChannel, ok := b.consumers[id]
	_ = eventChannel
	if !ok {
		return
	}

	//TODO: routine safe? need sync on map, verify
	log.V("%s removing consumer %s", b, c)
	close(eventChannel)
	delete(b.consumers, id)
}

func (b *broker) String() string {
	return "EventBroker"
}
