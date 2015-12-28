package gohome

import (
	"time"

	"github.com/nu7hatch/gouuid"
)

// At a certain time e.g. 8pm
// time no year, no month, no day, hour, minute, second
// After a certain delay every 5 minutes
// Iterations - certain number of times
// TODO: Be able to get sunrise/sunset time for a location: https://github.com/cpucycle/astrotime
// Days of week - e.g. Tues/Wed/Sun
type TimeTrigger struct {
	Iterations uint64
	Forever    bool
	At         time.Time
	Interval   time.Duration

	timer    *time.Timer
	ticker   *time.Ticker
	doneChan chan bool
	id       string
	enabled  bool
}

func (t *TimeTrigger) Type() string {
	return "gohome.TimeTrigger"
}

func (t *TimeTrigger) Name() string {
	return "Time Trigger"
}

func (t *TimeTrigger) Description() string {
	return "Triggers when the specified time or duration expires"
}

func (t *TimeTrigger) Enabled() bool {
	return t.enabled
}

func (t *TimeTrigger) SetEnabled(enabled bool) {
	t.enabled = enabled
}

//TODO: Create via reflection
func (t *TimeTrigger) New() Trigger {
	return &TimeTrigger{}
}

func (t *TimeTrigger) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			Identifiable: Identifiable{
				ID:          "Iterations",
				Name:        "Iterations",
				Description: "The number of times the trigger will fire before stopping",
			},
			Type: "integer",
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "Forever",
				Name:        "Forever",
				Description: "If true, the trigger will run forever",
			},
			Type: "boolean",
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "At",
				Name:        "At",
				Description: "The date and time to fire the trigger",
			},
			Type: "datetime",
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "Interval",
				Name:        "Interval",
				Description: "The time (in ms) between each trigger event",
			},
			Type: "integer",
		},
	}
}

func (t *TimeTrigger) EventConsumerID() string {
	if t.id == "" {
		id, err := uuid.NewV4()
		if err != nil {
			//TODO: error
		}
		t.id = id.String()
	}
	return t.id
}

func (t *TimeTrigger) Start() (<-chan bool, <-chan bool) {
	fireChan := make(chan bool)
	t.doneChan = make(chan bool)

	go func() {
		if !t.At.IsZero() {
			var count uint64 = 0
			finalAt := t.At
			for {
				t.timer = time.NewTimer(finalAt.Sub(time.Now()))
				<-t.timer.C
				fireChan <- true
				count++
				if !t.Forever && count >= t.Iterations {
					break
				}
				finalAt = t.At.Add(t.Interval * time.Duration(count))
			}
		} else if t.Interval != 0 {
			t.ticker = time.NewTicker(t.Interval)
			var count uint64 = 0
			for _ = range t.ticker.C {
				fireChan <- true
				count++
				if !t.Forever && count >= t.Iterations {
					break
				}
			}
		}
		doneChan := t.doneChan
		t.doneChan = nil
		close(doneChan)
	}()
	return fireChan, t.doneChan
}

func (t *TimeTrigger) Stop() {
	if t.timer != nil {
		t.timer.Stop()
	}
	if t.ticker != nil {
		t.ticker.Stop()
	}
	if t.doneChan != nil {
		close(t.doneChan)
	}
}
