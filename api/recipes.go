package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/store"
)

func RegisterRecipeHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/recipes",
		apiRecipesHandlerPost(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/recipes/{id}",
		apiRecipeHandler(s.system, s.recipeManager)).Methods("POST")
	r.HandleFunc("/api/v1/recipes/{id}",
		apiRecipeHandlerDelete(s.system, s.recipeManager)).Methods("DELETE")
	r.HandleFunc("/api/v1/recipes",
		apiRecipesHandlerGet(s.system, s.recipeManager)).Methods("GET")
}

func apiRecipesHandlerPost(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var data map[string]interface{}
		if err = json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipe, err := recipeManager.UnmarshalNewRecipe(data)
		if err != nil {
			errBad := err.(*gohome.ErrUnmarshalRecipe)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(struct {
				ParamID     string `json:"paramId"`
				ErrorType   string `json:"errorType"`
				Description string `json:"description"`
			}{
				ParamID:     errBad.ParamID,
				ErrorType:   errBad.ErrorType,
				Description: errBad.Description,
			})
			return
		}

		recipeManager.RegisterAndStart(recipe)
		err = store.SaveSystem(system, recipeManager)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct {
			ID string `json:"id"`
		}{ID: recipe.ID})
	}
}

func apiRecipeHandler(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var data struct {
			Enabled bool `json:"enabled"`
		}
		if err = json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipeID := mux.Vars(r)["id"]
		recipe := recipeManager.RecipeByID(recipeID)
		if recipe == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = recipeManager.EnableRecipe(recipe, data.Enabled)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		store.SaveSystem(system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiRecipeHandlerDelete(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		recipeID := mux.Vars(r)["id"]
		recipe := recipeManager.RecipeByID(recipeID)
		if recipe == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := recipeManager.DeleteRecipe(recipe)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		store.SaveSystem(system, recipeManager)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(struct{}{})
	}
}

func apiRecipesHandlerGet(system *gohome.System, recipeManager *gohome.RecipeManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		recipes := system.Recipes
		jsonRecipes := make(jsonRecipes, len(recipes))

		i := 0
		for _, recipe := range recipes {
			jsonRecipes[i] = jsonRecipe{
				ID:          recipe.ID,
				Name:        recipe.Name,
				Description: recipe.Description,
				Enabled:     recipe.Enabled(),
			}
			i++
		}
		sort.Sort(jsonRecipes)
		if err := json.NewEncoder(w).Encode(jsonRecipes); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
