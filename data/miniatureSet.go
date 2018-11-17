package data

import (
	"fmt"
	"strings"
)

type MiniatureSet struct {
	id   string
	name string
}

func (miniSet *MiniatureSet) ID() string {
	return miniSet.id
}
func (miniSet *MiniatureSet) Name() string {
	return miniSet.name
}

// GetMiniatureSetByID returns a set corresponding to the given ID (aka setCode)
func GetMiniatureSetByID(setCode string) (*MiniatureSet, error) {
	miniSet, exists := idToSet[strings.ToUpper(setCode)]
	if !exists {
		return nil, fmt.Errorf("Unable to find set with code %s", setCode)
	}
	return miniSet, nil
}

// GetMiniatureSets returns all of the sets
func GetMiniatureSets() ([]*MiniatureSet, error) {
	var copy []*MiniatureSet
	for i := range sets {
		copy = append(copy, &sets[i])
	}
	return copy, nil
}

var sets []MiniatureSet
var idToSet map[string]*MiniatureSet

func init() {
	sets = []MiniatureSet{
		{id: "A", name: "Anvilborn"},
		{id: "B", name: "Base"},
		{id: "BW", name: "Baxar's War"},
		{id: "CP", name: "Chrysotic Plague"},
		{id: "NF", name: "Night Fusion"},
		{id: "SD", name: "Serrated Dawn"},
	}
	idToSet = make(map[string]*MiniatureSet)
	for i := range sets {
		idToSet[strings.ToUpper(sets[i].id)] = &sets[i]
	}
}
