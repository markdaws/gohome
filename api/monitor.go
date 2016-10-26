package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

// RegisterMonitorHandlers registers all of the monitor specific REST API routes
func RegisterMonitorHandlers(r *mux.Router, s *apiServer) {
	r.HandleFunc("/api/v1/monitor/group", apiSubscribeHandler(s.system)).Methods("POST")
}

func apiSubscribeHandler(system *gohome.System) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: get the monitor group - call in to monitor
		w.WriteHeader(http.StatusOK)
	}
}
