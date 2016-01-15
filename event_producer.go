package gohome

type EventProducer interface {
	StartProducingEvents() (<-chan Event, <-chan bool)
}
