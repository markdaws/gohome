package gohome

type EventProducer interface {
	ProducesEvents() bool
	StartProducingEvents() (<-chan Event, <-chan bool)
}
