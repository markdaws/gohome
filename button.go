package gohome

// Button represents a button, which can be an actual physical button on a
// device or it might be a phantom button, which is exposed by a device to
// allow some kind of control action, but does not exist in the physical world
type Button struct {
	Address     string
	ID          string
	Name        string
	Description string
	Device      Device
}
