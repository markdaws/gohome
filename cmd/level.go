package cmd

import "fmt"

// Level represent the level of a zone.  It can contain a single level value
// or hold RGB values if the zone supports that information
type Level struct {
	Value float32 `json:"level"`
	R     byte    `json:"r"`
	G     byte    `json:"g"`
	B     byte    `json:"b"`
}

func (l Level) String() string {
	return fmt.Sprintf("v: %f, r:%d, g:%d, b:%d", l.Value, l.R, l.G, l.B)
}
