package comm

import "time"

type ConnectionPool interface {
	Get() Connection
	Release(Connection)
	SetNewConnectionCallback(func())
}

func NewConnectionPool(count int, newConnectionCb func() Connection) ConnectionPool {
	p := &connectionPool{
		pool:          make(map[Connection]bool, count),
		newConnection: newConnectionCb,
	}

	for i := 0; i < count; i++ {
		retryNewConnection(p, true)
	}
	return p
}

type connectionPool struct {
	pool            map[Connection]bool
	newConnection   func() Connection
	newConnectionCb func()
}

func (p *connectionPool) Get() Connection {
	if len(p.pool) == 0 {
		return nil
	}
	for c, _ := range p.pool {
		if p.pool[c] {
			p.pool[c] = false
			return c
		}
	}
	return nil
}

func (p *connectionPool) Release(c Connection) {
	if c.Status() == CSClosed {
		retryNewConnection(p, false)
		return
	}
	p.pool[c] = true
}

func (p *connectionPool) SetNewConnectionCallback(cb func()) {
	p.newConnectionCb = cb
}

func retryNewConnection(p *connectionPool, sync bool) {
	c := make(chan bool)
	go func() {
		for {
			err := appendAndOpenNewConnection(p)
			if err == nil {
				break
			}
			time.Sleep(time.Second * 30)
		}
		close(c)
	}()

	if sync {
		<-c
	}
}

func appendAndOpenNewConnection(p *connectionPool) error {
	c := p.newConnection()
	c.SetStatus(CSNew)
	err := c.Open()
	if err != nil {
		return err
	}

	p.pool[c] = true
	if p.newConnectionCb != nil {
		p.newConnectionCb()
	}

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer func() {
			ticker.Stop()
		}()

		for _ = range ticker.C {
			if err := c.PingCallback()(); err != nil {
				c.Close()
				c.SetStatus(CSClosed)

				// If the connection is sitting in the pool, release
				// it and open a new one
				if p.pool[c] {
					p.Release(c)
				}
				break
			}
		}
	}()

	return nil
}
