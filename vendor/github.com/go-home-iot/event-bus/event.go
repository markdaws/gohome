package evtbus

import "fmt"

// Event is an interface describing an event in the system
type Event interface {
	fmt.Stringer
}
