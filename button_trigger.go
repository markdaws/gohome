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

func (t *ButtonTrigger) GetName() string {
	return "Button Trigger"
}

func (t *ButtonTrigger) GetDescription() string {
	return "Triggers when a button is pressed one or more times"
}

func (t *ButtonTrigger) GetIngredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			Identifiable: Identifiable{
				ID:          "PressCount",
				Name:        "Press Count",
				Description: "The number of times the button should be pressed to activate the trigger",
			},
			Type: "number",
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "MaxDuration",
				Name:        "Max Duration (ms)",
				Description: "The maximum time in milliseconds from the first press that all presses must happen",
			},
			Type: "number",
		},
		//TODO: Button -> how to enumerate the right button?
		//Buttons live on devices
		//Devices live in systems
	}
}

//TODO: A recipe may or my not have the attributes set you can share those

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
