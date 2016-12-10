package gohome

import (
	"time"

	"github.com/go-home-iot/event-bus"
)

// FeatureTrigger is a trigger that can be used to fire based on a features attributes changing
type FeatureTrigger struct {
	Condition *condition
	Triggered func()

	// Count the number of times the trigger must evaluate to true before triggering
	Count int

	// Duration the maximum amount of time allowed for the trigger to reach the Count amount. If the
	// count amount is notreached in this time the time is reset e.g. you may say get 3 events within
	// 3000 milliseconds
	Duration time.Duration

	trueCount int
	startTime time.Time
}

func (e *FeatureTrigger) ConsumerName() string {
	return "EventTrigger"
}

func (e *FeatureTrigger) StartConsuming(ch chan evtbus.Event) {
	go func() {
		for evt := range ch {
			attrEvt, ok := evt.(*FeatureAttrsChangedEvt)
			if !ok {
				continue
			}

			isTrue := e.Condition.Evaluate(attrEvt)
			if isTrue {
				if time.Now().After(e.startTime.Add(e.Duration)) {
					e.trueCount = 1
					e.startTime = time.Now()
				} else {
					e.trueCount++
				}

				if e.Count == 0 {
					// User has not set a count, just trigger
					e.Triggered()
				} else if e.trueCount == e.Count {
					// Reached the trigger amount
					e.Triggered()
				}
			}
		}
	}()
}

func (e *FeatureTrigger) StopConsuming() {
	//TODO:
}

func (e *FeatureTrigger) Trigger() {
	e.Triggered()
}
