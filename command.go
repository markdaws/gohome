package gohome

import "fmt"

type Command interface {
	FriendlyString() string
	fmt.Stringer
}

type FCommand struct {
	Func     func() error
	Friendly string
}

func (c *FCommand) Execute() error {
	return c.Func()
}
func (c *FCommand) FriendlyString() string {
	return c.Friendly
}
func (c *FCommand) String() string {
	return "FCommand"
}

type ZoneSetLevelCommand struct {
	Zone  *Zone
	Level float32
}

func (c *ZoneSetLevelCommand) FriendlyString() string {
	return fmt.Sprintf("Zone [%s] \"%s\" set to %.2f%%", c.Zone.GlobalID, c.Zone.Name, c.Level)
}
func (c *ZoneSetLevelCommand) String() string {
	return "ZoneSetLevelCommand"
}

type ButtonPressCommand struct {
	Button *Button
}

func (c *ButtonPressCommand) FriendlyString() string {
	return fmt.Sprintf("Device [%s] \"%s\" press button %s [%s]",
		c.Button.Device.GlobalID(), c.Button.Device.Name(), c.Button.LocalID, c.Button.GlobalID)
}
func (c *ButtonPressCommand) String() string {
	return "ButtonPressCommand"
}

type ButtonReleaseCommand struct {
	Button *Button
}

func (c *ButtonReleaseCommand) FriendlyString() string {
	return fmt.Sprintf("Device [%s] \"%s\" press release %s [%s]",
		c.Button.Device.GlobalID(), c.Button.Device.Name(), c.Button.LocalID, c.Button.GlobalID)
}
func (c *ButtonReleaseCommand) String() string {
	return "ButtonReleaseCommand"
}

type SceneSetCommand struct {
	Scene *Scene
}

func (c *SceneSetCommand) FriendlyString() string {
	return fmt.Sprintf("Set scene \"%s\" [%s]", c.Scene.Name, c.Scene.GlobalID)
}
func (c *SceneSetCommand) String() string {
	return "SceneSetCommand"
}
