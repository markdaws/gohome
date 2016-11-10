package api

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

// RegisterButtonHandlers registers all of the button specific API REST routes
func RegisterButtonHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/buttons",
		apiButtonsHandler(s.system)).Methods("GET")
}

// ButtonsToJSON returns a jsonified slice of all the input buttons,
// this function also sorts the output
func ButtonsToJSON(btns map[string]*gohome.Button) []jsonButton {
	buttons := make(buttons, len(btns))
	var i int32
	for _, button := range btns {
		buttons[i] = jsonButton{
			ID:          button.ID,
			Name:        button.Name,
			Description: button.Description,
			Address:     button.Address,
			FullName:    button.Device.Name + " / " + button.Name,
		}
		i++
	}
	sort.Sort(buttons)
	return buttons
}

func apiButtonsHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(ButtonsToJSON(system.Buttons)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
