package gohome

import "fmt"

type Action interface {
	//TODO: Return error
	Execute()
}

type PrintAction struct {
}

func (a *PrintAction) Execute() {
	fmt.Println("I am a print action")
}

type FuncAction struct {
	Func func()
}

func (a *FuncAction) Execute() {
	a.Func()
}

type SetSceneAction struct {
	Scene *Scene
}

func (a *SetSceneAction) Execute() {
	a.Scene.Execute()
}

//ZoneSetLevelAction
