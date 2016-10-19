# connectedbytcp
A golang library for discovering and controlling ConnectedByTCP smart bulbs

##Documentation
See [godoc](https://godoc.org/github.com/go-home-iot/connectedbytcp)

##Installation
```bash
go get github.com/go-home-iot/connectedbytcp
```

##Package
```go
import "github.com/go-home-iot/connectedbytcp"
```

##Using the library
To control a bulb, you first have to call Scan() to get a list of all the ConnectedByTCP hubs.  The address of the hub will then be in ScanResponse.Location.  The next important step is that before you can talk to the bulbs, you have to press the physical "sync" button on your hub hardware, then call the GetToken() function.  This will then return a security token that must be used in all API calls.  If you don't press the physical "sync" button on the hub before calling this function you will get an error.

Once you have the address and the token, you can then call RoomGetCarousel, this returns a list of rooms and devices (the bulbs).  To control a bulb, you then need the hub address, the token and then the ID of the bulb you want to control, which was in the Device.DID field.

NOTE: it is assumed you configured the bulbs using the app provided by the hardware maker, once you have configured the bulbs you can then use this library to control them.

##Version History
###0.1.0
Initial release, support for scanning for devices and TurnOn/TurnOff/SetValue

