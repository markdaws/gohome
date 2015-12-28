package gohome

import "github.com/nu7hatch/gouuid"

//TODO: Rules for trigger/action writers, don't have pointers to objects, have ids, use the
//system object to get the items ou want access to, otherwise won't work on save/reload

type Recipe struct {
	Identifiable Identifiable
	Trigger      Trigger
	Action       Action
	Version      string
	system       *System
}

func NewRecipe(name, description string, enabled bool, t Trigger, a Action, s *System) (*Recipe, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	t.SetEnabled(enabled)
	return &Recipe{
		Identifiable: Identifiable{
			ID:          id.String(),
			Name:        name,
			Description: description,
		},
		Trigger: t,
		Action:  a,
		Version: "1",
		system:  s,
	}, nil
}

func (r *Recipe) Start() <-chan bool {
	fireChan, doneChan := r.Trigger.Start()
	go func() {
		for {
			select {
			case <-fireChan:
				r.Action.Execute(r.system)

			case <-doneChan:
				doneChan = nil
			}

			if doneChan == nil {
				break
			}
		}
	}()
	return doneChan
}

func (r *Recipe) Stop() {
	r.Trigger.Stop()
}
