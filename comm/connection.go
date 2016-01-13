package comm

import "io"

type Connection interface {
	Open() error
	Close()
	io.Reader
	io.Writer
	SetPingCallback(PingCallback)
	PingCallback() PingCallback
	Status() ConnectionStatus
	SetStatus(ConnectionStatus)
}

type PingCallback func() error

type ConnectionStatus string

const (
	CSNew        ConnectionStatus = "new"
	CSConnecting                  = "connecting"
	CSConnected                   = "connected"
	CSClosed                      = "closed"
)
