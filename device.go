package gohome

import (
	"fmt"
)

type Device struct {
	Identifiable
	System     System
	Connection Connection
	Scenes     []Scene
}

func (d *Device) ExecuteCommand(c Command) {
	fmt.Println(d.Connection)
	fmt.Println("c", c)
	fmt.Println("d", c.Data())
	d.Connection.Send(c.Data())
}
func (d *Device) SetScene(s *Scene) {
	s.Execute()
}
