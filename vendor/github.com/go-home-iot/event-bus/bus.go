package evtbus

import (
	"errors"
	"sync"
)

// Bus is a type that facilitates sending events from multiple producers to multiple consumers. It
// is non blocking, so fast producers won't break the bus and slow consumers won't block other consumers
// from receiving events
type Bus struct {
	consumers map[Consumer]chan Event
	producers []Producer
	events    chan Event
	stopped   bool
	mutex     sync.RWMutex

	// Capacity is the number of events that can be in the bus at one time before incoming
	// events will be ignored
	Capacity int

	// ConsumerCapacity is the number of events a consumer can queue. If the consumer is slow
	// at processing events, you should set this to some high number. Once this number is reached
	// the bus will throw away events to the slow consumer so that other consumers are not blocked
	ConsumerCapacity int
}

// ErrBusFull indicates the bus is full and can't handle any more events
// at this time.
var ErrBusFull = errors.New("Event bus is full, no more events can be added at this time")

// NewBus returns an initialized bus
func NewBus(capacity, consumerCapacity int) *Bus {
	b := &Bus{Capacity: capacity, ConsumerCapacity: consumerCapacity}
	b.consumers = make(map[Consumer]chan Event)
	b.events = make(chan Event, capacity)
	b.init()
	return b
}

func (b *Bus) init() {
	go func() {
		for {
			// Wait for events to process
			e, more := <-b.events
			if b.stopped || !more {
				return
			}

			// Each consumer has it's own buffered channel, that way we can send an event
			// to the consumer and not block the bus, so we can send to other consumers if
			// others are slow and at full capacity
			b.mutex.RLock()
			for _, q := range b.consumers {
				select {
				case q <- e:
				default:
					// Consumer queue was full, drop event, keep going for other consumers
				}
			}
			b.mutex.RUnlock()
		}
	}()
}

// Stop removes all of the consumers and producers and stops processing events. After calling
// this method the Bus is no longer usable and you should create a new one of you need another bus
func (b *Bus) Stop() {
	b.stopped = true

	for len(b.consumers) > 0 {
		// Keyed with the consumer, so get first key each time then break
		for c := range b.consumers {
			b.RemoveConsumer(c)
			break
		}
	}
	for len(b.producers) > 0 {
		b.RemoveProducer(b.producers[0])
	}
	close(b.events)
}

// AddConsumer adds a consumer to the bus, once added the consumer will start to
// receive events from the bus
func (b *Bus) AddConsumer(c Consumer) {
	if b.stopped {
		return
	}

	_, ok := b.consumers[c]
	if ok {
		// Already consuming, ignore
		return
	}

	b.mutex.Lock()
	b.consumers[c] = make(chan Event, b.ConsumerCapacity)
	b.mutex.Unlock()

	c.StartConsuming(b.consumers[c])
}

// RemoveConsumer removes a consumer from the bus, once removed consumers will no
// longer receive events
func (b *Bus) RemoveConsumer(c Consumer) {
	b.mutex.RLock()
	q, ok := b.consumers[c]
	b.mutex.RUnlock()

	if !ok {
		return
	}
	delete(b.consumers, c)
	close(q)
	c.StopConsuming()
}

// AddProducer adds a producer to the bus and registers it
func (b *Bus) AddProducer(p Producer) {
	b.mutex.Lock()

	for _, prod := range b.producers {
		// Already producing, ignore
		if prod == p {
			return
		}
	}
	b.producers = append(b.producers, p)
	b.mutex.Unlock()

	p.StartProducing(b)
}

// RemoveProducer removes a producer from the bus
func (b *Bus) RemoveProducer(p Producer) {
	for i, prod := range b.producers {
		if prod == p {
			b.mutex.Lock()
			b.producers = append(b.producers[:i], b.producers[i+1:]...)
			b.mutex.Unlock()

			p.StopProducing()
			return
		}
	}
}

// Enqueue adds an event to the event bus. It is non blocking, if there is not
// enough capacity in the bus to add a new event, the method returns an error
func (b *Bus) Enqueue(e Event) error {
	if b.stopped {
		return nil
	}

	select {
	case b.events <- e:
		return nil
	default:
		return ErrBusFull
	}
	return nil
}
