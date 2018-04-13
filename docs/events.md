Events are the backbone of the app.  The app has an Event Bus, comprising of consumer and producers.  Producers push events on to the event bus, and consumers read the events and potentially perform some action.

## Well Known Events
Extensions can push whatever events they want on the bus.  An event type simply has to implement the evtbus.Event interface found in the github.com/go-home-iot/event-bus package.  As well as custom events, there are common events:

### SensorAttrChangedEvt
This event is raised when a sensor attribute value changed.  For example, if a temperature sensor detects a temperate change it can raise this event on to the bus
```go
type SensorAttrChangedEvt struct {
  // SensorID is the ID of the sensor whos values have changed
  SensorID string
  
  // SensorName is the name of the sensor whos values have changed
  SensorName string
  
  // Information on the attribute that changed
  Attr SensorAttr
}
```

### ZoneLevelChangedEvt
This event is raised when a zones level has changed
```go
type ZoneLevelChangedEvt struct {
// ZoneID is the ID of the zone whos value has changed
ZoneID   string

// ZoneName is the name of the zone whos value changed
ZoneName string

// Level contains the current zone level information
Level    cmd.Level
}
```

### SensorsReportEvt
This event signifies that the system wishes to get the current sensor attribute values for all of the sensors listed in the SensorIDs field.  Consumers should look at this event and if they are responsible for any sensors in the list, get the current sensor values, then raise a SensorsReportingEvt event with the current value
```go
type SensorsReportEvt struct {
  SensorIDs map[string]bool
}
```

### SensorsReportingEvt
This event is fired when there are sensors in the system that are reporting their current values. Consumers should look at the values in the Sensors field and see if they are interested in any of the values, the field is keyed by the sensors ID.
```go
type SensorsReportingEvt struct {
  Sensors map[string]SensorAttr
}
```

### ZonesReportEvt
This event signifies that the system wishes to get the current zone leves for all of the zones listed in the ZoneIDs field.  Consumers should look at this event and if they are responsible for any zones in the list, get the current zone level, then raise a ZonesReportingEvt event with the current level
```go
type ZonesReportEvt struct {
  ZoneIDs map[string]bool
}
```

### ZonesReportingEvt
This event is fired when there are zones in the system that are reporting their current level. Consumers should look at the values in the Zones field and see if they are interested in any of the values, the field is keyed by the zones ID.
```go
type ZonesReportingEvt struct {
  Zones map[string]cmd.Level
}
```

### DeviceLostEvt
This event is raised if connection to a device is lost
```go
type DeviceLostEvt struct {
  DeviceID string
  DeviceName string
}
```

###ZoneLevelReporting
//TODO:
###SensorAttributeReporting
//TODO:
###ClientConnectedEvt
//TODO:
###ClientDisconnectedEvt
//TODO:
###DeviceProducingEvt
//TODO:
###DeviceConnectedEvt
//TODO:
###Sunrise
//TODO:
###Sunset
//TODO:
###UserLoginEvt
//TODO:
###UserLogoutEvt
//TODO:
