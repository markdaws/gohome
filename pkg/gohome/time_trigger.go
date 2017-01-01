package gohome

import (
	"fmt"
	"math"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/pkg/clock"
	"github.com/markdaws/gohome/pkg/log"
)

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
	Name      string
	Mode      string
	At        time.Time
	Days      uint32
	Time      clock.Time
	Triggered func()
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

	if hasDate {
		// Was this for before the current date, if so ignore
		delta := t.At.Sub(t.Time.Now())
		if delta < 0 {
			return
		}

		<-t.Time.After(delta)
		t.Triggered()
		return
	}

	for {
		now := t.Time.Now()
		absoluteAt := t.nextTriggerTime()

		log.V("TimeTrigger[%s] - next trigger time: %s", t.Name, absoluteAt)
		delta := absoluteAt.Sub(now)

		// Sleep until the correct time
		<-t.Time.After(delta)
		t.scheduleAction()

		// Small wait to make sure we don't re-run the automation on the same day
		time.Sleep(time.Second * 1)
	}
}

func (t *TimeTrigger) nextTriggerTime() time.Time {
	now := t.Time.Now()
	absoluteAt := time.Date(now.Year(), now.Month(), now.Day(), t.At.Hour(), t.At.Minute(), t.At.Second(), 0, now.Location())
	delta := absoluteAt.Sub(now)

	if delta < 0 {
		// Event was before this time on this day, go to the next day
		absoluteAt = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		absoluteAt = absoluteAt.Add(time.Hour * 24)
		absoluteAt = time.Date(absoluteAt.Year(), absoluteAt.Month(), absoluteAt.Day(),
			t.At.Hour(), t.At.Minute(), t.At.Second(), 0, now.Location())
	}
	return absoluteAt
}

func (t *TimeTrigger) scheduleAction() {
	now := t.Time.Now()

	// Make sure this is a day of the week we should execute this trigger
	dayOfWeekOrdinal := int(now.Weekday())

	// Convert time.Weekday to our representation of days of week
	daysValue := uint32(math.Pow(2, float64(dayOfWeekOrdinal)))

	if (t.Days & daysValue) != 0 {
		t.Triggered()
	}
}
