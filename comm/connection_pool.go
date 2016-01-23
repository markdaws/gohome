package comm

import (
	"fmt"
	"time"

	"github.com/markdaws/gohome/log"
)

type ConnectionPool interface {
	Name() string
	Get() Connection
	Release(Connection)
	SetNewConnectionCallback(func())
}

func NewConnectionPool(name string, count int, newConnectionCb func() Connection) ConnectionPool {
	p := &connectionPool{
		name:          name,
		pool:          make(map[Connection]bool, count),
		newConnection: newConnectionCb,
	}

	//TODO: Change to connect in parallel even if sync
	for i := 0; i < count; i++ {
		retryNewConnection(p, true)
	}
	return p
}

type connectionPool struct {
	name            string
	pool            map[Connection]bool
	newConnection   func() Connection
	newConnectionCb func()
}

func (p *connectionPool) Name() string {
	return p.name
}

func (p *connectionPool) Get() Connection {
	if len(p.pool) == 0 {
		return nil
	}
	for c := range p.pool {
		if p.pool[c] {
			p.pool[c] = false
			return c
		}
	}
	return nil
}

func (p *connectionPool) Release(c Connection) {
	if _, ok := p.pool[c]; !ok {
		return
	}

	if c.Status() == CSClosed {
		log.V("%s closed connection release, adding new connection", p)
		delete(p.pool, c)
		retryNewConnection(p, false)
		return
	}

	p.pool[c] = true
}

func (p *connectionPool) SetNewConnectionCallback(cb func()) {
	p.newConnectionCb = cb
}

func (p *connectionPool) String() string {
	return fmt.Sprintf("connectionPool[%s]", p.name)
}

func retryNewConnection(p *connectionPool, sync bool) {
	//TODO: This doesn't have to be async now we queue up commands???
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

	//TODO: sync issue, updating on different go routine
	p.pool[c] = true
	if p.newConnectionCb != nil {
		p.newConnectionCb()
	}

	// If the connection has a ping handler, the pool will call it to make
	// sure that the connection stays alive
	pingCB := c.PingCallback()
	if pingCB != nil {
		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer func() {
				ticker.Stop()
			}()

			for _ = range ticker.C {
				if err := pingCB(); err != nil {
					c.Close()
					c.SetStatus(CSClosed)

					log.V("%s ping failed for connection %s, releasing", p, c)

					// If the connection is sitting in the pool, release
					// it and open a new one
					if p.pool[c] {
						p.Release(c)
					}
					break
				}
			}
		}()
	}
	return nil
}
