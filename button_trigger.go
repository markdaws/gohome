package gohome

import (
	"fmt"
	"time"

	"github.com/nu7hatch/gouuid"
)

type ButtonTrigger struct {
	ButtonID   string
	PressCount int
	//TODO: ReleaseCount
	MaxDuration time.Duration
	//EventBroker or make this a consumer
	enabled    bool
	id         string
	startTime  time.Time
	pressCount int
	fireChan   chan bool
	doneChan   chan bool
}

func (t *ButtonTrigger) Type() string {
	return "gohome.ButtonTrigger"
}

func (t *ButtonTrigger) Name() string {
	return "Button Trigger"
}

func (t *ButtonTrigger) Description() string {
	return "Triggers when a button is pressed one or more times"
}

func (t *ButtonTrigger) Enabled() bool {
	return t.enabled
}

func (t *ButtonTrigger) SetEnabled(enabled bool) {
	t.enabled = enabled
}

func (t *ButtonTrigger) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			Identifiable: Identifiable{
				ID:          "ButtonID",
				Name:        "Button ID",
				Description: "The button ID associated with this trigger",
			},
			Type:     "string",
			Required: true,
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "PressCount",
				Name:        "Press Count",
				Description: "The number of times the button should be pressed to activate the trigger",
			},
			Type:     "integer",
			Required: true,
		},
		Ingredient{
			Identifiable: Identifiable{
				ID:          "MaxDuration",
				Name:        "Max Duration (ms)",
				Description: "The maximum time in milliseconds from the first press that all presses must happen",
			},
			Type:     "duration",
			Required: true,
		},
	}
}

func (t *ButtonTrigger) New() Trigger {
	return &ButtonTrigger{}
}

//TODO: Move the common code into a mixin
func (t *ButtonTrigger) Start() (<-chan bool, <-chan bool) {
	t.fireChan = make(chan bool)
	t.doneChan = make(chan bool)
	return t.fireChan, t.doneChan
}

func (t *ButtonTrigger) Stop() {
	//TODO: Should stop exit the function too?
	t.enabled = false
}

func (t *ButtonTrigger) EventConsumerID() string {
	if t.id == "" {
		id, err := uuid.NewV4()
		if err != nil {
			//TODO: error
		}
		t.id = id.String()
	}
	return t.id
}

func (t *ButtonTrigger) StartConsumingEvents() chan<- Event {
	c := make(chan Event)
	go func() {
		for e := range c {
			if !t.enabled {
				continue
			}

			if time.Now().After(t.startTime.Add(t.MaxDuration)) {
				t.startTime = time.Now()
				t.pressCount = 0
			}
			//fmt.Printf("ButtonTrigger: %s got event: %s\n", t.id, e.String())

			if e.ReplayCommand == nil ||
				e.ReplayCommand.GetType() != CTDevicePressButton {
				continue
			}

			// If matches the button ID associated with this trigger
			t.pressCount++
			//fmt.Printf("Current press count: %d\n", t.pressCount)
			if t.pressCount == t.PressCount {
				//fmt.Printf("Got %d presses\n", t.pressCount)
				//TODO: What if nobody listening
				t.fireChan <- true
				t.pressCount = 0
				t.startTime = time.Now()
			}
		}
		//TODO:
		fmt.Println("done")
		//t.doneChan <- true
	}()
	return c
}
