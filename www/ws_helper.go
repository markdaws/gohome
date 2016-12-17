package www

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-home-iot/event-bus"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/attr"
	"github.com/markdaws/gohome/log"
)

// Need check origin to allow cross-domain calls
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type WSHelper struct {
	monitor     *gohome.Monitor
	evtBus      *evtbus.Bus
	nextID      int64
	connections map[string]map[*connection]bool
	conn        *websocket.Conn
	mutex       sync.RWMutex
	updates     chan *gohome.ChangeBatch
}

type connection struct {
	monitorID    string
	connectionID string
	ws           *websocket.Conn
	writeChan    chan bool
	readChan     chan bool
}

func NewWSHelper(monitor *gohome.Monitor, evtBus *evtbus.Bus) *WSHelper {
	h := WSHelper{
		monitor:     monitor,
		evtBus:      evtBus,
		nextID:      time.Now().UnixNano(),
		connections: make(map[string]map[*connection]bool),
		updates:     make(chan *gohome.ChangeBatch, 1000),
	}
	h.processUpdates()
	return &h
}

func (h *WSHelper) register(c *connection) {
	log.V("WSHelper - registering connection, monitorID: " + c.monitorID)

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

	log.V("WSHelper - unregister connection, monitorID: " + c.monitorID)
	c.ws.Close()
	close(c.writeChan)
	close(c.readChan)

	h.evtBus.Enqueue(&gohome.ClientDisconnectedEvt{ConnectionID: c.connectionID})
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
			connectionID: strconv.FormatInt(h.nextID, 10),
			monitorID:    monitorID,
			ws:           c,
			writeChan:    make(chan bool),
			readChan:     make(chan bool),
		}
		h.nextID++

		h.register(conn)
		go conn.writeLoop(h)

		origin := ""
		if orig, ok := r.Header["Origin"]; ok && len(orig) == 1 {
			origin = orig[0]
		}

		// Let the system know a new client has connected
		h.evtBus.Enqueue(&gohome.ClientConnectedEvt{
			MonitorID:    monitorID,
			Origin:       origin,
			ConnectionID: conn.connectionID,
		})

		// When a connection registers, we need to ask the monitor to refresh all
		// values associated with it. Since we could have subscribed but not connected
		// yet and missed previous updates
		h.monitor.Refresh(monitorID, false)

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
				Features: make(map[string]map[string]*attr.Attribute),
			}
			for featureID, attrs := range update.Features {
				evt.Features[featureID] = attrs
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
		log.V("WSHelper - expired connection, monitorID: " + monitorID)

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
