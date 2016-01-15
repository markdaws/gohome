package gohome

import (
	"time"

	"github.com/nu7hatch/gouuid"
)

type ButtonClickTrigger struct {
	ButtonID    string
	ClickCount  int
	MaxDuration time.Duration

	enabled    bool
	id         string
	startTime  time.Time
	clickCount int
	fireChan   chan bool
	doneChan   chan bool
}

func (t *ButtonClickTrigger) Type() string {
	return "gohome.ButtonClickTrigger"
}

func (t *ButtonClickTrigger) Name() string {
	return "Button Trigger"
}

func (t *ButtonClickTrigger) Description() string {
	return "Triggers when a button is pressed one or more times"
}

func (t *ButtonClickTrigger) Enabled() bool {
	return t.enabled
}

func (t *ButtonClickTrigger) SetEnabled(enabled bool) {
	t.enabled = enabled
}

func (t *ButtonClickTrigger) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			ID:          "ButtonID",
			Name:        "Button ID",
			Description: "The button ID associated with this trigger",
			Type:        "string",
			Required:    true,
		},
		Ingredient{
			ID:          "ClickCount",
			Name:        "Click Count",
			Description: "The number of times the button should be clicked to activate the trigger",
			Type:        "integer",
			Required:    true,
		},
		Ingredient{
			ID:          "MaxDuration",
			Name:        "Max Duration (ms)",
			Description: "The maximum time in milliseconds from the first click that all clicks must happen",
			Type:        "duration",
			Required:    true,
		},
	}
}

func (t *ButtonClickTrigger) New() Trigger {
	return &ButtonClickTrigger{}
}

func (t *ButtonClickTrigger) Start() (<-chan bool, <-chan bool) {
	t.fireChan = make(chan bool)
	t.doneChan = make(chan bool)
	return t.fireChan, t.doneChan
}

func (t *ButtonClickTrigger) Stop() {
	//TODO: Should stop exit the function too?
	t.enabled = false
}

func (t *ButtonClickTrigger) EventConsumerID() string {
	if t.id == "" {
		id, err := uuid.NewV4()
		if err != nil {
			//TODO: error
		}
		t.id = id.String()
	}
	return t.id
}

func (t *ButtonClickTrigger) StartConsumingEvents() chan<- Event {
	c := make(chan Event)
	go func() {
		for e := range c {
			if !t.enabled {
				continue
			}

			if time.Now().After(t.startTime.Add(t.MaxDuration)) {
				t.startTime = time.Now()
				t.clickCount = 0
			}

			cmd := e.ReplayCommand
			if cmd == nil ||
				cmd.CMDType() != CTDeviceReleaseButton {
				continue
			}

			btn, ok := e.Source.(*Button)
			if !ok {
				continue
			}

			if btn.GlobalID != t.ButtonID {
				continue
			}

			t.clickCount++
			if t.clickCount == t.ClickCount {
				t.fireChan <- true
				t.clickCount = 0
				t.startTime = time.Now()
			}
		}

		//TODO:
		//t.doneChan <- true
	}()
	return c
}
