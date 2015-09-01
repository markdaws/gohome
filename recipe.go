package gohome

type Recipe struct {
	ID          string
	Name        string
	Description string
	Trigger     Trigger
	Action      Action
}

// Loads all recipes
func LoadRecipesFromPath(path string) []*Recipe {
	//TODO:
	return nil
}

// Saves/updates a recipe
func SaveRecipe(e *Recipe, path string) error {
	//TODO:
	return nil
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
