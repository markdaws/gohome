package comm

/*
import (
	"fmt"
	"net"
	"time"

	"github.com/markdaws/gohome/log"
)

type TelnetConnection struct {
	conn         net.Conn
	pingCallback PingCallback
	status       ConnectionStatus
	id           int
	addr         string
	auth         *TelnetAuthenticator
}

var id = 1

func NewTelnetConnection(addr string, auth *TelnetAuthenticator) *TelnetConnection {
	c := TelnetConnection{
		status: CSNew,
		addr:   addr,
		auth:   auth,
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

func (c *TelnetConnection) Open() error {
	c.status = CSConnecting

	//TODO: Is this re-using the same network connection under the hood?
	log.V("%s connecting", c)
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.V("%s connection failed %s", c, err)
		c.status = CSClosed
		return err
	}

	c.conn = conn

	if c.auth != nil {
		if err = c.auth.Authenticate(c); err != nil {
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
	c.conn.Close()
}

func (c *TelnetConnection) Read(p []byte) (n int, err error) {
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	n, err = c.conn.Read(p)
	if err != nil {
		c.status = CSClosed
	}
	return
}

func (c *TelnetConnection) Write(p []byte) (n int, err error) {
	c.conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
	n, err = c.conn.Write(p)
	if err != nil {
		c.status = CSClosed
	}
	return
}

func (c *TelnetConnection) String() string {
	return fmt.Sprintf("TelnetConnection[%d %s]", c.id, c.addr)
}
*/
