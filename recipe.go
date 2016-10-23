package gohome

import (
	"fmt"

	"github.com/go-home-iot/event-bus"
	"github.com/markdaws/gohome/log"
)

//TODO: Rules for trigger/action writers, don't have pointers to objects, have ids, use the
//system object to get the items ou want access to, otherwise won't work on save/reload

type Recipe struct {
	ID          string
	Name        string
	Description string
	Trigger     Trigger
	Action      Action
	Version     string

	system  *System
	enabled bool
}

func NewRecipe(name, description string, enabled bool, t Trigger, a Action, s *System) (*Recipe, error) {
	return &Recipe{
		ID:          s.NextGlobalID(),
		Name:        name,
		Description: description,
		Trigger:     t,
		Action:      a,
		Version:     "1",
		system:      s,
		enabled:     enabled,
	}, nil
}

func (r *Recipe) String() string {
	return fmt.Sprintf("Recipe[%s]", r.Name)
}

func (r *Recipe) Enabled() bool {
	return r.enabled
}

func (r *Recipe) SetEnabled(enabled bool) {
	r.enabled = enabled
}

func (r *Recipe) ConsumerName() string {
	return r.Name + " - " + r.ID
}

func (r *Recipe) StartConsuming(ch chan evtbus.Event) {
	/*
		log.V("%s started consuming events", r)

		triggerDone := make(chan bool)
		fire, triggerProcessesEvents := r.Trigger.Init(triggerDone)

		done := make(chan bool)

		// Trigger could be something like a timer, can fire a signal
		// to indicate if has triggered, need to be able to handle it
		if fire != nil {
			go func() {
				for {
					select {
					case _, ok := <-done:
						if !ok {
							done = nil
						}
					case f := <-fire:
						if r.enabled {
							if f {
								log.V("%s trigger fired", r)
								executeAction(r)
							}
						}
					}
					if done == nil {
						break
					}
				}
			}()
		}

		c := make(chan event.Event)
		go func() {
			for e := range c {
				if !r.enabled || !triggerProcessesEvents {
					continue
				}

				if r.Trigger.ProcessEvent(e) {
					log.V("%s trigger fired", r)
					executeAction(r)
				}
			}

			// No longer a consumer, signal to stop trigger and listening for any more events
			close(done)
			close(triggerDone)

			log.V("%s stopped consuming events", r)
		}()
		return c
	*/
}

func (r *Recipe) StopConsuming() {
	//TODO:
}

func executeAction(r *Recipe) {
	err := r.Action.Execute(r.system)
	if err != nil {
		log.E("%s action failed: %s", r, err)
	}
}
