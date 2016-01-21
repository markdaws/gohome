package event

type Producer interface {
	ProducesEvents() bool
	StartProducingEvents() (<-chan Event, <-chan bool)
}
