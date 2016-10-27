Events are the backbone of the app.  The app has an Event Bus, comprising of consumer and producers.  Producers push events on to the event bus, and consumers read the events and potentially perform some action.

##Well Known Events
Extensions can push whatever events they want on the bus.  An event type simply has to implement the evtbus.Event interface found in the github.com/go-home-iot/event-bus package.  As well as custom events, there are common events:

###SensorAttrChanged
This event is raised when a sensor attribute value changed.  For example, if a temperature sensor detects a temperate change it can raise this event on to the bus
```go
type SensorAttrChanged struct {
  // SensorID is the ID of the sensor whos values have changed
  SensorID string
  
  // SensorName is the name of the sensor whos values have changed
  SensorName string
  
  // Information on the attribute that changed
  Attr SensorAttr
}
```
###ZoneLevelChanged
This event is raised when a zones level has changed
```go
type ZoneLevelChanged struct {
// ZoneID is the ID of the zone whos value has changed
ZoneID   string

// ZoneName is the name of the zone whos value changed
ZoneName string

// Level contains the current zone level information
Level    cmd.Level
}
```
###SensorsReport
This event signifies that the system wishes to get the current sensor attribute values for all of the sensors listed in the SensorIDs field.  Consumers should look at this event and if they are responsible for any sensors in the list, get the current sensor values, then raise a SensorsReporting event with the current value
```go
type SensorsReport struct {
  SensorIDs map[string]bool
}
```
###SensorsReporting
//TODO:
###ZonesReport
//TODO:
###ZonesReporting
//TODO:
###Sunrise
//TODO:
###Sunset
//TODO:
