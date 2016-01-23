package event

// Consumer is an interface for types that consume events from event.Broker.
// For example a Consumer might be a type that logs commands to the app log
// or a type that listens for events and sends them on to a websocket.
type Consumer interface {
	// EventConsumerID should return a unique ID which will be used
	// to register the consumer
	EventConsumerID() string

	// StartConsumingEvents when called signals the consumer that it will
	// start to reveice events.  It should return a channel that will be
	// used to send events to it
	StartConsumingEvents() chan<- Event
}
