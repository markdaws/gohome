package gohome

type EventConsumer interface {
	EventConsumerID() string
	StartConsumingEvents() chan<- Event
}
