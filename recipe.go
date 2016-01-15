package gohome

import (
	"fmt"

	"github.com/markdaws/gohome/log"
	"github.com/nu7hatch/gouuid"
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

	system                 *System
	enabled                bool
	initedTrigger          bool
	triggerProcessesEvents bool
}

func NewRecipe(name, description string, enabled bool, t Trigger, a Action, s *System) (*Recipe, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &Recipe{
		ID:          id.String(),
		Name:        name,
		Description: description,
		Trigger:     t,
		Action:      a,
		Version:     "1",
		system:      s,
		enabled:     true,
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

func (r *Recipe) EventConsumerID() string {
	return r.Name + " - " + r.ID
}

func (r *Recipe) StartConsumingEvents() chan<- Event {
	log.V("%s started consuming events", r)

	c := make(chan Event)
	done := make(chan bool)

	var ignoreTrigger = false
	if !r.initedTrigger {
		fire, triggerProcessesEvents := r.Trigger.Init()
		r.initedTrigger = true
		r.triggerProcessesEvents = triggerProcessesEvents

		// Trigger could be something like a timer, can fire a signal
		// to indicate if has triggered, need to be able to handle it
		if fire != nil {
			go func() {
				for {
					select {
					case <-done:
						//TODO: Test
						break
					case f := <-fire:
						if r.enabled && !ignoreTrigger {
							if f {
								executeAction(r)
							}
						}
					}
				}
			}()
		}
	}

	if r.triggerProcessesEvents {
		go func() {
			for e := range c {
				if !r.enabled {
					continue
				}

				if r.Trigger.ProcessEvent(e) {
					executeAction(r)
				}
			}
			close(done)
			ignoreTrigger = true
			log.V("%s stopped consuming events", r)
		}()
	}
	return c
}

func executeAction(r *Recipe) {
	err := r.Action.Execute(r.system)
	if err != nil {
		log.E("%s action failed: %s", r, err)
	}
}
