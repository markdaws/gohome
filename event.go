package gohome

import "time"

type Event struct {
	Time           time.Time
	OriginalString string
	Device         *Device
	ReplayCommand  Command
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
	out := e.OriginalString
	if e.ReplayCommand != nil {
		out += e.ReplayCommand.FriendlyString()
	}
	return out
}
