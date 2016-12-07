package gohome_test

import (
	"testing"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome"
	"github.com/stretchr/testify/require"
)

type MockTime struct {
	now   time.Time
	after func(time.Duration) <-chan time.Time
}

func (mt MockTime) Now() time.Time {
	return mt.now
}

func (mt MockTime) After(d time.Duration) <-chan time.Time {
	return mt.after(d)
}

func TestSunrise(t *testing.T) {
	t.Parallel()

	mt := MockTime{
		now:   time.Date(2016, time.December, 5, 10, 0, 0, 0, time.UTC),
		after: func(d time.Duration) <-chan time.Time { return time.After(d) },
	}

	wasTriggered := false
	trigger := &gohome.TimeTrigger{
		Time:   mt,
		Mode:   gohome.TimeTriggerModeSunrise,
		Offset: 0,
		Days:   gohome.TimeTriggerDaysMon | gohome.TimeTriggerDaysFri,
		Triggered: func() {
			wasTriggered = true
		},
	}

	ch := make(chan evtbus.Event)
	trigger.StartConsuming(ch)

	// Send sunrise event, check offset
	ch <- &gohome.SunriseEvt{}

	// Send sunset - nothing should happen
	time.Sleep(time.Second * 1)
	require.True(t, wasTriggered)
}

func TestSunset(t *testing.T) {
	t.Parallel()

	mt := MockTime{
		now:   time.Date(2016, time.December, 5, 10, 0, 0, 0, time.UTC),
		after: func(d time.Duration) <-chan time.Time { return time.After(d) },
	}

	wasTriggered := false
	trigger := &gohome.TimeTrigger{
		Time:   mt,
		Mode:   gohome.TimeTriggerModeSunset,
		Offset: 0,
		Days:   gohome.TimeTriggerDaysMon | gohome.TimeTriggerDaysFri,
		Triggered: func() {
			wasTriggered = true
		},
	}

	ch := make(chan evtbus.Event)
	trigger.StartConsuming(ch)

	// Send sunrise event, check offset
	ch <- &gohome.SunsetEvt{}

	// Send sunset - nothing should happen
	time.Sleep(time.Second * 1)

	require.True(t, wasTriggered)
}

func TestExactWithDate(t *testing.T) {
	t.Parallel()

	mt := MockTime{
		now:   time.Date(2016, time.December, 5, 10, 0, 0, 0, time.UTC),
		after: func(d time.Duration) <-chan time.Time { return time.After(d) },
	}

	wasTriggered := false
	trigger := &gohome.TimeTrigger{
		Time: mt,
		Mode: gohome.TimeTriggerModeExact,
		At:   mt.Now().Add(time.Second * 1),
		Triggered: func() {
			wasTriggered = true
		},
	}

	ch := make(chan evtbus.Event)
	trigger.StartConsuming(ch)

	time.Sleep(time.Second * 2)

	require.True(t, wasTriggered)
}

func TestExactWithoutDate(t *testing.T) {
	t.Parallel()

	mt := MockTime{
		// This is a monday
		now: time.Date(2016, time.December, 5, 10, 10, 0, 0, time.UTC),
		after: func(d time.Duration) <-chan time.Time {
			// Return immediately
			c := make(chan time.Time, 1)
			c <- time.Now()
			return c
		},
	}

	// This is the same as the current mock time
	at := time.Date(
		0, 1, 1,
		mt.Now().Hour(), mt.Now().Minute(), mt.Now().Second(), 0, mt.Now().Location())

	// move in to the future
	at = at.Add(time.Second)

	evalCount := 0
	wasTriggered := false
	eval := func() {
		// This is called each time the trigger evaluates if it should run

		// If this is not the first time it is evaluating, then pretend we incremented
		// the current time by 23 hours to the next day (not 24 incase we jump past the
		// execution time), so that the test can run over multiple days. We only want the
		// test to execute on Monday and Friday
		if evalCount > 0 {
			mt.now = mt.now.Add(time.Hour * 23)
		}

		// Start on Monday, if eval is one, it means we ran on monday, so we should have got
		// an execute scene, then also on friday
		if evalCount == 1 || evalCount == 5 {
			require.True(t, wasTriggered)
			wasTriggered = false
		} else {
			require.False(t, wasTriggered)
		}
		evalCount++
	}

	trigger := &gohome.TimeTrigger{
		Time:       &mt,
		Mode:       gohome.TimeTriggerModeExact,
		Days:       gohome.TimeTriggerDaysMon | gohome.TimeTriggerDaysFri,
		At:         at,
		Evaluating: eval,
		Triggered: func() {
			wasTriggered = true
		},
	}

	ch := make(chan evtbus.Event)
	trigger.StartConsuming(ch)

	time.Sleep(time.Second * 2)
	require.Equal(t, 6, evalCount)
}
