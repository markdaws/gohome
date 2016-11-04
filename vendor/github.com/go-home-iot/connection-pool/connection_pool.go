package pool

import (
	"errors"
	"io/ioutil"
	"sync"
	"time"
)

// ErrTimeout represents a timeout error, for example you called Get and couldn't get
// a connection within the timeout period.
var ErrTimeout = errors.New("timeout")

// ConnectionPool provides the ability to pool connections
type ConnectionPool struct {
	Config Config
	pool   chan *Connection
	closed bool
}

// NewPool creates a new ConnectionPool.  The pool which is returned will still need to
// have Init() called in it before it can be used
func NewPool(config Config) *ConnectionPool {
	p := &ConnectionPool{
		Config: config,
		pool:   make(chan *Connection, config.Size),
	}
	return p
}

// Init should be called before using the pool, the call is non blocking, but you
// can wait on the returned channel if you want to know when all of the underlying
// connections have been created and are ready to use
func (p *ConnectionPool) Init() chan bool {

	done := make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(p.Config.Size)

	for i := 0; i < p.Config.Size; i++ {
		p.retryNewConnection(&wg)
	}

	// Return the channel to let the caller know when init has completed
	go func() {
		wg.Wait()
		done <- true
	}()
	return done
}

// Close closes all of the underlying connections, this is non blocking but you can
// wait on the returned channel if you need to know all the connections have closed
func (p *ConnectionPool) Close() chan bool {
	done := make(chan bool)
	go func() {
		for len(p.pool) > 0 {
			c := <-p.pool
			c.returnOnClose = false
			c.Close()
		}
		done <- true
	}()
	return done
}

// Get is a blocking function that waits to get an available connection.  If after the
// timeout duration a connection could not be fetched, the function returns with ErrTimeout.
// The flush parameter if set to true will read all of the outstanding data from the
// connection before returning it to the caller. Note there is a possible 100ms delay for this
// function to return if you set flush==true while the pool tries to read any existing content
// from the connection
func (p *ConnectionPool) Get(timeout time.Duration, flush bool) (*Connection, error) {
	expire := time.Now().Add(timeout)
	select {
	case conn := <-p.pool:
		if flush {
			// Read all the contents from the buffer, if there is any, then
			// reset the read deadline to infinity
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			_, _ = ioutil.ReadAll(conn)
			conn.SetReadDeadline(time.Time{})
		}
		return conn, nil

	case <-time.After(expire.Sub(time.Now())):
		return nil, ErrTimeout
	}
}

// Release returns the connection back to the pool. err is any error that was returned
// by the connection while it was being used, if there was an error the pool will then
// throw this connection away and create a new one
func (p *ConnectionPool) Release(c *Connection, err error) {
	if c == nil {
		return
	}

	if err != nil {
		p.retryNewConnection(nil)
		return
	}
	p.pool <- c
}

func (p *ConnectionPool) retryNewConnection(wg *sync.WaitGroup) {
	// Just keeps trying to open a new connection until it succeeds
	go func() {
		for !p.closed {
			c, err := p.Config.NewConnection(p.Config)
			if err == nil {
				p.pool <- NewConnection(c, p)
				if wg != nil {
					wg.Done()
				}
				return
			}

			// Wait for a small time then retry
			time.Sleep(p.Config.RetryDuration)
		}
	}()
}
