package cmd

import "fmt"

type SceneSet struct {
	SceneGlobalID string
	SceneName     string
}

func (c *SceneSet) FriendlyString() string {
	return fmt.Sprintf("Set scene \"%s\" [%s]", c.SceneName, c.SceneGlobalID)
}
func (c *SceneSet) String() string {
	return "cmd.SceneSet"
}
