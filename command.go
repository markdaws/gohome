package gohome

type Command interface {
	Data() []byte
}
