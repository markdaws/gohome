package gohome

import "time"

type Event struct {
	ID             int
	Time           time.Time
	OriginalString string
	Device         *Device
	ReplayCommand  Command
}

var nextId int = 0

func NewEvent(d *Device, cmd Command, orig string) Event {
	nextId++

	return Event{
		ID:             nextId,
		Time:           time.Now(),
		OriginalString: orig,
		Device:         d,
		ReplayCommand:  cmd,
	}
}

func (e *Event) String() string {
	out := ""
	if e.ReplayCommand != nil {
		out += e.ReplayCommand.FriendlyString()
	}
	return out
}
