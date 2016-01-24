package event

// Producer represents the interface for types that are capable of
// producing events
type Producer interface {

	// ProducesEvents returns true if the type will produce any events
	ProducesEvents() bool

	// StartProducingEvents will be called at some point and signals that
	// the caller is ready for the device to produce events. The function
	// returns one Channel that can be used to receive events and another
	// channel that the producer can use to signal that it won't produce
	// any more events
	StartProducingEvents() (<-chan Event, <-chan bool)
}
