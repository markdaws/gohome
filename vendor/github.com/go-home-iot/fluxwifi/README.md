# fluxwifi
A golang library to scan and control Flux WIFI smart bulbs

##Documentation
See [godoc](https://godoc.org/github.com/go-home-iot/fluxwifi)

##Installation
```bash
go get github.com/go-home-iot/fluxwifi
```

##Package
```go
import "github.com/go-home-iot/fluwifi"
```

##Using the library
Some functions take an io.Writer interface, i.e. the TurnOn, TurnOff functions.  The io.Writer needs to be an open connection to the Flux WIFI bulb.  To open a connection, you would first call Scan() then in the responses, the IP field contains the address of the device, using that address, you can do a net.Dial to get a connection, then pass in the connection object to the functions as the io.Writer interface.

```go
conn, err := net.Dial("tcp", scanResponse.IP) 
```

##Version History
###0.1.0
Initial release, support for scanning for devices and TurnOn/TurnOff/SetValue
