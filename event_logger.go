package gohome

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

type jsonEvent struct {
	ID              string    `json:"id"`
	Time            time.Time `json:"datetime"`
	RawMessage      string    `json:"rawMessage"`
	FriendlyMessage string    `json:"friendlyMessage"`
}

func (l *EventLogger) StartConsumingEvents() chan<- Event {
	c := make(chan Event)
	go func() {
		for e := range c {
			if l.conn == nil {
				continue
			}

			evt := jsonEvent{
				ID:              strconv.Itoa(e.ID),
				Time:            e.Time,
				RawMessage:      e.OriginalString,
				FriendlyMessage: e.String(),
			}
			b, err := json.Marshal(evt)
			if err != nil {
				//TODO: Log error
				continue
			}
			err = l.conn.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
	}()
	return c
}

func (l *EventLogger) HTTPHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got websocket request")
		c, err := l.upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("error upgrading websocket")
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
