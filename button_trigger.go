package gohome

import (
	"fmt"
	"time"
)

type ButtonTrigger struct {
	//TODO: Need to expose buttons on devices
	Button     Button
	PressCount int
	//ReleaseCount
	MaxDuration time.Duration
	//EventBroker or make this a consumer
	Enabled bool

	startTime  time.Time
	pressCount int
	fireChan   chan bool
	doneChan   chan bool
}

func (t *ButtonTrigger) Start() (<-chan bool, <-chan bool) {
	t.Enabled = true
	t.fireChan = make(chan bool)
	t.doneChan = make(chan bool)
	return t.fireChan, t.doneChan
}

func (t *ButtonTrigger) Stop() {
	t.Enabled = false
}

func (t *ButtonTrigger) StartConsumingEvents() chan<- Event {
	c := make(chan Event)
	go func() {
		for e := range c {
			if !t.Enabled {
				continue
			}

			if time.Now().After(t.startTime.Add(t.MaxDuration)) {
				t.startTime = time.Now()
				t.pressCount = 0
			}
			fmt.Println("ButtonTrigger: got event: ", e.String())

			if e.ReplayCommand == nil ||
				e.ReplayCommand.GetType() != CTDevicePressButton {
				continue
			}

			// If matches
			t.pressCount++
			if t.pressCount == t.PressCount {
				fmt.Printf("Got %d presses\n", t.pressCount)
				//TODO: What if nobody listening
				t.fireChan <- true
			}
		}
		//TODO:
		fmt.Println("done")
		//t.doneChan <- true
	}()
	return c
}
