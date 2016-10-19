/*
 * Copyright (c) 2013-2015, fromkeith
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice, this
 *   list of conditions and the following disclaimer in the documentation and/or
 *   other materials provided with the distribution.
 *
 * * Neither the name of the fromkeith nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
 * ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

/*
A very simple SSDP implementation.

Client
======
    // create the client, passing in the listener, binding the socket
    // Notice: the client does not listen to broadcasts.. see client doc
    client err := gossdp.NewSsdpClient(b)
    if err != nil {
        log.Println("Failed to start client: ", err)
        return
    }
    // call stop  when we are done
    defer c.Stop()
    // run! this will block until stop is called. so open it in a goroutine here
    go c.Start()

    // send a request for the server type we are listening for.
    err = c.ListenFor("urn:fromkeith:test:web:0")
    if err != nil {
        log.Println("Error ", err)
    }

Server
======
    // create the server, binding the socket
    s, err := gossdp.NewSsdp(nil)
    if err != nil {
        log.Println("Error creating ssdp server: ", err)
        return
    }
    // call stop  when we are done
    defer s.Stop()
    // run! this will block until stop is called. so open it in a goroutine here
    go s.Start()

    // Define the service we want to advertise
    serverDef := gossdp.AdvertisableServer{
        ServiceType: "urn:fromkeith:test:web:0",            // define the service type
        DeviceUuid: "hh0c2981-0029-44b7-4u04-27f187aecf78", // make this unique!
        Location: "http://192.168.1.1:8080",                // this is the location of the service we are advertising
        MaxAge: 3600,                                       // Max age this advertisment is valid for
    }
    // start advertising it!
    s.AdvertiseServer(serverDef)




Misc
====

USN:
    uuid:device-UUID::upnp:rootdevice
        Snet once per root device
    uuid:device-UUID
        Sent once per device. Device-UUID unique for all devices.
    uuid:device-UUID::urn:domain-name:device:deviceType:v
        Sent once per device. device-UUID, domain-name, device, deviceType and v (version)
        defined by vendor. Periods in domainname should be replaced with '-'


*/
package gossdp

import (
    "strings"
    "regexp"
    "log"
    "time"
    "net"
    "fmt"
    "bytes"
    "errors"
    "strconv"
    "net/http"
    "bufio"
    "runtime"
    "sync"
)


// a small interface to intercept all of my logs
type LoggerInterface interface {
    Tracef(fmt string, args ... interface{})
    Infof(fmt string, args ... interface{})
    Warnf(fmt string, args ... interface{})
    Errorf(fmt string, args ... interface{})
}

// a default implementation of the LoggerInterface, simply using the 'log' library
type DefaultLogger struct {}

func (l DefaultLogger) Tracef(fmt string, args ... interface{}) {
    log.Printf(fmt + "\n", args...)
}

func (l DefaultLogger) Infof(fmt string, args ... interface{}) {
    log.Printf(fmt + "\n", args...)
}
func (l DefaultLogger) Warnf(fmt string, args ... interface{}) {
    log.Printf(fmt + "\n", args...)
}
func (l DefaultLogger) Errorf(fmt string, args ... interface{}) {
    log.Printf(fmt + "\n", args...)
}


var (
    cacheControlAge = regexp.MustCompile(`.*max-age=([0-9]+).*`)
    serverName = fmt.Sprintf("%s/0.0 UPnP/1.0 gossdp/0.1", runtime.GOOS)
)

// a SSDP defintion
type Ssdp struct {
    advertisableServers     map[string][]*AdvertisableServer
    deviceIdToServer        map[string]*AdvertisableServer
    socket                  theSocket
    listener                SsdpListener
    listenSearchTargets     map[string]bool
    writeChannel            chan writeMessage
    exitWriteWaitGroup      sync.WaitGroup
    exitReadWaitGroup       sync.WaitGroup
    interactionLock         sync.Mutex
    isRunning               bool
    logger                  LoggerInterface
}

type writeMessage struct {
    message             []byte
    to                  *net.UDPAddr
    shouldExit          bool
}


// The common SSDP fields in the Notify ssdp:alive message.
//
// Notify (alive)
//      NOTIFY * HTTP/1.1
//      Host: 239.255.255.250:1900
//      NT: blenderassociation:blender               // notification type. Aka search target.
//      NTS: ssdp:alive                              // message sub-type. Either ssdp:alive or ssdp:byebye
//      USN: someunique:idscheme3                    // Unique Service Name. An instance of a device
//      LOCATION: <blender:ixl><http://foo/bar>      // location of the service being advertised. Eg. http://hello.com
//      Cache-Control: max-age = 7393                // how long this is valid for. as defined by http standards
//      SERVER: WIN/8.1 UPnP/1.0 gossdp/0.1                  // Concat of OS, UPnP, and product.
type AliveMessage struct {
    // Search Target. The urn: that defines what type of resource it is
    SearchType      string
    // Its unique identifier
    DeviceId        string
    // The USN of the service. uuid:DeviceId:SearchType
    Usn             string
    // The location of the service being advertised
    Location        string
    // How long this message should be considered valid for
    MaxAge          int
    // The os/generic info about the SSDP server
    Server          string
    // The parsed request
    RawRequest      *http.Request
    // The urn part of the USN
    Urn             string
}

// Notify (bye):
//      NOTIFY * HTTP/1.1
//      Host: 239.255.255.250:1900
//      NT: search:target
//      NTS: ssdp:byebye
//      USN: uuid:the:unique
type ByeMessage struct {
    // Search Target. The urn: that defines what type of resource it is
    SearchType      string
    // Its unique identifier
    DeviceId        string
    // The USN of the service. uuid:DeviceId:SearchType
    Usn             string
    // The parsed request
    RawRequest      *http.Request
    // The urn part of the USN
    Urn             string
}

// M-Search Response:
//      HTTP/1.1 200 OK
//      Ext:                                                 // required by http extension framework. just key, no value
//      Cache-Control: max-age = 5000                        // number of seconds this message is valid for
//      ST: ge:fridge                                        // Search target. respond with all matching targets. Same as NT in Notify messages
//      USN: uuid:abcdefgh-7dec-11d0-a765-00a0c91e6bf6       // Unique Service name
//      LOCATION: <blender:ixl><http://foo/bar>              // location of the service being advertised. Eg. http://hello.com
//      SERVER: WIN/8.1 UPnP/1.0 gossdp/0.1                  // Concat of OS, UPnP, and product.
//      DATE: date of response                               // rfc1123-date of the response
type ResponseMessage struct {
    // How long this message should be considered valid for
    MaxAge              int
    // Search Target. The urn: that defines what type of resource it is
    SearchType          string
    // Its unique identifier
    DeviceId            string
    // The USN of the service. uuid:DeviceId:SearchType
    Usn                 string
    // The location of the service being advertised
    Location            string
    // The os/generic info about the SSDP server
    Server              string
    // The parsed response
    RawResponse         *http.Response
    // The urn part of the USN
    Urn             string
}

// Listener to recieve events.
type SsdpListener interface {
    // Notified on ssdp:alive messages. Only for those we are listening for.
    NotifyAlive(message AliveMessage)
    // Notified on ssdp:byebye messages. Only for those we are listening for.
    NotifyBye(message ByeMessage)
    // Notified on M-SEARCH responses.
    Response(message ResponseMessage)
}

// reference doc: http://www.upnp.org/specs/arch/UPnP-arch-DeviceArchitecture-v1.0-20081015.pdf


// search: client-only:
// M-SEARCH * HTTP/1.1
// Host: 239.255.255.250:1900
// Man: "ssdp:discover"                                 // message sub-type
// ST: ge:fridge                                        // search target
                                                        //  ssdp:all -> all targets
                                                        //  uuid:device-UUID    -> particular target
                                                        //  urn:domainname:service:servicetype:v
// MX: 3                                                // maximum wait time in seconds.
                                                        //  Response time should be random between 0 and this number



// Describes the server/service we wish to advertise
type AdvertisableServer struct {
    // The type of this service. In the URN it is pasted after the device-UUID.
    //  It is what devices will search for
    ServiceType             string
    // The unique identifier of this device.
    DeviceUuid              string
    // The location of the service we are advertising. Eg. http://192.168.0.2:3434
    Location                string
    // The max number of seconds we want advertise and responses to be valid for.
    MaxAge                  int

    usn                     string
    lastTimer               *time.Timer
    last3sTimer             *time.Timer
}

// Register a service to advertise
// Should only be called once per server
// This implementation will automatically adverise when maxAge expires.
func (s *Ssdp) AdvertiseServer(ads AdvertisableServer) {
    s.interactionLock.Lock()
    defer s.interactionLock.Unlock()
    if !s.isRunning {
        return
    }

    adsPointer := &ads

    adsPointer.usn = fmt.Sprintf("uuid:%s::%s", adsPointer.DeviceUuid, adsPointer.ServiceType)
    if v, ok := s.advertisableServers[adsPointer.ServiceType]; ok {
        s.advertisableServers[adsPointer.ServiceType] = append(v, adsPointer)
    } else {
        s.advertisableServers[adsPointer.ServiceType] = []*AdvertisableServer{adsPointer}
    }
    s.deviceIdToServer[adsPointer.DeviceUuid] = adsPointer
    adsPointer.lastTimer = s.advertiseTimer(adsPointer, 1 * time.Second, adsPointer.MaxAge)
    adsPointer.last3sTimer = s.advertiseTimer(adsPointer, 3 * time.Second, adsPointer.MaxAge)
}


func (s *Ssdp) RemoveServer(deviceUuid string) {
    s.interactionLock.Lock()
    defer s.interactionLock.Unlock()
    if !s.isRunning {
        return
    }


    var ads *AdvertisableServer
    var ok bool
    if ads, ok = s.deviceIdToServer[deviceUuid]; !ok {
        return
    }
    ads.lastTimer.Stop()
    ads.last3sTimer.Stop()
    delete(s.deviceIdToServer, deviceUuid)
    var group []*AdvertisableServer
    if group, ok = s.advertisableServers[ads.ServiceType]; !ok {
        return
    }
    if len(group) == 1 {
        delete(s.advertisableServers, ads.ServiceType)
        return
    }
    for i := range group {
        if group[i].DeviceUuid == ads.DeviceUuid {
            newGroup := make([]*AdvertisableServer, len(group) - 1)
            if i > 0 {
                copy(newGroup, group[:i])
            }
            if i < len(group) - 1 {
                copy(newGroup[i:], group[i+1:len(group)])
            }
            s.advertisableServers[ads.ServiceType] = newGroup
            break
        }
    }
}


// Creates a new server
func NewSsdp(l SsdpListener) (*Ssdp, error) {
    return NewSsdpWithLogger(l, DefaultLogger{})
}

func NewSsdpWithLogger(l SsdpListener, lg LoggerInterface) (*Ssdp, error) {
    var s Ssdp
    s.advertisableServers = make(map[string][]*AdvertisableServer)
    s.deviceIdToServer = make(map[string]*AdvertisableServer)
    s.listenSearchTargets = make(map[string]bool)
    s.listener = l
    s.writeChannel = make(chan writeMessage)
    s.logger = lg
    if err := s.createSocket(); err != nil {
        return nil, err
    }
    s.isRunning = true

    return &s, nil
}

func (s *Ssdp) parseMessage(message, hostPort string) {
    if strings.HasPrefix(message, "HTTP") {
        if s.listener == nil {
            return
        }
        respData := parseResponse(message, hostPort)
        if respData != nil {
            s.listener.Response(*respData)
        }
        return
    }
    req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(message)))
    if err != nil {
        s.logger.Warnf("Error reading request: ", err)
        return
    }

    if req.URL.Path != "*" {
        s.logger.Warnf("Unknown path requested: ", req.URL.Path)
        return
    }

    s.parseCommand(req, hostPort)
}

func (s *Ssdp) parseCommand(req * http.Request, hostPort string) {
    if req.Method == "NOTIFY" {
        s.notify(req)
        return
    }
    if req.Method == "M-SEARCH" {
        s.msearch(req, hostPort)
        return
    }
    s.logger.Warnf("Unknown message type!. Message: " + req.Method)
}

func extractUrnDeviceIdFromUsn(usn string) (deviceId, urn string) {
    if len(usn) > 0 {
        parts := strings.Split(usn, ":")
        if len(parts) > 2 {
            if parts[0] == "uuid" {
                deviceId = parts[1]
                urn = strings.TrimPrefix(usn, "uuid:" + deviceId + ":")
                if parts[2] == "" {
                    urn = strings.TrimPrefix(urn, ":")
                }
            } else {
                urn = usn
            }
        }
    }
    return
}

func (s *Ssdp) notify(req * http.Request) {
    if s.listener == nil {
        return
    }
    nts := req.Header.Get("NTS")
    if nts == "" {
        s.logger.Warnf("Missing NTS in NOTIFY")
        return
    }
    searchType := req.Header.Get("NT")
    if searchType == "" {
        s.logger.Warnf("Missing NT in NOTIFY")
        return
    }
    usn := req.Header.Get("USN")
    deviceId, urn := extractUrnDeviceIdFromUsn(usn)

    nts = strings.ToLower(nts)
    if nts == "ssdp:alive" {
        location := req.Header.Get("LOCATION")
        server := req.Header.Get("SERVER")
        maxAge := -1
        if cc := req.Header.Get("CACHE-CONTROL"); cc != "" {
            subMatch := cacheControlAge.FindStringSubmatch(cc)
            if len(subMatch) == 2 {
                maxAgeInt64, err := strconv.ParseInt(subMatch[1], 10, 0)
                if err == nil {
                    maxAge = int(maxAgeInt64)
                }
            }
        }
        // don't notify alive for people we aren't listening to
        if len(s.listenSearchTargets) > 0 {
            if _, ok := s.listenSearchTargets[urn]; !ok {
                return
            }
        }
        message := AliveMessage{
            SearchType      : searchType,
            DeviceId        : deviceId,
            Usn             : usn,
            Urn             : urn,
            Location        : location,
            MaxAge          : maxAge,
            Server          : server,
            RawRequest      : req,
        }
        s.listener.NotifyAlive(message)
        return
    }
    if nts == "ssdp:byebye" {
        message := ByeMessage{
            SearchType      : searchType,
            Usn             : usn,
            Urn             : urn,
            DeviceId        : deviceId,
            RawRequest      : req,
        }
        s.listener.NotifyBye(message)
        return
    }
    s.logger.Warnf("Could not identify NTS header!: " + nts)
}


func (s *Ssdp) msearch(req * http.Request, hostPort string) {
    if v := req.Header.Get("MAN"); v == "" {
        return
    }
    if v := req.Header.Get("MX"); v == "" {
        return
    }
    if st := req.Header.Get("ST"); st == "" {
        return
    } else {
        s.inMSearch(st, req, hostPort) // TODO: extract MX
    }
}


func parseResponse(msg, hostPort string) (*ResponseMessage) {
    resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(msg)), nil)
    if err != nil {
        return nil
    }
    defer resp.Body.Close()

    maxAge := -1
    if cc := resp.Header.Get("CACHE-CONTROL"); cc != "" {
        subMatch := cacheControlAge.FindStringSubmatch(cc)
        if len(subMatch) == 2 {
            maxAgeInt64, err := strconv.ParseInt(subMatch[1], 10, 0)
            if err == nil {
                maxAge = int(maxAgeInt64)
            }
        }
    }
    usn := resp.Header.Get("USN")
    deviceId, urn := extractUrnDeviceIdFromUsn(usn)

    respMessage := ResponseMessage{
        MaxAge              : maxAge,
        SearchType          : resp.Header.Get("ST"),
        Usn                 : usn,
        Urn                 : urn,
        DeviceId            : deviceId,
        Location            : resp.Header.Get("LOCATION"),
        Server              : resp.Header.Get("SERVER"),
        RawResponse         : resp,
    }
    return &respMessage
}


func (s *Ssdp) inMSearch(st string, req * http.Request, sendTo string) {
    if st[0] == '"' && st[len(st) - 1] == '"' {
        st = st[1:len(st) - 2]
    }
    mx := 1
    if mxStr := req.Header.Get("MX"); mxStr != "" {
        mxInt64, err := strconv.ParseInt(mxStr, 10, 0)
        if err != nil {
            mx = int(mxInt64)
        }
        if mx < 1 {
            mx = 1
        } else if mx > 5 {
            mx = 5
        }
    }

    if st == "ssdp:all" {
        for _, v := range s.advertisableServers {
            for _, d := range v {
                s.respondToMSearch(d, sendTo, mx)
            }
        }
    } else if d, ok := s.deviceIdToServer[st]; ok {
        s.respondToMSearch(d, sendTo, mx)
    } else if v, ok := s.advertisableServers[st]; ok {
        for _, d := range v {
            s.respondToMSearch(d, sendTo, mx)
        }
    }
}

func (s *Ssdp) respondToMSearch(ads *AdvertisableServer, sendTo string, mx int) {
    time.Sleep(time.Duration(mx) * time.Second)


    msg := createSsdpHeader(
        "200 OK",
        map[string]string{
            "ST": ads.ServiceType,
            "USN": ads.usn,
            "LOCATION": ads.Location,
            "CACHE-CONTROL": fmt.Sprintf("max-age=%d", ads.MaxAge),
            "DATE": time.Now().Format(time.RFC1123),
            "SERVER": serverName,
            "EXT": "",
        },
        true,
    )

    addr, err := net.ResolveUDPAddr("udp4", sendTo)
    if err != nil {
        s.logger.Errorf("Error resolving UDP addr: ", err)
        return
    }

    s.interactionLock.Lock()
    defer s.interactionLock.Unlock()
    if !s.isRunning  {
        return
    }

    s.writeChannel <- writeMessage{msg, addr, false}
}

// Filters the NOTIFIES to only be returned for the given target.
func (s *Ssdp) ListenFor(searchTarget string) error {
    s.interactionLock.Lock()
    defer s.interactionLock.Unlock()
    if !s.isRunning {
        return errors.New("Not running. Can't listen")
    }

    // listen directly for their search target
    s.listenSearchTargets[searchTarget] = true

    return nil
}


func (s *Ssdp) advertiseTimer(ads *AdvertisableServer, d time.Duration, age int) *time.Timer {
    var timer *time.Timer
    timer = time.AfterFunc(d, func () {
        s.advertiseServer(ads, true)
        timer.Reset(d + time.Duration(age) * time.Second)
    })
    return timer
}


// Kills the server by closing the socket.
// If any servers are being advertised they will NOTIFY a byebye
func (s *Ssdp) Stop() {
    s.interactionLock.Lock()
    s.isRunning = false
    s.interactionLock.Unlock()

    if s.socket.IsValid() {
        if len(s.advertisableServers) > 0 {
            s.advertiseClosed()
        }
        s.writeChannel <- writeMessage{nil, nil, true}
        s.exitWriteWaitGroup.Wait()
        close(s.writeChannel)
        //s.socket.Close()
        s.closeSocket()
        //s.rawSocket.Close()
        s.exitReadWaitGroup.Wait()
        //s.rawSocket = nil
    }
    s.logger.Tracef("Stop exiting")
}

func (s *Ssdp) advertiseClosed() {
    for _, ad := range s.deviceIdToServer {
        ad.lastTimer.Stop()
        ad.last3sTimer.Stop()
        s.advertiseServer(ad, false)
    }
}

func (s *Ssdp) advertiseServer(ads *AdvertisableServer, alive bool) {
    s.interactionLock.Lock()
    defer s.interactionLock.Unlock()

    if !s.isRunning  {
        return
    }

    ntsString := "ssdp:alive"
    if !alive {
        ntsString = "ssdp:byebye"
    }

    heads := map[string]string{
        "HOST": "239.255.255.250:1900",
        "NT": ads.ServiceType,
        "NTS": ntsString,
        "USN": ads.usn,
    }
    if alive {
        heads["LOCATION"] = ads.Location
        heads["CACHE-CONTROL"] = fmt.Sprintf("max-age=%d", ads.MaxAge)
        heads["SERVER"] = serverName
    }
    msg := createSsdpHeader(
            "NOTIFY",
            heads,
            false,
        )

    to, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
    if err == nil {
        s.writeChannel <- writeMessage{msg, to, false}
    } else {
        s.logger.Warnf("Error sending advertisement: ", err)
    }
}

func createSsdpHeader(head string, vars map[string]string, isResponse bool) []byte {
    buf := bytes.Buffer{}
    if isResponse {
        buf.WriteString(fmt.Sprintf("HTTP/1.1 %s\r\n", head))
    } else {
        buf.WriteString(fmt.Sprintf("%s * HTTP/1.1\r\n", head))
    }
    for k, v := range vars {
        buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
    }
    buf.WriteString("\r\n")
    return []byte(buf.String())
}

// Starts listening to packets on the network.
func (s *Ssdp) Start() {
    go s.socketWriter()
    s.socketReader()
}


func (s *Ssdp) socketReader() {
    s.exitReadWaitGroup.Add(1)
    defer s.exitReadWaitGroup.Add(-1)

    for {
        msg, src, err := s.read()
        if err != nil {
            s.logger.Warnf("Error reading from SSDP socket: %v", err)
            return
        }
        if len(msg) > 0 {
            //s.logger.Warnf("Received: %s", string(msg))
            s.parseMessage(string(msg), src)
            //s.logger.Warnf("Done parsing")
        }
    }
}

func (s *Ssdp) socketWriter() {
    s.exitWriteWaitGroup.Add(1)
    defer s.exitWriteWaitGroup.Add(-1)
    for {
        msg, more := <- s.writeChannel
        if !more {
            return
        }
        if msg.shouldExit {
            return
        }
        if err := s.write(msg); err != nil {
            s.logger.Warnf("Error sending message. %v", err)
        }
    }
}
