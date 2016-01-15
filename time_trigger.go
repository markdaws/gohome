package gohome

import "time"

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

	timer  *time.Timer
	ticker *time.Ticker
	id     string
	fire   chan bool
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

func (t *TimeTrigger) New() Trigger {
	return &TimeTrigger{}
}

func (t *TimeTrigger) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			ID:          "Iterations",
			Name:        "Iterations",
			Description: "The number of times the trigger will fire before stopping",
			Type:        "integer",
		},
		Ingredient{
			ID:          "Forever",
			Name:        "Forever",
			Description: "If true, the trigger will run forever",
			Type:        "boolean",
		},
		Ingredient{
			ID:          "At",
			Name:        "At",
			Description: "The date and time to fire the trigger",
			Type:        "datetime",
		},
		Ingredient{
			ID:          "Interval",
			Name:        "Interval",
			Description: "The time (in ms) between each trigger event",
			Type:        "integer",
		},
	}
}

func (t *TimeTrigger) Init() (<-chan bool, bool) {
	t.fire = make(chan bool)
	return t.fire, false
}

func (t *TimeTrigger) ProcessEvent(e Event) bool {
	if !t.At.IsZero() {
		var count uint64 = 0
		finalAt := t.At
		for {
			t.timer = time.NewTimer(finalAt.Sub(time.Now()))
			<-t.timer.C
			t.fire <- true
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
			t.fire <- true
			count++
			if !t.Forever && count >= t.Iterations {
				break
			}
		}
	}
	return false
}
