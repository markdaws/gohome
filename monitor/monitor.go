package monitor

/*
Thoughts:

/api/v1/monitor -> POST
- list of zone ids, monitor type?
- Can a zone have more than one value you can monitor?

Monitor:
- light: intensity
- shade: level
- switch: on/off
- sensor: on/off, temperature -> name/value/unit_of_measure
- camera

Cache values for quick retrieval
Have a refresh rate to refresh values
Pull or push


send a list of ids you are interested in

1.
[
    {
      zoneId: 1234,
      parameters: ['value']  //optional, if ommited returns all values
    }, ...
]

2.
return a subscription id to the caller or hook this up to a push event model

3. As soon as we get values or changes, we push that to the client



Events
 - how relate to command
 - what are events useful for?

New Command to get levels

ZoneSetLevel
ZoneGetLevel()

SensorGetAttributes()
*/

type Monitor struct {
}

func (m *Monitor) Values() {
	// return what we have, fetch ?
}

// returns subscription id - if after a period of time haven't renewed your subscription
// we eject the subscription from book keeping
func (m *Monitor) Subscribe() string {
}

//Unsubscribe
