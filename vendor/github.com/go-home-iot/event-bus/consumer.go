package evtbus

// Consumer is an interface that describes a type that consumes events
type Consumer interface {
	// Name should return a friendly name for the consumer.  This field
	// may be useful for debugging
	ConsumerName() string

	// StartConsuming is called when the consumer is added to the event bus.
	// The channel which is passed in will recieve all of the events in the bus
	StartConsuming(chan Event)

	// StopConsuming is called when the consumer is removed from the event bus
	StopConsuming()
}
