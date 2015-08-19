package gohome

import "time"

type Trigger interface {
	Start() (<-chan bool, <-chan bool)
	Stop()
}

// At a certain time e.g. 8pm
// time no year, no month, no day, hour, minute, second
// After a certain delay every 5 minutes
// Iterations - certain number of times
// TODO: Be able to get sunrise/sunset time for a location: https://github.com/cpucycle/astrotime
type TimeTrigger struct {
	Iterations uint64
	Forever    bool
	At         time.Time
	Interval   time.Duration
	timer      *time.Timer
	ticker     *time.Ticker
	doneChan   chan bool
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
