package gohome

type Trigger interface {
	Start() (<-chan bool, <-chan bool)
	Stop()
}
