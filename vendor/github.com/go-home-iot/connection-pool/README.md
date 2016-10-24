# connection-pool
A golang based connection pool for net.Conn connections

##Documentation
See [godoc](https://godoc.org/github.com/go-home-iot/connection-pool)

##Installation
```bash
go get github.com/go-home-iot/connection-pool
```

##Package
```go
import "github.com/go-home-iot/pool"
```

##Usage
See connection_pool_test.go for examples of how to use this library.  The basic concepts are that you create a pool, in the config you pass in a NewConnection function that allows the pool to create new connections.  If at any point a connection is found to be bad, you just set the connections IsBad field to true and return it to the pool, this indicates to the pool that the connection should be thrown away and a new one created in its place.

##Version History
###0.1.0
Initial release



