package evtbus

// Producer is an interface that describes a type that produces events on the event bus
type Producer interface {
	// Name should return a friendly name for the consumer.  This field
	// may be useful for debugging
	ProducerName() string

	// StartProducing is called when the procuder is added to the bus
	StartProducing(*Bus)

	// StopProducing is called when the producer is removed from the bus
	StopProducing()
}
