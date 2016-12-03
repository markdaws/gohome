package gohome

import "github.com/markdaws/gohome/feature"

// Area represents a physical space e.g. Bathroom, garden etc
type Area struct {
	ID          string
	Name        string
	Description string
	Areas       []*Area
	Parent      *Area
	Features    []*feature.Feature
}

func (a *Area) AddArea(area *Area) {
	//TODO:
}

func (a *Area) AddFeature(f *feature.Feature) {
	//TODO:
}
