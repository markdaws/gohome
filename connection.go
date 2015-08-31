package gohome

import "io"

type Connection interface {
	Open() error
	Close() error
	io.Reader
	io.Writer
	SetConnectionInfo(ConnectionInfo)
}
