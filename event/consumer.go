package event

type Consumer interface {
	EventConsumerID() string
	StartConsumingEvents() chan<- Event
}
