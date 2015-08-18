package gohome

type Recipe struct {
	Id          string
	Name        string
	Description string
	Trigger     Trigger
	Action      Action
}

func (r *Recipe) Start() <-chan bool {
	fireChan, doneChan := r.Trigger.Start()
	go func() {
		for {
			select {
			case <-fireChan:
				r.Action.Execute()

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
