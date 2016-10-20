package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

// RegisterCookBookHandlers registers all of the CookBook specific API REST routes
func RegisterCookBookHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/cookbooks",
		apiCookBooksHandler(s.recipeManager.CookBooks)).Methods("GET")
	r.HandleFunc("/api/v1/cookbooks/{id}",
		apiCookBookHandler(s.recipeManager.CookBooks)).Methods("GET")
}

func apiCookBooksHandler(cookBooks []*gohome.CookBook) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")

		type jsonCookBook struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			LogoURL     string `json:"logoUrl"`
		}

		//TODO: Return in a consistent order
		jsonCookBooks := make([]jsonCookBook, len(cookBooks))
		for i, cookBook := range cookBooks {
			jsonCookBooks[i] = jsonCookBook{
				ID:          cookBook.ID,
				Name:        cookBook.Name,
				Description: cookBook.Description,
				LogoURL:     cookBook.LogoURL,
			}
		}
		if err := json.NewEncoder(w).Encode(jsonCookBooks); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiCookBookHandler(cookBooks []*gohome.CookBook) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")

		type jsonIngredient struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Type        string `json:"type"`
			Reference   string `json:"reference"`
		}
		type jsonTrigger struct {
			ID          string           `json:"id"`
			Name        string           `json:"name"`
			Description string           `json:"description"`
			Ingredients []jsonIngredient `json:"ingredients"`
		}
		type jsonAction struct {
			ID          string           `json:"id"`
			Name        string           `json:"name"`
			Description string           `json:"description"`
			Ingredients []jsonIngredient `json:"ingredients"`
		}
		type jsonCookBook struct {
			ID          string        `json:"id"`
			Name        string        `json:"name"`
			Description string        `json:"description"`
			LogoURL     string        `json:"logoUrl"`
			Triggers    []jsonTrigger `json:"triggers"`
			Actions     []jsonAction  `json:"actions"`
		}

		vars := mux.Vars(r)
		cbID := vars["id"]
		var found = false
		for _, c := range cookBooks {
			if c.ID != cbID {
				continue
			}

			jsonTriggers := make([]jsonTrigger, len(c.Triggers))
			for i, t := range c.Triggers {
				jsonTriggers[i] = jsonTrigger{
					ID:          t.Type(),
					Name:        t.Name(),
					Description: t.Description(),
					Ingredients: make([]jsonIngredient, len(t.Ingredients())),
				}

				for j, ing := range t.Ingredients() {
					jsonTriggers[i].Ingredients[j] = jsonIngredient{
						ID:          ing.ID,
						Name:        ing.Name,
						Description: ing.Description,
						Type:        ing.Type,
					}
				}
			}

			// for each trigger need to json all ingredients
			jsonActions := make([]jsonAction, len(c.Actions))
			for i, a := range c.Actions {
				jsonActions[i] = jsonAction{
					ID:          a.Type(),
					Name:        a.Name(),
					Description: a.Description(),
					Ingredients: make([]jsonIngredient, len(a.Ingredients())),
				}

				for j, ing := range a.Ingredients() {
					jsonActions[i].Ingredients[j] = jsonIngredient{
						ID:          ing.ID,
						Name:        ing.Name,
						Description: ing.Description,
						Type:        ing.Type,
					}
				}
			}

			// for each action need to json all ingredients
			jsonCookBook := jsonCookBook{
				ID:          c.ID,
				Name:        c.Name,
				Description: c.Description,
				LogoURL:     c.LogoURL,
				Triggers:    jsonTriggers,
				Actions:     jsonActions,
			}
			if err := json.NewEncoder(w).Encode(jsonCookBook); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			found = true
			break
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
