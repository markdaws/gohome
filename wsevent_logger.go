package gohome

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/markdaws/gohome/event"
	"github.com/nu7hatch/gouuid"
)

type WSEventLogger interface {
	HTTPHandler() func(http.ResponseWriter, *http.Request)
	event.Consumer
}

var upgrader websocket.Upgrader

type wsEventLogger struct {
	id          string
	connections map[*connection]bool
	conn        *websocket.Conn
	devicer     Devicer
}

func NewWSEventLogger(devicer Devicer) WSEventLogger {
	c := wsEventLogger{
		devicer:     devicer,
		connections: make(map[*connection]bool),
	}
	return &c
}

func (l *wsEventLogger) register(c *connection) {
	l.connections[c] = true
}
func (l *wsEventLogger) unregister(c *connection) {
	if _, ok := l.connections[c]; ok {
		delete(l.connections, c)
		c.ws.Close()
		close(c.writeChan)
		close(c.readChan)
	}
}

func (l *wsEventLogger) EventConsumerID() string {
	if l.id == "" {
		id, err := uuid.NewV4()
		if err != nil {
			//TODO: error
		}
		l.id = id.String()
	}
	return l.id
}

type jsonEvent struct {
	ID              string    `json:"id"`
	Time            time.Time `json:"datetime"`
	RawMessage      string    `json:"rawMessage"`
	FriendlyMessage string    `json:"friendlyMessage"`
	DeviceName      string    `json:"deviceName"`
}

func (l *wsEventLogger) StartConsumingEvents() chan<- event.Event {
	c := make(chan event.Event)

	go func() {
		for {
			select {
			case e := <-c:
				// Don't block event broker
				go func() {
					//TODO: parellelize?
					for conn, _ := range l.connections {
						dev := l.devicer.FromID(e.DeviceID)
						var devName string = "unknown"
						if dev != nil {
							devName = dev.Name()
						}
						evt := jsonEvent{
							ID:              strconv.Itoa(e.ID),
							Time:            e.Time,
							RawMessage:      e.OriginalString,
							DeviceName:      devName,
							FriendlyMessage: e.String(),
						}
						b, err := json.Marshal(evt)
						if err != nil {
							//TODO: Log error
							continue
						}

						conn.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
						err = conn.ws.WriteMessage(websocket.TextMessage, b)
						if err != nil {
							l.unregister(conn)
						}
					}
				}()
			}
		}
	}()
	return c
}

func (l *wsEventLogger) HTTPHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conn := &connection{
			ws:        c,
			writeChan: make(chan bool),
			readChan:  make(chan bool),
		}
		l.register(conn)
		go conn.writeLoop(l)
		conn.readLoop(l)
	}
}

type connection struct {
	ws        *websocket.Conn
	writeChan chan bool
	readChan  chan bool
}

func (c *connection) writeLoop(l *wsEventLogger) {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	var exit bool = false
	for {
		select {
		case _, ok := <-c.writeChan:
			if !ok {
				exit = true
			}
		case <-ticker.C:
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

func (c *connection) readLoop(l *wsEventLogger) {
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
			fmt.Println(err)
			break
		}
	}
}
