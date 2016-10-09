package cmd

import "github.com/markdaws/gohome/comm"

type Builder interface {
	Build(Command) (*Func, error)
	ID() string
	//TODO: Shouldn't be here
	Connections(name, address string) comm.ConnectionPool
}
