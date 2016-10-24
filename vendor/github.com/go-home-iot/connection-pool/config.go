package pool

import (
	"net"
	"time"
)

// Config contains all of the configuration parameters for the connection pool
type Config struct {
	// Name is a friendly name associated with the pool, cab be useful for debugging
	Name          string
	
	// Size is the number of connections to open
	Size          int
	
	// RetryDuration specifies how long the pool will wait to try to create a new connection
	// if the previous new conneciton attempt failed
	RetryDuration time.Duration
	
	// NewConnection takes in the pool config information and returns an open net.Conn connection
	// that can be added to the pool.
	NewConnection func(Config) (net.Conn, error)
}
