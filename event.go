package gohome

import "time"

type Event struct {
	Time           time.Time
	OriginalString string
	//TODO: Clarify what device this is
	Device        *Device
	ReplayCommand Command
}

func NewEvent(d *Device, cmd Command, orig string) Event {
	return Event{
		Time:           time.Now(),
		OriginalString: orig,
		Device:         d,
		ReplayCommand:  cmd,
	}
}

func (e *Event) String() string {
	//TODO: Time + Device
	return e.OriginalString + " : " + e.ReplayCommand.FriendlyString()
}
