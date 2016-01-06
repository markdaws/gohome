package gohome

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
)

//TODO: rename
type EventLogger struct {
	id       string
	upgrader websocket.Upgrader
	conn     *websocket.Conn
}

func (l *EventLogger) EventConsumerID() string {
	if l.id == "" {
		id, err := uuid.NewV4()
		if err != nil {
			//TODO: error
		}
		l.id = id.String()
	}
	return l.id
}

func (l *EventLogger) StartConsumingEvents() chan<- Event {
	c := make(chan Event)
	go func() {
		for e := range c {
			if l.conn == nil {
				continue
			}

			l.conn.WriteMessage(websocket.TextMessage, []byte(e.String()))
		}
	}()
	return c
}

func (l *EventLogger) HTTPHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := l.upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		l.conn = c
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				break
			}
			fmt.Printf("recv: %s", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		}
	}
}
