package event

//TODO: Delete - refactor recipes to use go-home-iot/event-bus

import (
	"time"

	"github.com/markdaws/gohome/cmd"
)

var nextID int

// Event represent an event that occurs in the system, such as a light being
// turned on or off.
type Event struct {
	ID             int
	Time           time.Time
	OriginalString string
	DeviceID       string
	ReplayCommand  cmd.Command
}

// New returns a new Event instance
func New(deviceID string, cmd cmd.Command, orig string) Event {
	nextID++

	return Event{
		ID:             nextID,
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
