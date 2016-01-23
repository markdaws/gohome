package gohome

// Bulb represents a light bulb. Bulbs are owned by zones. Bulbs are not
// used to set and get values, but might contain information such as power
// usage for a particular bulb that can then be used for energy calculations etc.
type Bulb struct {
}

//TODO: Zones own bulbs, bulbs are not directly controllable, but
//can contain info such as type, power usage etc that can be used
//for other calculations
