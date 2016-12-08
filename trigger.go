package gohome

import (
	"fmt"
	"math"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/clock"
)

type Trigger interface {
	evtbus.Consumer
	Trigger()
}

const (
	// TimeTriggerModeSunrise - the time trigger is relative to sunrise
	TimeTriggerModeSunrise string = "sunrise"

	// TimeTriggerModeSunset - the time trigger is relative to sunset
	TimeTriggerModeSunset string = "sunset"

	// TimeTriggerModeExact - the time trigger is an exact time
	TimeTriggerModeExact string = "exact"
)

const (
	TimeTriggerDaysSun   uint32 = 1
	TimeTriggerDaysMon   uint32 = 2
	TimeTriggerDaysTues  uint32 = 4
	TimeTriggerDaysWed   uint32 = 8
	TimeTriggerDaysThurs uint32 = 16
	TimeTriggerDaysFri   uint32 = 32
	TimeTriggerDaysSat   uint32 = 64
)

// TimeTrigger is a trigger that can be used to execute actions either at sunrise/sunset or at an exact
// time.  You can also specify for the trigger to only fire on certain days of the week
type TimeTrigger struct {
	Mode       string
	Offset     time.Duration
	At         time.Time
	Days       uint32
	Time       clock.Time
	Evaluating func()
	Triggered  func()
	//TODO:
	//Start       time.Time
	//TODO:
	//End         time.Time
}

func (t *TimeTrigger) Trigger() {
	t.Triggered()
}

func (t *TimeTrigger) ConsumerName() string {
	return fmt.Sprintf("timetrigger")
}

func (t *TimeTrigger) StartConsuming(ch chan evtbus.Event) {
	go func() {
		switch t.Mode {
		case TimeTriggerModeSunrise:
			t.scheduleSunrise(ch)
		case TimeTriggerModeSunset:
			t.scheduleSunset(ch)
		case TimeTriggerModeExact:
			t.scheduleExact()
		}
	}()
}

func (t *TimeTrigger) StopConsuming() {
	//TODO:
}

func (t *TimeTrigger) scheduleSunrise(ch chan evtbus.Event) {
	for e := range ch {
		if _, ok := e.(*SunriseEvt); !ok {
			continue
		}
		t.scheduleAction()
	}
}

func (t *TimeTrigger) scheduleSunset(ch chan evtbus.Event) {
	for e := range ch {
		if _, ok := e.(*SunsetEvt); !ok {
			continue
		}
		t.scheduleAction()
	}
}

func (t *TimeTrigger) scheduleExact() {
	// If the time does not have a date it will be 0000 as the year (the null time)
	// if we have a date then this fires only once, otherwise if it doesn't have a
	// date it is just a time so we look at the days of the week to see if it should
	// execute
	hasDate := t.At.Year() != 0

	count := 0
	for {
		// Hook to let callers know the trigger is evaluating
		if t.Evaluating != nil {
			t.Evaluating()
		}

		now := t.Time.Now()
		count++
		if count > 5 {
			return
		}

		if hasDate {
			// Was this for before the current date, if so ignore
			delta := t.At.Sub(now)
			if delta < 0 {
				return
			}

			// If the At value has a date then this trigger is for an explicit date/time
			// so it's only executed once
			<-t.Time.After(delta)
			t.Triggered()

			// don't run again, was an explicit date/time trigger
			return
		} else {
			// Make sure this is a day of the week we should execute this trigger
			dayOfWeekOrdinal := int(now.Weekday())

			// Convert time.Weekday to our representation of days of week
			daysValue := uint32(math.Pow(2, float64(dayOfWeekOrdinal)))

			// This is a day of the week when we should execute this actions
			if (t.Days & daysValue) != 0 {
				// Figure out when the next time will be
				absoluteAt := time.Date(now.Year(), now.Month(), now.Day(), t.At.Hour(), t.At.Minute(), t.At.Second(), 0, now.Location())
				delta := absoluteAt.Sub(now)

				// Make sure we haven't gone past the time
				if delta >= 0 {
					<-t.Time.After(delta)
					t.Triggered()
				}
			}

			// Wait until the following day before we check again
			now = t.Time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), t.At.Hour(), t.At.Minute(), t.At.Second(), 0, now.Location())
			next = next.Add(time.Hour * 24)

			<-t.Time.After(next.Sub(now))
		}
	}
}

func (t *TimeTrigger) scheduleAction() {
	// Got a sunrise event - see if we should fire the trigger now or if we
	// need to add an offset and if it should fire on this day of the week
	now := t.Time.Now()

	// Make sure this is a day of the week we should execute this trigger
	dayOfWeekOrdinal := int(now.Weekday())

	// Convert time.Weekday to our representation of days of week
	daysValue := uint32(math.Pow(2, float64(dayOfWeekOrdinal)))

	if (t.Days & daysValue) != 0 {
		go func() {
			<-t.Time.After(t.Offset)
			t.Triggered()
		}()
	}
}
