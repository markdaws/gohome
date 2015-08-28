package gohome

import "net"

type Connection interface {
	Connect() (net.Conn, error)
	Send([]byte)
}
