package gossdp

import (
    "net"
    "sync"
    "strings"
)


// Listener to recieve events.
type ClientListener interface {
    // Notified on M-SEARCH responses.
    Response(message ResponseMessage)
}


type ClientSsdp struct {
    socket                  *net.UDPConn
    listener                ClientListener
    writeChannel            chan writeMessage
    exitWriteWaitGroup      sync.WaitGroup
    exitReadWaitGroup       sync.WaitGroup
    interactionLock         sync.Mutex
    isRunning               bool
    logger                  LoggerInterface
}


// Creates a new client
// the client will not listen for NOTIFY events.
// it will only listen for direct replies.
//
// Why?
// 1. We are not binding to :1900. I have found that unicast replies
//      to :1900 get eaten by other processes. So the best method
//      is for the client to bind to a random port.
//      A future improvment would be to bind both to the random port
//      and :1900 so we can listen for broadcasts, and get replies.
func NewSsdpClient(l ClientListener) (*ClientSsdp, error) {
    return NewSsdpClientWithLogger(l, DefaultLogger{})
}

func NewSsdpClientWithLogger(l ClientListener, lg LoggerInterface) (*ClientSsdp, error) {
    var c ClientSsdp
    c.listener = l
    c.writeChannel = make(chan writeMessage)
    c.logger = lg
    if err := c.createSocket(); err != nil {
        return nil, err
    }
    c.isRunning = true

    return &c, nil
}

func (c *ClientSsdp) createSocket() error {
    addr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
    if err != nil {
        return err
    }
    c.socket, err = net.ListenUDP("udp4", addr)
    if err != nil {
        return err
    }
    return nil
}

func (c *ClientSsdp) parseMessage(message, hostPort string) {
    if strings.HasPrefix(message, "HTTP") {
        respData := parseResponse(message, hostPort)
        if respData != nil && c.listener != nil {
            c.listener.Response(*respData)
        }
        return
    }
    c.logger.Warnf("Unknown message. We only expect replies.")
    return
}

// Starts listening to packets on the network.
func (c *ClientSsdp) Start() {
    go c.socketWriter()
    c.socketReader()
}

func (c *ClientSsdp) socketReader() {
    c.exitReadWaitGroup.Add(1)
    defer c.exitReadWaitGroup.Add(-1)
    readBytes := make([]byte, 2048)
    for {
        n, src, err := c.socket.ReadFrom(readBytes)
        if err != nil {
            c.logger.Warnf("Error reading from socket: %v", err)
            return
        }
        if n > 0 {
            c.parseMessage(string(readBytes[0:n]), src.String())
        }
    }
}

func (c *ClientSsdp) socketWriter() {
    c.exitWriteWaitGroup.Add(1)
    defer c.exitWriteWaitGroup.Add(-1)
    for {
        msg, more := <- c.writeChannel
        if !more {
            return
        }
        _, err := c.socket.WriteTo(msg.message, msg.to)
        if err != nil {
            c.logger.Warnf("Error sending message. %v", err)
        }
    }
}

// Kills the client by closing the socket.
// If any servers are being advertised they will NOTIFY a byebye
func (c *ClientSsdp) Stop() {
    c.interactionLock.Lock()
    c.isRunning = false
    c.interactionLock.Unlock()

    if c.socket != nil {
        close(c.writeChannel)
        c.exitWriteWaitGroup.Wait()
        c.socket.Close()
        c.exitReadWaitGroup.Wait()
        c.socket = nil
    }
    c.logger.Tracef("Stop exiting")
}


// Sends out 1 M-SEARCH request for the specified target.
func (c *ClientSsdp) ListenFor(searchTarget string) error {
    msg := createSsdpHeader(
        "M-SEARCH",
        map[string]string{
            "HOST": "239.255.255.250:1900",
            "ST": searchTarget,
            "MAN": `"ssdp:discover"`,
            "MX": "3",
        },
        false,
    )

    addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
    if err != nil {
        return err
    }
    // run in a goroutine, because Start may not have been called yet
    // and thus s.writeChannel will block!
    go func() {
        c.interactionLock.Lock()
        defer c.interactionLock.Unlock()
        if !c.isRunning {
            return
        }
        c.writeChannel <- writeMessage{msg, addr, false}
    }()

    return err
}
