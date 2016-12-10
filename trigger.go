package gohome

import "github.com/go-home-iot/event-bus"

// Trigger is an interface representing an automation trigger
type Trigger interface {
	evtbus.Consumer
	Trigger()
}
