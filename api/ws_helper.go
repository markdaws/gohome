package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/log"
)

// Need check origin to allow cross-domain calls
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type WSHelper struct {
	monitor     *gohome.Monitor
	connections map[string]map[*connection]bool
	conn        *websocket.Conn
	mutex       sync.RWMutex
	updates     chan *gohome.ChangeBatch
}

type connection struct {
	monitorID string
	ws        *websocket.Conn
	writeChan chan bool
	readChan  chan bool
}

func NewWSHelper(monitor *gohome.Monitor) *WSHelper {
	h := WSHelper{
		monitor:     monitor,
		connections: make(map[string]map[*connection]bool),
		updates:     make(chan *gohome.ChangeBatch, 1000),
	}
	h.processUpdates()
	return &h
}

func (h *WSHelper) register(c *connection) {
	h.mutex.Lock()
	conns, ok := h.connections[c.monitorID]
	if !ok {
		conns = make(map[*connection]bool)
		h.connections[c.monitorID] = conns
	}
	conns[c] = true
	h.mutex.Unlock()
}

func (h *WSHelper) unregister(c *connection) {
	h.mutex.Lock()

	var conns map[*connection]bool
	conns, ok := h.connections[c.monitorID]
	if !ok {
		h.mutex.Unlock()
		return
	}

	if _, ok := conns[c]; !ok {
		h.mutex.Unlock()
		return
	}

	delete(conns, c)
	if len(conns) == 0 {
		delete(h.connections, c.monitorID)
	}
	h.mutex.Unlock()

	c.ws.Close()
	close(c.writeChan)
	close(c.readChan)
}

func (h *WSHelper) HTTPHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check the monitorID, use has to first subscribe and get an ID
		// before trying to stream the values
		monitorID := mux.Vars(r)["monitorID"]
		if _, ok := h.monitor.Group(monitorID); !ok {
			c.Close()
			return
		}

		conn := &connection{
			monitorID: monitorID,
			ws:        c,
			writeChan: make(chan bool),
			readChan:  make(chan bool),
		}
		h.register(conn)

		// When a connection registers, we need to ask the monitor to refresh all
		// values associated with it. Since we could have subscribed but not connected
		// yet and missed previous updates
		h.monitor.Refresh(monitorID, false)

		go conn.writeLoop(h)
		conn.readLoop(h)
	}
}

func (h *WSHelper) processUpdates() {
	go func() {
		for update := range h.updates {
			h.mutex.RLock()
			conns, ok := h.connections[update.MonitorID]
			if !ok || len(conns) == 0 {
				h.mutex.RUnlock()
				continue
			}

			connList := make([]*connection, 0, len(conns))
			for conn := range conns {
				connList = append(connList, conn)
			}
			h.mutex.RUnlock()

			evt := jsonMonitorGroupResponse{
				Sensors: make(map[string]jsonSensorAttr),
				Zones:   make(map[string]jsonZoneLevel),
			}
			for sensorID, attr := range update.Sensors {
				evt.Sensors[sensorID] = jsonSensorAttr{
					Name:     attr.Name,
					Value:    attr.Value,
					DataType: string(attr.DataType),
				}
			}
			for zoneID, level := range update.Zones {
				evt.Zones[zoneID] = jsonZoneLevel{
					Value: level.Value,
					R:     level.R,
					G:     level.G,
					B:     level.B,
				}
			}

			bytes, err := json.Marshal(evt)
			if err != nil {
				log.E("failed to marshal change batch to JSON for update: %s", err)
				continue
			}

			// Serial, if we ever get a lot of conncurrent users, would want to push
			// these in parallel
			for _, conn := range connList {
				conn.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err = conn.ws.WriteMessage(websocket.TextMessage, bytes)
				if err != nil {
					h.unregister(conn)
				}
			}
		}
	}()
}

// ========= gohome.MonitorDelegate interface ==============

// Update is the callback to the monitor service, it will get change notifications
// when zones and sensors update
func (h *WSHelper) Update(b *gohome.ChangeBatch) {
	h.updates <- b
}

func (h *WSHelper) Expired(monitorID string) {
	// The monitor ID has expired, close any connections associated
	// with this monitorID
	go func() {
		h.mutex.RLock()
		conns, ok := h.connections[monitorID]

		if !ok || len(conns) == 0 {
			h.mutex.RUnlock()
			return
		}

		connList := make([]*connection, 0, len(conns))
		for conn := range conns {
			connList = append(connList, conn)
		}
		h.mutex.RUnlock()

		for _, conn := range connList {
			h.unregister(conn)
		}
	}()
}

// =========================================================

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
			break
		}
	}
}
