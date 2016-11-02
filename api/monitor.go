package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/markdaws/gohome"
)

// RegisterMonitorHandlers registers all of the monitor specific REST API routes
func RegisterMonitorHandlers(r *mux.Router, s *apiServer) {
	wsHelper := NewWSHelper(s.system.Services.Monitor, s.system.Services.EvtBus)

	// Clients call to subscribe to items, api returns a monitorID that can then be used
	// to subscribe and unsubscribe to notifications
	r.HandleFunc("/api/v1/monitor/groups", apiSubscribeHandler(s.system, wsHelper)).Methods("POST")

	// extends the timeout period for a monitor groups
	r.HandleFunc("/api/v1/monitor/groups/{monitorID}", apiRefreshSubscribeHandler(s.system, wsHelper)).Methods("PUT")

	// deletes a monitor group
	r.HandleFunc("/api/v1/monitor/groups/{monitorID}", apiUnsubscribeHandler(s.system, wsHelper)).Methods("DELETE")

	// web socket for receiving new events
	r.HandleFunc("/api/v1/monitor/groups/{monitorID}", wsHelper.HTTPHandler())

}

func apiUnsubscribeHandler(system *gohome.System, wsHelper *WSHelper) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		monitorID := mux.Vars(r)["monitorID"]
		if _, ok := system.Services.Monitor.Group(monitorID); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		system.Services.Monitor.Unsubscribe(monitorID)
		w.WriteHeader(http.StatusOK)
	}
}

func apiRefreshSubscribeHandler(system *gohome.System, wsHelper *WSHelper) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		monitorID := mux.Vars(r)["monitorID"]
		if _, ok := system.Services.Monitor.Group(monitorID); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		system.Services.Monitor.SubscribeRenew(monitorID)
		w.WriteHeader(http.StatusOK)
	}
}

func apiSubscribeHandler(system *gohome.System, wsHelper *WSHelper) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var groupJSON jsonMonitorGroup
		if err = json.Unmarshal(body, &groupJSON); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		group := &gohome.MonitorGroup{
			Timeout: time.Duration(groupJSON.TimeoutInSeconds) * time.Second,
			Sensors: make(map[string]bool),
			Zones:   make(map[string]bool),
			Handler: wsHelper,
		}
		for _, sensorID := range groupJSON.SensorIDs {
			group.Sensors[sensorID] = true
		}
		for _, zoneID := range groupJSON.ZoneIDs {
			group.Zones[zoneID] = true
		}

		mID, err := system.Services.Monitor.Subscribe(group, false)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			MonitorID string `json:"monitorId"`
		}{MonitorID: mID})
	}
}
