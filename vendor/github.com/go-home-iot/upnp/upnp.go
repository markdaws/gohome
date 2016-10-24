package upnp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// NotifyEvent contains information about a notify event instance
type NotifyEvent struct {
	// SID is the unique ID returned to the called when they subscribed to the event
	SID string

	// Body is the response from the server, the caller is responsible for parsing it
	Body string
}

// Subscriber is an interface for types that want to get events from a device
type Subscriber interface {
	UPNPNotify(e NotifyEvent)
}

// SubServer is a upnp subscription server.  Do not create directly, you must call NewSubServer to
// return a working instance of this type
type SubServer struct {
	listenAddr string
	subs       map[string]*subInfo
}

type subInfo struct {
	SID              string
	ServiceURL       string
	Subscriber       Subscriber
	TimeoutInSeconds int
	Timer            *time.Timer
}

// NewSubServer returns a new an inited SubServer instance
func NewSubServer() *SubServer {
	s := &SubServer{}
	s.subs = make(map[string]*subInfo)
	return s
}

// Start starts a web server at the specified address which listens for upnp NOTIFY events. listenAddr is
// the address the upnp type uses to listen for notifications as the callback to the SUBSCRIBE call. e.g.
// 192.168.0.9:9000
func (s *SubServer) Start(listenAddr string) error {
	s.listenAddr = listenAddr

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Example notify response
		/*
			Method:NOTIFY URL:/ Proto:HTTP/1.1 ProtoMajor:1 ProtoMinor:1 Header:map[Nt:[upnp:event] Nts:[upnp:propchange] Sid:[uuid:054ce1c6-1dd2-11b2-864f-d1dbdfd9a2e2] Seq:[6] Content-Type:[text/xml; charset="utf-8"] Content-Length:[296]] Body:0xc4200dc300 ContentLength:296 TransferEncoding:[] Close:false Host:192.168.0.9:9000 Form:map[] PostForm:map[] MultipartForm:<nil> Trailer:map[] RemoteAddr:192.168.0.34:3663 RequestURI:/ TLS:<nil> Cancel:<nil> Response:<nil> ctx:0xc4200dc340}
		*/

		defer func() {
			w.WriteHeader(http.StatusOK)
		}()

		if r.Method != "NOTIFY" {
			return
		}

		sid, ok := r.Header["Sid"]
		if !ok {
			return
		}
		if len(sid) == 0 {
			return
		}
		sub, ok := s.subs[sid[0]]
		if !ok {
			return
		}

		// Just return a string, caller responsible for disassembling it
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		sub.Subscriber.UPNPNotify(NotifyEvent{SID: sid[0], Body: string(body)})
	})

	return http.ListenAndServe(listenAddr, mux)
}

// Subscribe signs up to receive events when there are new NOTIFY messages from the device. The function returns the SID
// returned by the device, once it has been successfully subscribed to.  You will need to keep the SID to unsubscribe.
func (s *SubServer) Subscribe(serviceURL, sid string, timeoutInSeconds int, autoRefresh bool, sub Subscriber) (string, error) {
	req, err := http.NewRequest(
		"SUBSCRIBE",
		serviceURL,
		nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("TIMEOUT", "Second-"+strconv.Itoa(timeoutInSeconds))
	req.Close = true

	if sid != "" {
		// Renew subscription
		req.Header.Add("SID", sid)
	} else {
		req.Header.Add("CALLBACK", "<http://"+s.listenAddr+">")
		req.Header.Add("NT", "upnp:event")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	sids, ok := resp.Header["Sid"]
	if !ok || len(sids) == 0 {
		return "", errors.New("Sid header missing in response")
	}

	// If the sid == "" this is a new subscribe otherwise it is a refresh in
	// which case we don't need to insert this book keeping
	if sid == "" {
		s.subs[sids[0]] = &subInfo{
			SID:              sids[0],
			ServiceURL:       serviceURL,
			Subscriber:       sub,
			TimeoutInSeconds: timeoutInSeconds,
		}
	}

	if autoRefresh {
		s.renewRefreshTimer(sids[0])
	}

	return sids[0], nil
}

// RefreshSubscription refreshes the subscription to the service.
func (s *SubServer) RefreshSubscription(sid string, autoRefresh bool) error {
	sub, ok := s.subs[sid]
	if !ok {
		return fmt.Errorf("unknown SID: %s", sid)
	}

	_, err := s.Subscribe(sub.ServiceURL, sid, sub.TimeoutInSeconds, autoRefresh, sub.Subscriber)
	return err
}

// Unsubscribe unsubscribes from the events associated with the specified sid
func (s *SubServer) Unsubscribe(sid string) error {
	sub, ok := s.subs[sid]

	if !ok {
		return nil
	}

	if sub.Timer != nil {
		sub.Timer.Stop()
	}

	delete(s.subs, sid)

	req, err := http.NewRequest(
		"UNSUBSCRIBE",
		sub.ServiceURL,
		nil)
	req.Close = true
	req.Header.Add("SID", sub.SID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error unsubscribing, status code: %d", resp.StatusCode)
	}
	return nil
}

func (s *SubServer) renewRefreshTimer(sid string) {
	sub, ok := s.subs[sid]
	if !ok {
		return
	}

	expTimer := time.AfterFunc(
		time.Duration(float32(sub.TimeoutInSeconds)*0.75)*time.Second,
		func() {
			s.RefreshSubscription(sid, true)
		},
	)
	sub.Timer = expTimer
}
