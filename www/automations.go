package www

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

// RegisterAutomationHandlers registers all of the automation specific API REST routes
func RegisterAutomationHandlers(r *mux.Router, s *Server) {
	r.HandleFunc("/v1/automations", apiAutomationHandler(s.system)).Methods("GET")
	r.HandleFunc("/v1/automations/{ID}/test", apiAutomationTestHandler(s.system)).Methods("POST")
}

func apiAutomationHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")

		i := 0
		automations := system.Automations()
		items := make([]jsonAutomation, len(automations))
		for _, automation := range automations {
			items[i] = jsonAutomation{
				Name:   automation.Name,
				TempID: automation.TempID,
			}
			i++
		}

		if err := json.NewEncoder(w).Encode(items); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func apiAutomationTestHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		automationID := mux.Vars(r)["ID"]
		automation := system.AutomationByTempID(automationID)
		if automation == nil {
			respBadRequest(fmt.Sprintf("invalid automation ID: %s", automationID), w)
			return
		}

		automation.Trigger.Trigger()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}
