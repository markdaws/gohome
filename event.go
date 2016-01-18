package gohome

import "time"

type Event struct {
	ID             int
	Time           time.Time
	OriginalString string
	Device         Device
	ReplayCommand  Command
	Type           EventType
}

var nextId int = 0

//TODO: Needed or just use command type?
type EventType uint32

const (
	//TODO: Define event types
	ETUnknown = iota

	// Ping event to a device
	ETPing
)

func NewEvent(d Device, cmd Command, orig string, t EventType) Event {
	nextId++

	return Event{
		ID:             nextId,
		Time:           time.Now(),
		OriginalString: orig,
		Device:         d,
		ReplayCommand:  cmd,
		Type:           t,
	}
}

func (e *Event) String() string {
	out := ""
	if e.ReplayCommand != nil {
		out += e.ReplayCommand.FriendlyString()
	}
	return out
}
