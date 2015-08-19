package gohome

type Connection interface {
	Connect() error
	Send([]byte)
}
