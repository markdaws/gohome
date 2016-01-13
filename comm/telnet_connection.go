package comm

import (
	"fmt"
	"net"
	"time"
)

type telnetConnection struct {
	conn         net.Conn
	info         ConnectionInfo
	pingCallback PingCallback
	status       ConnectionStatus
	id           int
}

var id = 1

func NewTelnetConnection(i ConnectionInfo) *telnetConnection {
	c := telnetConnection{
		status: CSNew,
		info:   i,
		id:     id,
	}
	id++
	return &c
}

func (c *telnetConnection) Status() ConnectionStatus {
	return c.status
}

func (c *telnetConnection) SetStatus(s ConnectionStatus) {
	c.status = s
}

func (c *telnetConnection) SetPingCallback(cb PingCallback) {
	c.pingCallback = cb
}

func (c *telnetConnection) PingCallback() PingCallback {
	return c.pingCallback
}

func (c *telnetConnection) Info() ConnectionInfo {
	return c.info
}

func (c *telnetConnection) Open() error {
	c.status = CSConnecting

	fmt.Println("telnetConnection - trying to connect")
	conn, err := net.Dial(c.info.Network, c.info.Address)
	if err != nil {
		fmt.Printf("Dial failed\n")
		c.status = CSClosed
		return err
	}

	c.conn = conn

	if c.info.Authenticator != nil {
		if err = c.info.Authenticator.Authenticate(c); err != nil {
			fmt.Println("Authentication failed")
			c.Close()
			return err
		}
	}

	fmt.Println("telnetConnection - connected")
	c.status = CSConnected
	return nil
}

func (c *telnetConnection) Close() {
	c.status = CSClosed
	c.conn.Close()
}

func (c *telnetConnection) Read(p []byte) (n int, err error) {
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	n, err = c.conn.Read(p)
	if err != nil {
		c.status = CSClosed
	}
	return
}

func (c *telnetConnection) Write(p []byte) (n int, err error) {
	c.conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	n, err = c.conn.Write(p)
	if err != nil {
		c.status = CSClosed
	}
	return
}

func (c *telnetConnection) String() string {
	return fmt.Sprintf("telnetConnection[%d]", c.id)
}
