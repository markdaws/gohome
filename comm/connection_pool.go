package comm

/*
import (
	"fmt"
	"time"

	"github.com/markdaws/gohome/log"
)

// Contains pool configuration information
type ConnectionPoolConfig struct {
	Name           string
	Size           int
	ConnectionType string
	Address        string
	TelnetPingCmd  string
	TelnetAuth     *TelnetAuthenticator
}

type ConnectionPool interface {
	Name() string
	Init() error
	Get() Connection
	Release(Connection)
	SetNewConnectionCallback(func())
	Config() ConnectionPoolConfig
}

//TODO: Split New and initializing connections
func NewConnectionPool(config ConnectionPoolConfig) (ConnectionPool, error) {

	var newConnection func() Connection
	switch config.ConnectionType {
	case "telnet":
		newConnection = func() Connection {
			conn := NewTelnetConnection(config.Address, config.TelnetAuth)

			if config.TelnetPingCmd != "" {
				conn.SetPingCallback(func() error {
					if _, err := conn.Write([]byte(config.TelnetPingCmd)); err != nil {
						return fmt.Errorf("%s ping failed: %s", config.Address, err)
					}
					return nil
				})
			}
			return conn
		}
	default:
		return nil, fmt.Errorf("unknown connection type: %s", config.ConnectionType)
	}

	p := &connectionPool{
		name:          config.Name,
		config:        config,
		pool:          make(map[Connection]bool, config.Size),
		newConnection: newConnection,
	}
	return p, nil
}

type connectionPool struct {
	name            string
	config          ConnectionPoolConfig
	pool            map[Connection]bool
	newConnection   func() Connection
	newConnectionCb func()
}

func (p *connectionPool) Name() string {
	return p.name
}

func (p *connectionPool) Config() ConnectionPoolConfig {
	return p.config
}

func (p *connectionPool) Init() error {
	for i := 0; i < p.config.Size; i++ {
		retryNewConnection(p, true)
	}

	//TODO: Errors
	return nil
}

//TODO: Have concept of maxTimeout waiting for a connection
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
*/
