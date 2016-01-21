package gohome

type Devicer interface {
	FromID(ID string) Device
}
