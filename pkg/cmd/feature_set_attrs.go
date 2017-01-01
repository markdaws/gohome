package cmd

import (
	"fmt"

	"github.com/markdaws/gohome/pkg/attr"
)

// FeatureSetAttrs indicates there are attributes we need to update on the hardware
type FeatureSetAttrs struct {
	ID          string
	FeatureID   string
	FeatureType string
	FeatureName string
	Attrs       map[string]*attr.Attribute
}

func (c *FeatureSetAttrs) GetID() string {
	return c.ID
}
func (c *FeatureSetAttrs) FriendlyString() string {
	return fmt.Sprintf("FeatureSetAttrs[ID: %s, Type:%s, Name: %s]", c.FeatureID, c.FeatureType, c.FeatureName)
}
func (c *FeatureSetAttrs) String() string {
	return "cmd.FeatureSetAttrs"
}
