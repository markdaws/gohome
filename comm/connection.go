package comm

import "io"

type Connection interface {
	Open() error
	Close()
	//TODO: Why add these here?
	io.Reader
	io.Writer

	//TODO: Needed?
	SetPingCallback(PingCallback)
	PingCallback() PingCallback
	Status() ConnectionStatus
	SetStatus(ConnectionStatus)
	Info() ConnectionInfo
}

type PingCallback func() error

type ConnectionStatus string

const (
	CSNew        ConnectionStatus = "new"
	CSConnecting                  = "connecting"
	CSConnected                   = "connected"
	CSClosed                      = "closed"
)
