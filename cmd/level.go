package cmd

// Level represent the level of a zone.  It can contain a single levle value
// or hold RGB values if the zone supports that information
type Level struct {
	Value float32
	R     byte
	G     byte
	B     byte
}
