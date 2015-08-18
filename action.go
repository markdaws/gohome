package gohome

import "fmt"

type Action interface {
	Execute()
}

type PrintAction struct {
}

func (a *PrintAction) Execute() {
	fmt.Println("I am a print action")
}
