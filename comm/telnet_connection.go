package comm

import (
	"fmt"
	"net"
	"time"

	"github.com/markdaws/gohome/log"
)

type TelnetConnection struct {
	conn         net.Conn
	info         TelnetConnectionInfo
	pingCallback PingCallback
	status       ConnectionStatus
	id           int
}

var id = 1

func NewTelnetConnection(i TelnetConnectionInfo) *TelnetConnection {
	c := TelnetConnection{
		status: CSNew,
		info:   i,
		id:     id,
	}
	id++
	return &c
}

func (c *TelnetConnection) Status() ConnectionStatus {
	return c.status
}

func (c *TelnetConnection) SetStatus(s ConnectionStatus) {
	c.status = s
}

func (c *TelnetConnection) SetPingCallback(cb PingCallback) {
	c.pingCallback = cb
}

func (c *TelnetConnection) PingCallback() PingCallback {
	return c.pingCallback
}

func (c *TelnetConnection) Info() ConnectionInfo {
	return c.info
}

func (c *TelnetConnection) Open() error {
	c.status = CSConnecting

	//TODO: Is this re-using the same network connection under the hood?
	log.V("%s connecting", c)
	conn, err := net.Dial(c.info.Network, c.info.Address)
	if err != nil {
		log.V("%s connection failed %s", c, err)
		c.status = CSClosed
		return err
	}

	c.Conn = conn

	if c.info.Authenticator != nil {
		if err = c.info.Authenticator.Authenticate(c); err != nil {
			log.V("%s authenticate failed %s", c, err)
			c.Close()
			return err
		}
	}

	log.V("%s connected successfully", c)
	c.status = CSConnected
	return nil
}

func (c *TelnetConnection) Close() {
	log.V("%s closed", c)
	c.status = CSClosed
	c.Conn.Close()
}

func (c *TelnetConnection) Read(p []byte) (n int, err error) {
	c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	n, err = c.Conn.Read(p)
	if err != nil {
		c.status = CSClosed
	}
	return
}

func (c *TelnetConnection) Write(p []byte) (n int, err error) {
	c.Conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	n, err = c.Conn.Write(p)
	if err != nil {
		c.status = CSClosed
	}
	return
}

func (c *TelnetConnection) String() string {
	return fmt.Sprintf("TelnetConnection[%d %s %s]", c.id, c.info.Network, c.info.Address)
}
