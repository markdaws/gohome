package gohome

type Command interface {
	Execute(args ...interface{})
}
