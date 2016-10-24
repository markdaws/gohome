# event-bus
An event bus for consuming and producing events.  The library allows for slow consumers not to block the bus and stop other consumers from receiving events, it also handles having fast producers that are overwhelming the system.

##Documentation
See [godoc](https://godoc.org/github.com/go-home-iot/event-bus)

##Installation
```bash
go get github.com/go-home-iot/event-bus
```

##Package
```go
import "github.com/go-home-iot/event-bus"
```

##Usage
See bus_test.go for examples on how to use this library.

##Version History
###0.1.0
Initial release



