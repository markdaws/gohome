package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/markdaws/gohome"
)

// RegisterMonitorHandlers registers all of the monitor specific REST API routes
func RegisterMonitorHandlers(r *mux.Router, s *apiServer) {
	// Web Socket helper
	wsHelper := NewWSHelper(s.system.Services.Monitor)

	r.HandleFunc("/api/v1/monitor/groups/{monitorID}", wsHelper.HTTPHandler())
	r.HandleFunc("/api/v1/monitor/groups", apiSubscribeHandler(s.system, wsHelper)).Methods("POST")

	//TODO:
	//r.HandleFunc("/api/v1/monitor/groups/{monitorID}", apiUnsubscribeHandler(s.system)).Methods("DELETE")
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

		if len(groupJSON.SensorIDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		group := &gohome.MonitorGroup{
			TimeoutInSeconds: groupJSON.TimeoutInSeconds,
			Sensors:          make(map[string]bool),
			Zones:            make(map[string]bool),
			Handler:          wsHelper,
		}
		for _, sensorID := range groupJSON.SensorIDs {
			group.Sensors[sensorID] = true
		}
		for _, zoneID := range groupJSON.ZoneIDs {
			group.Zones[zoneID] = true
		}

		mID := system.Services.Monitor.Subscribe(group, true)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			MonitorID string `json:"monitorId"`
		}{MonitorID: mID})
	}
}

//TODO: Move
// Need check origin to allow cross-domain calls
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type WSHelper struct {
	//index connetions by monitor id
	monitor     *gohome.Monitor
	connections map[*connection]bool
	conn        *websocket.Conn
}

func NewWSHelper(monitor *gohome.Monitor) *WSHelper {
	h := WSHelper{
		monitor:     monitor,
		connections: make(map[*connection]bool),
	}
	return &h
}

func (l *WSHelper) register(c *connection) {
	l.connections[c] = true
}
func (l *WSHelper) unregister(c *connection) {
	//TODO: goroutine safe
	if _, ok := l.connections[c]; ok {
		delete(l.connections, c)
		c.ws.Close()
		close(c.writeChan)
		close(c.readChan)
	}
}

func (h *WSHelper) Update(monitorID string, b *gohome.ChangeBatch) {
	//TODO: Index connections by monitor id -> could be multiple connections per monitor id
	for conn := range h.connections {
		if conn.monitorID == monitorID {
			evt := jsonMonitorGroupResponse{
				Sensors: make(map[string]jsonSensorAttr),
				Zones:   make(map[string]jsonZoneLevel),
			}
			for sensorID, attr := range b.Sensors {
				evt.Sensors[sensorID] = jsonSensorAttr{
					Name:     attr.Name,
					Value:    attr.Value,
					DataType: string(attr.DataType),
				}
			}
			fmt.Println("ws here")
			for zoneID, level := range b.Zones {
				evt.Zones[zoneID] = jsonZoneLevel{
					Value: level.Value,
					R:     level.R,
					G:     level.G,
					B:     level.B,
				}
			}

			b, err := json.Marshal(evt)
			if err != nil {
				//TODO: Log error
				return
			}

			conn.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err = conn.ws.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				h.unregister(conn)
			}
			return
		}
	}
}

func (h *WSHelper) HTTPHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		monitorID := mux.Vars(r)["monitorID"]
		conn := &connection{
			monitorID: monitorID,
			ws:        c,
			writeChan: make(chan bool),
			readChan:  make(chan bool),
		}
		h.register(conn)
		go conn.writeLoop(h)
		conn.readLoop(h)
	}
}

type connection struct {
	monitorID string
	ws        *websocket.Conn
	writeChan chan bool
	readChan  chan bool
}

func (c *connection) writeLoop(l *WSHelper) {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	var exit = false
	for {
		select {
		case _, ok := <-c.writeChan:
			if !ok {
				exit = true
			}
		case <-ticker.C:
			// Making sure the client is still alive
			c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				l.unregister(c)
				exit = true
			}
		}

		if exit {
			break
		}
	}
}

func (c *connection) readLoop(l *WSHelper) {
	// have to have a read loop otherwise ping/pong don't work
	defer func() {
		l.unregister(c)
	}()
	c.ws.SetReadLimit(1024)

	maxWait := 60 * time.Second
	c.ws.SetReadDeadline(time.Now().Add(maxWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(maxWait))
		return nil
	})
	for {
		// If the client closes we get a 1001 error here
		if _, _, err := c.ws.ReadMessage(); err != nil {
			//TODO: remove
			fmt.Println(err)
			break
		}
	}
}
