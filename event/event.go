package event

import (
	"time"

	"github.com/markdaws/gohome/cmd"
)

type Event struct {
	ID             int
	Time           time.Time
	OriginalString string
	DeviceID       string
	ReplayCommand  cmd.Command
}

var nextId int

func New(deviceID string, cmd cmd.Command, orig string) Event {
	nextId++

	return Event{
		ID:             nextId,
		Time:           time.Now(),
		OriginalString: orig,
		DeviceID:       deviceID,
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
