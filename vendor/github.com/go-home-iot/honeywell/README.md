# honeywell
golang library for monitoring and controlling Honeywell WIFI thermostat products that use the https://mytotalconnectcomfort.com portal.

##Documentation
See [godoc](https://godoc.org/github.com/go-home-iot/honeywell)

##Installation
```bash
go get github.com/go-home-iot/honeywell
```

##Using the library
To use the library, see the examples in honeywell_test.go  You will need to know your device ID for the device you wish to control. The deviceID has to be determined manually, log in to the mytotalconnectcomfort website,navigate to your device, then the URL will look something like https://mytotalconnectcomfort.com/portal/Device/CheckDataSession/123456, you need to copy the number that is in place of the 123456 and use that as your device ID. 

Once you have your device ID you can create a thermostat instance, call Connect() first, then call the other methods, e.g. if we want to temporarily set the temperature to 68 degrees for 30 minutes we would do:

```go
ctx := context.TODO()
ts := honeywell.NewThermostat(DEVICEID)
err := ts.Connect(ctx, LOGIN, PASSWORD)
if err != nil {
  //handle error
}

err = ts.HeatMode(ctx, 68.0, time.Minute*30)
if err != nil {
  //handle error
}
```

##Version History
###0.1.0
Initial release, support for Heat/Cool/Cancel and getting the current state

