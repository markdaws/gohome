package cmd

import "fmt"

type SceneSet struct {
	ID        string
	SceneID   string
	SceneName string
}

func (c *SceneSet) FriendlyString() string {
	return fmt.Sprintf("Set scene \"%s\" [%s]", c.SceneName, c.SceneID)
}
func (c *SceneSet) String() string {
	return fmt.Sprintf("cmd.SceneSet: %s, %s", c.SceneID, c.SceneName)
}
