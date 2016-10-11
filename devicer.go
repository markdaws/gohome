package gohome

// Devicer is an interface that describes a type which returns a Device
type Devicer interface {
	FromID(ID string) *Device
}
