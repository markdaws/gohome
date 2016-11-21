package testing

import (
	"errors"

	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/cmd"
)

type cmdBuilder struct {
	Device      *gohome.Device
	ModelNumber string
}

func (b *cmdBuilder) Build(c cmd.Command) (*cmd.Func, error) {
	switch b.ModelNumber {
	default:
		return nil, errors.New("unsupported hardware found")
	}
}
