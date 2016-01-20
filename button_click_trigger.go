package gohome

import (
	"time"

	"github.com/markdaws/gohome/cmd"
)

type ButtonClickTrigger struct {
	ButtonID    string
	ClickCount  int
	MaxDuration time.Duration

	startTime  time.Time
	clickCount int
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

func (t *ButtonClickTrigger) Ingredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			ID:          "ButtonID",
			Name:        "Button ID",
			Description: "The button ID associated with this trigger",
			Type:        "string",
			Required:    true,
			Reference:   "button",
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

func (t *ButtonClickTrigger) Init(done <-chan bool) (<-chan bool, bool) {
	return nil, true
}

func (t *ButtonClickTrigger) ProcessEvent(e Event) bool {
	if time.Now().After(t.startTime.Add(t.MaxDuration)) {
		t.startTime = time.Now()
		t.clickCount = 0
	}

	cmd, ok := e.ReplayCommand.(*cmd.ButtonRelease)
	if !ok {
		return false
	}

	if cmd.ButtonGlobalID != t.ButtonID {
		return false
	}

	t.clickCount++
	if t.clickCount == t.ClickCount {
		t.clickCount = 0
		t.startTime = time.Now()
		return true
	}
	return false
}
