package gohome

import "fmt"

type Command interface {
	Execute() error
	FriendlyString() string
	//TODO: Remove
	CMDType() CommandType
	fmt.Stringer
}

type FCommand struct {
	Func func() error
}

func (c *FCommand) Execute() error {
	return c.Func()
}
func (c *FCommand) FriendlyString() string {
	return ""
}
func (c *FCommand) String() string {
	return "FCommand"
}
func (c *FCommand) CMDType() CommandType {
	return CTUnknown
}

//TODO: rename ZoneLetLevelCmd
type ZoneSetLevelCommand struct {
	Zone  *Zone
	Level float32
	Func  func() error
}

func (c *ZoneSetLevelCommand) Execute() error {
	return c.Func()
}
func (c *ZoneSetLevelCommand) FriendlyString() string {
	return "TODO: FriendlyString"
}
func (c *ZoneSetLevelCommand) String() string {
	return "ZoneSetLevelCommand"
}
func (c *ZoneSetLevelCommand) CMDType() CommandType {
	return CTZoneSetLevel
}

type ButtonPressCommand struct {
	Button *Button
	Func   func() error
}

func (c *ButtonPressCommand) Execute() error {
	return c.Func()
}
func (c *ButtonPressCommand) FriendlyString() string {
	return "//TODO: friendly string"
}
func (c *ButtonPressCommand) String() string {
	return "ButtonPressCommand"
}
func (c *ButtonPressCommand) CMDType() CommandType {
	return CTDevicePressButton
}

type ButtonReleaseCommand struct {
	Button *Button
	Func   func() error
}

func (c *ButtonReleaseCommand) Execute() error {
	return c.Func()
}
func (c *ButtonReleaseCommand) FriendlyString() string {
	return "//TODO: friendly string"
}
func (c *ButtonReleaseCommand) String() string {
	return "ButtonReleaseCommand"
}
func (c *ButtonReleaseCommand) CMDType() CommandType {
	return CTDeviceReleaseButton
}

type SceneSetCommand struct {
	Scene *Scene
}

func (c *SceneSetCommand) Execute() error {
	//TODO: Doesn't know what device the command belongs to ...
	/*
		for _, cmd := range c.Scene.Commands {

			//Send command to device ...
			cmd.Execute()
		}*/
	//TODO: error
	return nil
}
func (c *SceneSetCommand) FriendlyString() string {
	return "//TODO: friendly string"
}
func (c *SceneSetCommand) String() string {
	return "SceneSetCommand"
}
func (c *SceneSetCommand) CMDType() CommandType {
	return CTSystemSetScene
}
