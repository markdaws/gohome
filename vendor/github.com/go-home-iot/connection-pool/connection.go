package pool

import "net"

// Connection represents a connection to a network resource
type Connection struct {
	net.Conn
	owner         *ConnectionPool
	returnOnClose bool
	IsBad         bool
}

// NewConnection returns an initialized Connection instance
func NewConnection(c net.Conn, p *ConnectionPool) *Connection {
	return &Connection{
		Conn:          c,
		owner:         p,
		returnOnClose: true,
	}
}

// Close returns the connection to the pool, the connection stays open
func (c *Connection) Close() error {
	if !c.returnOnClose {
		if c.Conn != nil {
			return c.Conn.Close()
		}
		return nil
	}
	c.owner.Release(c)
	return nil
}
